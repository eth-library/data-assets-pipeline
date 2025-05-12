# Data Archive Assets Pipeline

## Overview

This project implements a digital asset processing pipeline that implements the Submission Information Package (SIP)
component of the Open Archival Information System (OAIS) reference model and METS (Metadata Encoding and Transmission
Standard) specifications. It processes and manages digital assets within a data archive by extracting metadata from METS
files and organizing them into structured SIPs.

The system utilizes [Dagster](https://dagster.io/) as its core data orchestrator, providing robust workflow management
for complex archiving processes. The implementation ensures:

- **OAIS SIP Processing**: Implements the OAIS Submission Information Package model with structured metadata handling
- **METS Standard Support**: Full parsing and processing of METS XML files
- **Data Validation**: Robust validation using Pydantic models
- **Scalable Architecture**: Modular design for handling complex archiving workflows

## Installation

1. Clone the repository
2. Install the package with development dependencies:

```bash
pip install -e ".[dev]"
```

Required dependencies:

- [dagster](https://pypi.org/project/dagster/): Core orchestration framework
- [dagster-cloud](https://pypi.org/project/dagster-cloud/): Cloud deployment support
- [dagster-webserver](https://pypi.org/project/dagster-webserver/): Development UI
- [pydantic](https://pypi.org/project/pydantic/): Data validation and modeling
- [pytest](https://pypi.org/project/pytest/): Testing framework

## Usage

### Starting the Dagster UI

Launch the Dagster web interface:

```bash
dagster dev
```

Access the UI at http://localhost:3000

### Running Tests

Execute the test suite:

```bash
pytest da_pipeline_tests
```

## Project Configuration

The project configuration is split between package installation and runtime settings:

**Package Installation** 

```bash
pip install -e ".[dev]"
```

- `pyproject.toml`: Specifies the build system requirements and configuration for Python projects. As established by the
 [PEP 518](https://peps.python.org/pep-0518/) standard.
- [setuptools](https://setuptools.pypa.io/en/latest/userguide/declarative_config.html): Tools for package configuration and distribution:
  - `setup.py`: Script for defining package dependencies, installation settings, and other setup instructions.
  - `setup.cfg`: Declarative configuration file for package metadata and settings, enabling configuration without extensive Python code.

**Dagster Runtime Configuration**:

- `workspace.yaml`: Configures code locations for Dagster, telling it where to find your code and how to load it. For
  more details, refer to
  the [Dagster documentation on workspace.yaml](https://docs.dagster.io/guides/deploy/code-locations/workspace-yaml).