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

#### Prerequisites

- [Nix](https://nixos.org/download.html) package manager with [flakes](https://wiki.nixos.org/wiki/Flakes) enabled
- [direnv](https://direnv.net/docs/installation.html) for environment management

#### Steps

1. Clone the repository
2. Allow direnv in the project directory:

```bash
direnv allow
```

This will automatically:
- Create a Python 3.12 virtual environment in `.venv`
- Install all dependencies using UV package manager
- Set up the development environment

If you need to manually activate the environment without direnv:

```bash
nix develop
```


## Dependency Management

Dependencies are managed using [UV](https://github.com/astral-sh/uv), a modern Python package manager:

- `pyproject.toml`: Defines project dependencies (requires Python 3.12+)
- `uv.lock`: Locks dependencies to specific versions

Common UV commands:

```bash
# Update dependencies
uv sync

# Update lock file
uv lock

# Install dependencies (for manual setup)
uv install
```

## Usage

### Starting the Dagster UI

Launch the Dagster web interface:

```bash
dagster dev
```

Access the UI at http://localhost:3000

### Pipeline Structure

The pipeline consists of the following components:

1. **Assets**:
   - `sip_asset`: Parses METS XML files into a structured SIP model
   - `intellectual_entities`: Extracts and processes Intellectual Entity models
   - `representations`: Collects and processes file representations
   - `files`: Extracts and processes file metadata
   - `fixities`: Extracts and processes file checksums

2. **Jobs**:
   - `ingest_sip_job`: Orchestrates the complete SIP creation process

3. **Sensors**:
   - `xml_file_sensor`: Monitors for new METS XML files and triggers processing

### Running Tests

Execute the test suite:

```bash
pytest da_pipeline_tests
```

## Project Configuration

- `flake.nix`: Defines the development environment and dependencies
- `.envrc`: Configures direnv to use the Nix flake
- `pyproject.toml`: Defines Python package metadata and dependencies
- `workspace.yaml`: Configures Dagster code locations
- `uv.lock`: Locks dependencies to specific versions
