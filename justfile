# Justfile for da_pipeline development
# Run `just` or `just --list` to see available commands

set shell := ["bash", "-euo", "pipefail", "-c"]

# Configuration
namespace := "dagster"
release := "dagster"
image := "da-pipeline:local"
helm_chart := "dagster/dagster"
helm_version := "1.10.14"

# Default: show available commands
default:
    @just --list

# Show environment info (versions + commands)
info:
    @echo "da_pipeline"
    @echo "==========="
    @just versions
    @echo ""
    @just --list

# Show tool versions
versions:
    @echo "python:  $(python --version 2>&1 | cut -d' ' -f2)"
    @echo "uv:      $(uv --version 2>&1 | cut -d' ' -f2 || echo 'n/a')"
    @echo "just:    $(just --version 2>&1 | cut -d' ' -f2 || echo 'n/a')"
    @echo "kubectl: $(kubectl version --client -o json 2>/dev/null | python -c 'import sys,json; print(json.load(sys.stdin)["clientVersion"]["gitVersion"])' 2>/dev/null || echo 'n/a')"
    @echo "helm:    $(helm version --short 2>/dev/null | cut -d'+' -f1 || echo 'n/a')"
    @echo "dagster: $(python -c 'import dagster; print(dagster.__version__)' 2>/dev/null || echo '⚠ not installed - run: just setup')"

# ═══════════════════════════════════════════════════════════════════════════════
# Local Development
# ═══════════════════════════════════════════════════════════════════════════════

# Setup Python environment (creates venv if needed, syncs deps)
setup:
    #!/usr/bin/env bash
    set -euo pipefail
    [ ! -d .venv ] && uv venv --python "$(which python)"
    [ ! -f uv.lock ] && uv lock
    source .venv/bin/activate
    uv sync --extra dev --quiet
    echo "Done."

# Start Dagster development server (localhost:3000)
dev:
    dagster dev

# Run tests
test *ARGS:
    pytest da_pipeline_tests {{ ARGS }}

# ═══════════════════════════════════════════════════════════════════════════════
# Kubernetes (k8s-*)
# ═══════════════════════════════════════════════════════════════════════════════

# Deploy to local Kubernetes (localhost:8080)
k8s-up: _check-k8s _build _helm-up
    #!/usr/bin/env bash
    set -euo pipefail
    echo ""
    just k8s-status
    echo ""
    pkill -f "kubectl.*port-forward.*8080" 2>/dev/null || true
    # Start port-forward in background with retries (suppress all output)
    (while ! kubectl port-forward svc/dagster-dagster-webserver 8080:80 -n {{ namespace }} &>/dev/null; do sleep 2; done) &

    echo "Waiting for webserver to start..."
    attempt=1
    while true; do
        if curl -s -o /dev/null -w '' --connect-timeout 1 http://localhost:8080 2>/dev/null; then
            echo ""
            echo "✓ Dagster UI is ready!"
            echo "  URL: http://localhost:8080"
            exit 0
        fi
        echo "[$(date +%H:%M:%S)] Attempt $attempt - Waiting for webserver... (retrying in 2s)"
        attempt=$((attempt + 1))
        sleep 2
    done

# Tear down Kubernetes deployment
k8s-down: _check-k8s
    #!/usr/bin/env bash
    set -euo pipefail
    echo "Tearing down..."
    helm uninstall {{ release }} -n {{ namespace }} --wait=false 2>/dev/null || true
    # Delete all dagster run jobs and pods (prevents PVC from being stuck)
    kubectl delete jobs -n {{ namespace }} -l dagster/run-id --timeout=10s 2>/dev/null || true
    kubectl delete pods -n {{ namespace }} -l dagster/run-id --grace-period=0 --force --timeout=10s 2>/dev/null || true
    kubectl delete pvc dagster-storage -n {{ namespace }} --timeout=10s 2>/dev/null || true
    kubectl delete configmap test-data-xml -n {{ namespace }} --timeout=10s 2>/dev/null || true
    echo "Done."

# Rebuild and restart user code pod
k8s-restart: _check-k8s _build
    @kubectl rollout restart deployment -n {{ namespace }} -l "app.kubernetes.io/name=dagster-user-deployments"
    @kubectl rollout status deployment -n {{ namespace }} -l "app.kubernetes.io/name=dagster-user-deployments" --timeout=120s
    @echo "Restarted."

# Show pods and services
k8s-status:
    #!/usr/bin/env bash
    set -euo pipefail
    if ! kubectl get namespace {{ namespace }} &>/dev/null; then
        echo "Kubernetes deployment not running. Use 'just k8s-up' to deploy."
        exit 0
    fi
    echo "=== Pods ==="
    kubectl get pods -n {{ namespace }} -o wide 2>/dev/null || echo "No pods"
    echo ""
    echo "=== Services ==="
    kubectl get svc -n {{ namespace }} 2>/dev/null || echo "No services"

# Stream logs from user code pod
k8s-logs *ARGS:
    kubectl logs -n {{ namespace }} -l "app.kubernetes.io/name=dagster-user-deployments" --tail=100 -f {{ ARGS }}

# Port-forward to Dagster UI (localhost:8080)
k8s-ui:
    #!/usr/bin/env bash
    set -euo pipefail
    if lsof -i :8080 -sTCP:LISTEN &>/dev/null; then
        echo "✓ Port-forward already active: http://localhost:8080"
        exit 0
    fi

    # Start port-forward in background
    (kubectl port-forward svc/dagster-dagster-webserver 8080:80 -n {{ namespace }} &>/dev/null &)

    echo "Waiting for webserver to start..."
    attempt=1
    while true; do
        if curl -s -o /dev/null -w '' --connect-timeout 1 http://localhost:8080 2>/dev/null; then
            echo ""
            echo "✓ Dagster UI is ready!"
            echo "  URL: http://localhost:8080"
            exit 0
        fi
        echo "[$(date +%H:%M:%S)] Waiting for webserver... (retrying in 2s)"
        sleep 2
    done

# Open shell in user code pod
k8s-shell:
    kubectl exec -it -n {{ namespace }} $(kubectl get pods -n {{ namespace }} -l "app.kubernetes.io/name=dagster-user-deployments" -o jsonpath='{.items[0].metadata.name}') -- /bin/bash

# ═══════════════════════════════════════════════════════════════════════════════
# Internal recipes (prefixed with _)
# ═══════════════════════════════════════════════════════════════════════════════

[private]
_check-k8s:
    @kubectl cluster-info &>/dev/null || (echo "Error: Kubernetes not available. Enable it in Docker Desktop." && exit 1)
    @kubectl config use-context docker-desktop &>/dev/null || true

[private]
_build:
    @echo "Building {{ image }}..."
    @docker build -t {{ image }} -q .
    @echo "Built."

[private]
_helm-up:
    #!/usr/bin/env bash
    set -euo pipefail

    # Ensure namespace
    kubectl create namespace {{ namespace }} --dry-run=client -o yaml | kubectl apply -f -

    # Setup helm repo (only update if not updated in last hour)
    helm repo add dagster https://dagster-io.github.io/helm 2>/dev/null || true
    REPO_CACHE="$HOME/.cache/helm/repository/dagster-index.yaml"
    if [[ ! -f "$REPO_CACHE" ]] || [[ $(find "$REPO_CACHE" -mmin +60 2>/dev/null) ]]; then
        helm repo update dagster >/dev/null
    fi

    # Apply PVC from manifest (creates or updates)
    kubectl apply -n {{ namespace }} -f helm/pvc.yaml

    # Create ConfigMap for test data (only if changed)
    if ls da_pipeline_tests/test_data/*.xml &>/dev/null; then
        kubectl create configmap test-data-xml --from-file=da_pipeline_tests/test_data/ -n {{ namespace }} --dry-run=client -o yaml | kubectl apply -f -
    fi

    # Deploy with Helm (no --wait for faster return)
    echo "Deploying Dagster..."
    helm upgrade --install {{ release }} {{ helm_chart }} \
        -f helm/values.yaml \
        -f helm/values-local.yaml \
        -n {{ namespace }} \
        --version {{ helm_version }} \
        --skip-schema-validation
    echo "Deployed (pods starting in background)."
