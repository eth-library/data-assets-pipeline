#!/bin/bash
set -euo pipefail

# Only run in Claude Code on the web
if [ "${CLAUDE_CODE_REMOTE:-}" != "true" ]; then
  exit 0
fi

PROJECT_DIR="${CLAUDE_PROJECT_DIR:-$(pwd)}"
NIX_INSTALLER="https://nixos.org/nix/install"

# ─── Install Nix (single-user, with flakes) ──────────────────────────────────
if ! command -v nix &>/dev/null && ! ls /nix/store/*/bin/nix &>/dev/null; then
  echo "Installing Nix (single-user)..."
  curl -fsSL "$NIX_INSTALLER" | bash -s -- --no-daemon 2>&1 || true
fi

# Find the nix binary (it may not be on PATH if profile setup crashed)
NIX_BIN=""
if command -v nix &>/dev/null; then
  NIX_BIN="$(command -v nix)"
elif [ -f "$HOME/.nix-profile/bin/nix" ]; then
  NIX_BIN="$HOME/.nix-profile/bin/nix"
else
  # Find nix in the store directly (profile setup may have crashed in sandbox)
  NIX_BIN="$(ls -t /nix/store/*/bin/nix 2>/dev/null | head -1)"
fi

if [ -z "$NIX_BIN" ]; then
  echo "ERROR: Could not find nix binary"
  exit 1
fi

echo "Using nix: $NIX_BIN ($($NIX_BIN --version 2>/dev/null | tail -1))"

# ─── Fetch flake dependencies into /nix/store ────────────────────────────────
# nix derivation show evaluates the flake and fetches all store paths from the
# binary cache without needing to build locally. This works around the
# sandbox's PID restriction which prevents nix build/develop from completing.
echo "Fetching Nix flake dependencies..."
cd "$PROJECT_DIR"
DRV_JSON=$($NIX_BIN --extra-experimental-features "nix-command flakes" \
  derivation show .#devShells.x86_64-linux.default 2>/dev/null)

# ─── Extract package paths and LD_LIBRARY_PATH from derivation ───────────────
NATIVE_BUILD_INPUTS=$(echo "$DRV_JSON" | python3 -c "
import json, sys
data = json.load(sys.stdin)
for drv in data.get('derivations', data).values():
    if isinstance(drv, dict) and 'env' in drv:
        print(drv['env'].get('nativeBuildInputs', ''))
        break
" 2>/dev/null)

NIX_LD_LIBRARY_PATH=$(echo "$DRV_JSON" | python3 -c "
import json, sys
data = json.load(sys.stdin)
for drv in data.get('derivations', data).values():
    if isinstance(drv, dict) and 'env' in drv:
        print(drv['env'].get('LD_LIBRARY_PATH', ''))
        break
" 2>/dev/null)

# Build PATH from nix store packages (find bin dirs for each package)
NIX_PATH_ENTRIES=""
for pkg in $NATIVE_BUILD_INPUTS; do
  if [ -d "$pkg/bin" ]; then
    NIX_PATH_ENTRIES="${NIX_PATH_ENTRIES:+$NIX_PATH_ENTRIES:}$pkg/bin"
  fi
  # For -dev packages, look for the companion -bin package
  if [[ "$pkg" == *-dev ]]; then
    base_name=$(echo "$pkg" | sed 's|-dev$|-bin|; s|/nix/store/[a-z0-9]*-||')
    bin_pkg=$(ls -d /nix/store/*-"$base_name" 2>/dev/null | head -1)
    if [ -n "$bin_pkg" ] && [ -d "$bin_pkg/bin" ]; then
      NIX_PATH_ENTRIES="${NIX_PATH_ENTRIES:+$NIX_PATH_ENTRIES:}$bin_pkg/bin"
    fi
  fi
done

export PATH="$NIX_PATH_ENTRIES:$PATH"
export LD_LIBRARY_PATH="${NIX_LD_LIBRARY_PATH:-}"

echo "Nix tools available:"
echo "  python: $(python3 --version 2>&1)"
echo "  uv:     $(uv --version 2>&1)"
echo "  just:   $(just --version 2>&1)"

# ─── Create venv and install Python deps ─────────────────────────────────────
if [ ! -d "$PROJECT_DIR/.venv" ]; then
  echo "Creating Python virtual environment..."
  uv venv --quiet --directory "$PROJECT_DIR"
fi

echo "Syncing Python dependencies..."
source "$PROJECT_DIR/.venv/bin/activate"
uv sync --extra dev --frozen --quiet --directory "$PROJECT_DIR" 2>&1

# ─── Export environment for the Claude session ───────────────────────────────
if [ -n "${CLAUDE_ENV_FILE:-}" ]; then
  {
    echo "export DAGSTER_HOME=\"$PROJECT_DIR\""
    echo "export LD_LIBRARY_PATH=\"$NIX_LD_LIBRARY_PATH\""
    echo "export PATH=\"$PROJECT_DIR/.venv/bin:$NIX_PATH_ENTRIES:\$PATH\""
  } >> "$CLAUDE_ENV_FILE"
fi

echo "Session hook complete — Nix flake environment loaded, Python deps installed."
