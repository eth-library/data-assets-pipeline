{
  description = "Development environment for da_pipeline using Nix, uv, and .venv";

  inputs.nixpkgs.url = "github:NixOS/nixpkgs/nixos-25.11";

  outputs = { self, nixpkgs, flake-utils, ... }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        # Import nixpkgs for the current system
        pkgs = import nixpkgs {
          inherit system;
          # No overlays needed for this project
          overlays = [];
        };

        # Define the path to the Python virtual environment
        venvDir = "./.venv";

        # Define functions for structured logging
        log = {
          info = message: "echo -e \"\\033[0;34m[INFO]\\033[0m ${message}\"";
          success = message: "echo -e \"\\033[0;32m[SUCCESS]\\033[0m ${message}\"";
          warning = message: "echo -e \"\\033[0;33m[WARNING]\\033[0m ${message}\"";
          error = message: "echo -e \"\\033[0;31m[ERROR]\\033[0m ${message}\"";
        };
      in
      {
        devShell = pkgs.mkShell {
          # Define packages that should be available in the development shell
          buildInputs = [
            # Python 3.12 and core build tools
            pkgs.python312
            pkgs.python312Packages.hatchling
            pkgs.uv  # Modern Python package manager and installer
          ];

          shellHook = ''
            ${log.info "Initializing Python environment..."}

            # Create virtual environment if it doesn't exist
            if [ ! -d ${venvDir} ]; then
              ${log.info "Creating virtual environment in .venv..."}
              uv venv --python "$(which python)"
            fi

            # Activate virtual environment first to ensure we're using the right Python
            ${log.info "Activating virtual environment..."}
            # shellcheck disable=SC1091
            source ${venvDir}/bin/activate

            # Add the virtual environment's bin directory to PATH
            # This ensures that Dagster tools and uv commands are available
            export PATH="${venvDir}/bin:$PATH"

            # Generate lock file if it doesn't exist
            if [ ! -f uv.lock ]; then
              ${log.info "Generating lock file..."}
              uv lock
            fi

            # Install dependencies including development tools
            ${log.info "Installing Python dependencies with uv..."}
            uv sync --extra dev

            # Print helpful information
            echo ""
            ${log.success "Development environment activated!"}
            echo "Python version: $(python --version)"
            echo "Python interpreter path: $(python -c 'import sys; print(sys.executable)')"
            echo "uv version: $(uv --version)"
            echo ""
            echo "Available commands:"
            echo "  - uv sync: Update dependencies"
            echo "  - uv lock: Update lock file"
            echo "  - dagster dev: Start Dagster development server"
            echo ""
          '';
        };
      }
    );
}