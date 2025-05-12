# Data Archive Assets Pipeline

## Overview

This project implements a digital asset processing pipeline that implements the Submission Information Package (SIP)
component of the Open Archival Information System (OAIS) reference model and METS (Metadata Encoding and Transmission
Standard) specifications. It processes and manages digital assets within a data archive by extracting metadata from METS
files and organizing them into structured SIPs.

The system uses [Dagster](https://dagster.io/) as its core data orchestrator, providing robust workflow management
for complex archiving processes. The implementation ensures:

- **OAIS SIP Processing**: Implements the OAIS Submission Information Package model with structured metadata handling
- **METS Standard Support**: Full parsing and processing of METS XML files
- **Data Validation**: Robust validation using Pydantic models
- **Scalable Architecture**: Modular design for handling complex archiving workflows

## Setup

### Recommended Nix + Direnv Setup

We recommend using the fully automatic setup method using Nix Flakes and Direnv:

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
