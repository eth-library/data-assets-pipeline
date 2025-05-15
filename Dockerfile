# Base image - using Python 3.12 as specified in pyproject.toml (requires-python >= 3.12)
FROM python:3.12

# Update system packages
RUN apt-get update && \
    apt-get upgrade -y && \
    apt-get install -y && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

# Copy only the files not excluded by .dockerignore into the root directory of the container
COPY . /

# Configure Python to run in unbuffered mode
# This makes sure that logs show up immediately instead of being buffered
ENV PYTHONUNBUFFERED=1

# Update pip to the latest version
RUN pip install --upgrade pip

# Install the project with production dependencies
# This will install:
# - Core dependencies
# - Production dependencies
# as specified in pyproject.toml
RUN pip install .[production]

# Set the working directory to the main package directory
# This aligns with the module_name in pyproject.toml [tool.dagster] section
WORKDIR /da_pipeline/

# Expose port 80 for the Dagster webserver
# This allows external access to the Dagster UI
EXPOSE 80
