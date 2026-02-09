# Data Archive Pipeline (DAP) Orchestrator

A [Dagster](https://dagster.io/)-based orchestrator for processing digital assets following the OAIS reference model. Manages Submission Information Packages (SIPs) with support for multiple metadata standards including Dublin Core, METS, and PREMIS for comprehensive digital preservation.
> **Domain Terms**
>
> - **METS** (Metadata Encoding and Transmission Standard): XML schema for encoding metadata about digital objects
> - **OAIS** (Open Archival Information System): ISO reference model for digital preservation systems
> - **SIP** (Submission Information Package): A package submitted to an archive for ingest
> - **IE** (Intellectual Entity): A distinct unit of information to be preserved (e.g., a document, dataset)

## Features

- Parses METS XML files into validated Pydantic models following OAIS standards
- Extracts Dublin Core descriptive metadata and PREMIS preservation metadata
- Validates file fixity information (MD5, SHA-1, SHA-256, SHA-512 checksums)
- Automated file monitoring via Dagster sensors
- Hot-reload development with instant feedback
- Kubernetes-ready with Helm chart configuration

> ### Project Status Notice
> This project is in **alpha** and actively evolving.
>
> - **Stability:** Bugs, breaking changes, or incomplete features may occur.
> - **Evolution:** APIs and behaviors can change as we refine functionality.
>
> Provided as‑is, with no guarantees on stability.
>
> We appreciate your patience and feedback!

## Quick Start

### With Nix (Recommended)

Clone the repository and enter the development environment:

```bash
git clone https://github.com/eth-library/data-assets-pipeline.git
cd data-assets-pipeline
direnv allow
```

If not using [direnv](https://direnv.net/), activate the environment manually:

```bash
nix develop
```

Start the Dagster development server:

```bash
dap dev
```

Open http://localhost:3000 in your browser.

### Without Nix

Ensure you have the following installed:

- Python 3.12+
- [uv](https://github.com/astral-sh/uv) (Python package manager)
- [Go 1.22+](https://go.dev/dl/) (for the dap CLI)

Set up the project:

```bash
git clone https://github.com/eth-library/data-assets-pipeline.git
cd data-assets-pipeline
uv venv --python python3.12
source .venv/bin/activate
uv sync --extra dev

# Build the dap CLI
cd .cli && go build -o bin/dap . && cd ..
export PATH="$PWD/.cli/bin:$PATH"
```

Start the Dagster development server:

```bash
dap dev
```

Open http://localhost:3000 in your browser.

## Prerequisites

| Tool | Purpose | Installation |
|------|---------|--------------|
| [Nix](https://nixos.org/download.html) | Development environment (optional but recommended) | [Install guide](https://nixos.org/download.html) |
| [direnv](https://direnv.net/) | Automatic environment loading | [Install guide](https://direnv.net/docs/installation.html) |
| [nix-direnv](https://github.com/nix-community/nix-direnv) | Fast Nix + direnv integration | [Install guide](https://github.com/nix-community/nix-direnv#installation) |
| Python 3.12+ | Runtime | [python.org](https://www.python.org/downloads/) |
| [uv](https://github.com/astral-sh/uv) | Fast Python package manager | `curl -LsSf https://astral.sh/uv/install.sh \| sh` |
| [Go 1.22+](https://go.dev/dl/) | dap CLI build | [go.dev](https://go.dev/dl/) |

For Kubernetes development, you'll also need:

- Docker Desktop with Kubernetes enabled
- kubectl
- helm

## Architecture

The pipeline processes XML files through a sequence of Dagster [assets](https://docs.dagster.io/concepts/assets/software-defined-assets), each extracting a specific layer of the preservation package hierarchy:

### Assets

| Asset | Input | Output | Description                                                    |
|-------|-------|--------|----------------------------------------------------------------|
| `sip_asset` | XML file paths | `SIPModel` | Parses METS XML into a structured SIP model                    |
| `intellectual_entities` | `SIPModel` | `list[IEModel]` | Extracts intellectual entity metadata with Dublin Core         |
| `representations` | `list[IEModel]` | `list[RepresentationModel]` | Collects file representations (preservation, access, original) |
| `files` | `list[RepresentationModel]` | `list[FileModel]` | Extracts file metadata including MIME types and paths          |
| `fixities` | `list[FileModel]` | `list[FixityModel]` | Extracts and validates file checksums                          |

### Sensor

The `xml_file_sensor` monitors a configured directory for new XML files and automatically triggers the pipeline. By default, it watches `da_pipeline_tests/test_data/` every 30 seconds.

## Usage

### Local Development

Start the Dagster UI with hot-reload:

```bash
dap dev
```

Run the test suite:

```bash
dap test
```

### Manual Pipeline Runs

To run the pipeline manually via the Dagster UI launchpad, provide a run configuration:

```yaml
ops:
  sip_asset:
    config:
      file_paths:
        - /path/to/your/mets.xml
```

### Kubernetes Deployment

Requires Docker Desktop with Kubernetes enabled (Settings → Kubernetes → Enable).

Deploy to local Kubernetes:

```bash
dap k8s up
```

The Dagster UI will be available at http://localhost:8080.

Rebuild and restart after code changes:

```bash
dap k8s restart
```

Tear down the deployment:

```bash
dap k8s down
```

## Configuration

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `DAGSTER_HOME` | Dagster instance directory | Project root (set by `.envrc`) |
| `DAGSTER_TEST_DATA_PATH` | Directory containing METS XML files for the sensor to monitor | `da_pipeline_tests/test_data` |
| `DAP_QUIET` | Set to `1` to hide the "Commands" hint in welcome banner | unset |

Copy `.env.example` to `.env` and modify as needed.

### Configuration Files

| File | Purpose |
|------|---------|
| `flake.nix` | Nix development environment (Python, uv, Go, kubectl, helm) |
| `.cli/` | Go CLI source code (see [.cli/CONTRIBUTING.md](.cli/CONTRIBUTING.md)) |
| `pyproject.toml` | Python project metadata and dependencies |
| `dagster.yaml` | Dagster instance configuration |
| `config.yaml` | Example run configuration for manual pipeline execution |
| `helm/values.yaml` | Base Helm values for Kubernetes deployment |
| `helm/values-local.yaml` | Local Kubernetes overrides |
| `helm/pvc.yaml` | Persistent volume claim for Dagster storage |

## Commands Reference

Run `dap --help` to see all available commands.

| Command | Description |
|---------|-------------|
| `dap dev` | Start Dagster development server (localhost:3000) |
| `dap test` | Run pytest test suite |
| `dap check` | Run all quality checks (lint + typecheck + test) |
| `dap lint` | Check code style with Ruff (use `--fix` to auto-format) |
| `dap typecheck` | Type check with mypy |
| `dap versions` | Show tool versions |
| `dap clean` | Remove `.venv` and caches |
| `dap reset` | Clean and reinstall dependencies |
| `dap materialize` | Materialize all Dagster assets |
| `dap run` | Run the ingest_sip_job |
| `dap k8s up` | Build and deploy to local Kubernetes (localhost:8080) |
| `dap k8s down` | Tear down Kubernetes deployment |
| `dap k8s restart` | Rebuild and restart user code pod |
| `dap k8s status` | Show pods and services |
| `dap k8s logs` | Stream logs from user code pod |
| `dap k8s shell` | Open shell in user code pod |

For CLI development, see [.cli/CONTRIBUTING.md](.cli/CONTRIBUTING.md).

## Project Structure

```
.cli/                        # dap CLI (Go) - see .cli/CONTRIBUTING.md
├── cmd/                     # Command implementations
├── internal/                # Internal packages (ui, exec)
└── main.go                  # Entry point

da_pipeline/                 # Main package (Python)
├── definitions.py           # Dagster entry point (Definitions)
├── assets.py                # Pipeline assets
├── sensors.py               # File monitoring sensor and job definition
├── mets_parser.py           # METS XML parsing logic
├── pydantic_models.py       # OAIS-compliant data models
└── utils.py                 # Helper functions for metadata collection

da_pipeline_tests/           # Tests
├── test_data/               # Sample METS XML files
│   └── synthetic_sip.xml    # Synthetic test data
├── test_assets.py           # Asset tests
└── conftest.py              # Pytest fixtures

helm/                        # Kubernetes configuration
├── values.yaml              # Base Helm values
├── values-local.yaml        # Local development overrides
└── pvc.yaml                 # Persistent volume claim
```

## Dependency Management

Dependencies are managed with [uv](https://github.com/astral-sh/uv):

Install or update dependencies:

```bash
uv sync
```

Update the lock file:

```bash
uv lock
```

Add a new dependency:

```bash
uv add <package>
```

## Troubleshooting

### Dagster UI not loading

Ensure no other process is using port 3000 (local) or 8080 (Kubernetes):

```bash
lsof -i :3000
```

### Kubernetes deployment fails

Verify Kubernetes is running:

```bash
kubectl cluster-info
```

Ensure Docker Desktop Kubernetes is enabled in Settings → Kubernetes → Enable Kubernetes.

### Sensor not detecting files

Check the `DAGSTER_TEST_DATA_PATH` environment variable points to a directory containing `.xml` files. The sensor only processes files with the `.xml` extension.

### Test failures

Ensure dependencies are installed with dev extras:

```bash
uv sync --extra dev
```

## License

[Apache License 2.0](LICENSE) - © 2026 ETH Zurich

## Further Reading

### Core Technologies

- [Dagster Documentation](https://docs.dagster.io/) - Data orchestration platform
  - [Assets](https://docs.dagster.io/concepts/assets/software-defined-assets)
  - [Sensors](https://docs.dagster.io/concepts/partitions-schedules-sensors/sensors)
  - [Definitions](https://docs.dagster.io/concepts/code-locations)
- [Pydantic Documentation](https://docs.pydantic.dev/) - Data validation library

### Digital Preservation Standards

- [METS (Metadata Encoding and Transmission Standard)](https://www.loc.gov/standards/mets/)
- [OAIS Reference Model (ISO 14721)](https://www.iso.org/standard/57284.html)
- [Dublin Core Metadata](https://www.dublincore.org/specifications/dublin-core/)
- [PREMIS (Preservation Metadata)](https://www.loc.gov/standards/premis/)

### Development Tools

- [uv](https://github.com/astral-sh/uv) - Fast Python package manager
- [Go](https://go.dev/) - Programming language for the dap CLI
- [Cobra](https://github.com/spf13/cobra) - Go CLI framework
- [Charmbracelet](https://charm.sh/) - Terminal UI libraries
- [Nix](https://nixos.org/) - Reproducible development environments
- [direnv](https://direnv.net/) - Automatic environment loading
- [Ruff](https://docs.astral.sh/ruff/) - Python linter and formatter

### Kubernetes

- [Helm](https://helm.sh/docs/) - Kubernetes package manager
- [Dagster Helm Chart](https://docs.dagster.io/deployment/guides/kubernetes/deploying-with-helm)
