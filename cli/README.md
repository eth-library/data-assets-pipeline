# dap CLI

Developer tools for the **Data Archive Pipeline (DAP) Orchestrator** at ETH Library Zurich.

The `dap` CLI wraps common development tasks into short, memorable commands. It doesn't replace the underlying tools — it provides ergonomic shortcuts and a consistent interface.

## The dev toolchain stack

The project uses a layered toolchain where each layer builds on the one below:

```
nix flakes        Reproducible packages — pins exact versions of Python, uv, kubectl, etc.
  └─ direnv       Auto-loading shell env — activates when you cd into the project
      └─ nix-direnv   Cached flake evaluation — avoids re-evaluating the flake on every shell
          └─ uv       Fast Python deps — installs packages from the lockfile in milliseconds
              └─ dap   Ergonomic commands — wraps pytest, ruff, mypy, helm, kubectl
```

**Why each layer exists:**

- **Nix flakes** (`flake.nix`): ensures every developer has identical tool versions. No "works on my machine".
- **direnv** (`.envrc`): automatically loads the nix environment when you enter the project directory. No manual `nix develop`.
- **nix-direnv**: caches the nix evaluation so shell startup stays fast. Without it, entering the directory would re-evaluate the flake every time.
- **uv** (`pyproject.toml`, `uv.lock`): manages Python dependencies. Fast, deterministic, lockfile-based.
- **dap CLI** (`cli/`): the commands documented below. Wraps quality checks, environment management, and Kubernetes deployment.

## Installation

Handled automatically by direnv and nix:

```bash
git clone <repo-url> && cd dap
direnv allow  # one-time approval
# Everything is now available: dap, python, uv, dagster, kubectl, helm
```

## Commands

### Development

| Command | Description |
|---------|-------------|
| `dap test [--scope core\|cli\|all]` | Run tests with pytest |
| `dap lint [--fix] [--scope ...]` | Check code style and formatting with ruff |
| `dap typecheck [--scope ...]` | Run type checking with mypy |
| `dap check [--scope ...]` | Run all quality checks (ruff, mypy, pytest) |

The `--scope` flag controls which code is checked:
- `core` (default): `da_pipeline`, `da_pipeline_tests`
- `cli`: `cli/dap_cli`, `cli/tests`
- `all`: both core and CLI

### Environment

| Command | Description |
|---------|-------------|
| `dap welcome` | Show welcome banner and environment info |
| `dap tools [--all]` | Show installed tool versions and paths |
| `dap clean [--yes]` | Remove `.venv` and caches |
| `dap reset [--yes]` | Clean and reinstall dependencies |

`dap tools` shows Python toolchain by default. Pass `--all` to include nix, direnv, kubectl, and helm.

`dap clean` and `dap reset` prompt for confirmation. Pass `--yes` / `-y` to skip (for CI/scripts).

### Kubernetes

| Command | Description |
|---------|-------------|
| `dap k8s up` | Build and deploy to local Kubernetes |
| `dap k8s down [--yes]` | Tear down deployment |
| `dap k8s restart` | Rebuild image and rollout restart |
| `dap k8s status` | Show pods and services |
| `dap k8s logs` | Stream user code pod logs |
| `dap k8s shell` | Interactive shell in user code pod |

Requires Docker Desktop with Kubernetes enabled.

### Hint commands

These commands show how to use tools that are available directly in your shell:

| Command | Description |
|---------|-------------|
| `dap uv` | Common uv commands (sync, lock, add, run) |
| `dap dagster` | Common dagster/dg commands |
| `dap direnv` | Common direnv commands (allow, reload, status) |

### Global flags

| Flag | Description |
|------|-------------|
| `--version` / `-V` | Show CLI version |
| `--help` | Show help |

## Environment variables

| Variable | Effect |
|----------|--------|
| `DAP_THEME=light\|dark` | Override terminal background detection for colours |
| `DAP_QUIET=1` | Suppress Quick Start section in `dap welcome` |
| `NO_COLOR=1` | Disable all colour output (also respected in CI) |

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for CLI architecture, how to add commands, and testing patterns.