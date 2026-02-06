#!/bin/bash
set -euo pipefail

# Only run in Claude Code on the web
[ "${CLAUDE_CODE_REMOTE:-}" = "true" ] || exit 0

PROJECT_DIR="${CLAUDE_PROJECT_DIR:-$(pwd)}"

# ─── Install Nix ─────────────────────────────────────────────────────────────
# Profile setup may crash in sandboxed containers — that's OK, the store works.
if ! command -v nix &>/dev/null; then
  curl -fsSL https://nixos.org/nix/install | bash -s -- --no-daemon 2>&1 || true
fi

NIX="$(command -v nix 2>/dev/null || echo "$HOME/.nix-profile/bin/nix")"
[ -x "$NIX" ] || NIX="$(ls /nix/store/*/bin/nix 2>/dev/null | head -1)"
[ -x "$NIX" ] || { echo "nix not found"; exit 1; }

# ─── Load flake packages into PATH ───────────────────────────────────────────
cd "$PROJECT_DIR"
DRV=$("$NIX" --extra-experimental-features "nix-command flakes" \
  derivation show .#devShells.x86_64-linux.default 2>/dev/null)

eval "$(echo "$DRV" | python3 -c '
import json, sys, os, glob
env = next(v["env"] for v in json.load(sys.stdin).get("derivations",{}).values()
           if isinstance(v, dict) and "env" in v)
paths = []
for pkg in env.get("nativeBuildInputs", "").split():
    if os.path.isdir(pkg + "/bin"):
        paths.append(pkg + "/bin")
    if pkg.endswith("-dev"):
        name = pkg.split("/")[-1].replace("-dev", "-bin")
        paths.extend(glob.glob("/nix/store/*-" + name + "/bin"))
sep = ":"
print("NIX_PATHS=" + repr(sep.join(paths)))
print("NIX_LDPATH=" + repr(env.get("LD_LIBRARY_PATH", "")))
')"

export PATH="$NIX_PATHS:$PATH"
export LD_LIBRARY_PATH="$NIX_LDPATH"

# ─── Setup project using justfile ────────────────────────────────────────────
just setup

# ─── Export for the Claude session ───────────────────────────────────────────
[ -n "${CLAUDE_ENV_FILE:-}" ] && cat >> "$CLAUDE_ENV_FILE" <<EOF
export DAGSTER_HOME="$PROJECT_DIR"
export LD_LIBRARY_PATH="$NIX_LDPATH"
export PATH="$PROJECT_DIR/.venv/bin:$NIX_PATHS:\$PATH"
EOF

echo "Done — nix flake loaded, just setup complete."
