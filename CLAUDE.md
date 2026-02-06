# Claude Code — Project Setup

This project uses a **Nix flake** (`flake.nix`) to provide all development tools. A SessionStart hook (`.claude/hooks/session-start.sh`) installs Nix and loads the flake automatically in Claude Code on the web.

## Environment

After the hook runs, these tools are available from the Nix flake:

- `python` (3.12), `uv`, `just`, `jq`, `curl`, `openssl`, `kubectl`, `helm`
- Python venv at `.venv/` with all dev dependencies installed

## Commands

Use the justfile for all development tasks:

```bash
just test          # run tests (pytest)
just lint          # check code (ruff check + format --check)
just fmt           # auto-fix and format code
just dev           # start Dagster dev server (localhost:3000)
just setup         # recreate venv and sync dependencies
```

## Project structure

- `da_pipeline/` — main package (Dagster assets, sensors, METS parser, Pydantic models)
- `da_pipeline_tests/` — tests and test data
- `flake.nix` — Nix dev environment definition
- `pyproject.toml` — Python project config and dependencies (managed with `uv`)

## Key notes

- Always run `just lint` before committing
- The `.envrc` + `direnv` setup is for local development; Claude Code uses the SessionStart hook instead
- `DAGSTER_HOME` is set to the project root
