# Data Archive Assets Pipeline

## Overview

This project aims to process digital assets within a data archive by extracting metadata
from [METS (Metadata Encoding and Transmission Standard)](https://www.loc.gov/standards/mets/) files and assigning
values to assets using a data model tailored for long-term archiving.

The workflows are managed using [Dagster](https://dagster.io/), a robust data orchestrator designed for machine
learning, analytics, and ETL processes. The project is built to be **modular**, **scalable**, and **extensible**, making
it well-suited for handling complex archiving workflows.

## Getting started

First, install your Dagster code location as a Python package. By using the --editable flag, pip will install your
Python package in ["editable mode"](https://pip.pypa.io/en/latest/topics/local-project-installs/#editable-installs) so
that as you develop, local code changes will automatically apply.

```bash
pip install -e ".[dev]"
```

Then, start the Dagster UI web server:

```bash
dagster dev
```

Open http://localhost:3000 with your browser to see the project.

You can start writing assets in `da_pipeline/assets.py`. The assets are automatically loaded into the Dagster code
location as you define them.

## Development

### Adding new Python dependencies

You can specify new Python dependencies in `setup.py`.

### Unit testing

Tests are in the `da_pipeline_tests` directory and you can run tests using `pytest`:

```bash
pytest da_pipeline_tests
```

### Schedules and sensors

If you want to enable Dagster [Schedules](https://docs.dagster.io/concepts/partitions-schedules-sensors/schedules)
or [Sensors](https://docs.dagster.io/concepts/partitions-schedules-sensors/sensors) for your jobs,
the [Dagster Daemon](https://docs.dagster.io/deployment/dagster-daemon) process must be running. This is done
automatically when you run `dagster dev`.

Once your Dagster Daemon is running, you can start turning on schedules and sensors for your jobs.