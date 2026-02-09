#!/usr/bin/env bash
# Setup Python environment using uv
set -euo pipefail

if ! command -v uv &>/dev/null; then
    echo "error: uv is required to manage Python dependencies" >&2
    echo "       install it from https://docs.astral.sh/uv/" >&2
    exit 1
fi

uv sync --extra dev --quiet
