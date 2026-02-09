# Contributing to the DAP CLI

This guide is for developers who want to modify or extend the `dap` CLI tool.

## Quick Start

```bash
# Rebuild the CLI (go mod tidy + gomod2nix + nix build)
dap cli build

# Run tests
dap cli test
```

## Prerequisites

- Go 1.22+ (provided by `nix develop`)
- Understanding of [Cobra](https://github.com/spf13/cobra) CLI framework
- Familiarity with [Lipgloss](https://github.com/charmbracelet/lipgloss) for terminal styling

## Project Structure

```
.cli/
├── main.go                 # Entry point
├── go.mod                  # Go module definition
├── go.sum                  # Dependency checksums
├── gomod2nix.toml          # Nix dependency hashes (auto-generated)
├── package.nix             # Nix build definition
├── bin/dap                 # Development binary
│
├── cmd/                    # Command implementations
│   ├── root.go             # Root command, groups, global flags
│   │
│   ├── dev/                # Development commands
│   │   ├── dev.go          # Package entry, PythonTargets constant
│   │   ├── check.go        # dap check
│   │   ├── devserver.go    # dap dev
│   │   ├── fmt.go          # dap fmt
│   │   ├── lint.go         # dap lint
│   │   ├── test.go         # dap test
│   │   └── typecheck.go    # dap typecheck
│   │
│   ├── env/                # Environment commands
│   │   ├── env.go          # Package entry
│   │   ├── clean.go        # dap clean
│   │   ├── reset.go        # dap reset
│   │   ├── versions.go     # dap versions
│   │   └── welcome.go      # dap welcome
│   │
│   ├── dagster/            # Dagster commands
│   │   ├── dagster.go      # Package entry
│   │   ├── materialize.go  # dap materialize
│   │   └── run.go          # dap run
│   │
│   ├── k8s/                # Kubernetes commands
│   │   ├── k8s.go          # Parent command, constants
│   │   ├── up.go           # dap k8s up
│   │   ├── down.go         # dap k8s down
│   │   ├── restart.go      # dap k8s restart
│   │   ├── status.go       # dap k8s status
│   │   ├── logs.go         # dap k8s logs
│   │   └── shell.go        # dap k8s shell
│   │
│   └── meta/               # Hidden CLI maintenance commands
│       ├── meta.go         # Package entry
│       ├── build.go        # dap build (full rebuild)
│       └── go.go           # dap go build, dap go test
│
└── internal/               # Internal packages
    ├── exec/               # Command execution utilities
    │   ├── run.go          # Run, RunInteractive, RunPassthrough, Which
    │   └── run_test.go
    └── ui/                 # Terminal UI components
        ├── styles.go       # Colors, styles, symbols
        ├── output.go       # Banner, Section, Success, Error, etc.
        └── styles_test.go
```

## Adding a New Command

### 1. Choose the package

Commands are organized by domain:

| Package | Group | Commands |
|---------|-------|----------|
| `cmd/dev` | Development | check, dev, fmt, lint, test, typecheck |
| `cmd/env` | Environment | clean, reset, versions, welcome |
| `cmd/dagster` | Dagster | materialize, run |
| `cmd/k8s` | Kubernetes | up, down, status, restart, logs, shell |
| `cmd/meta` | Hidden | build, go |

### 2. Create the command file

Create a new file in the appropriate package (e.g., `cmd/dev/mycommand.go`):

```go
package dev

import (
    "fmt"

    "github.com/eth-library/dap/cli/internal/exec"
    "github.com/eth-library/dap/cli/internal/ui"
    "github.com/spf13/cobra"
)

var MyCmd = &cobra.Command{
    Use:     "mycommand",
    Short:   "Brief description",
    Long:    "Detailed description of what this command does.",
    GroupID: GroupID,
    RunE: func(cmd *cobra.Command, args []string) error {
        ui.Info("Doing something...")

        if err := exec.RunPassthrough("some-tool", "arg1"); err != nil {
            ui.Error("Something failed")
            return fmt.Errorf("some-tool failed: %w", err)
        }

        ui.Success("Done!")
        return nil
    },
}
```

### 3. Register the command

Add to the `Commands()` function in the package's main file (e.g., `cmd/dev/dev.go`):

```go
func Commands() []*cobra.Command {
    return []*cobra.Command{
        CheckCmd,
        DevServerCmd,
        // ...
        MyCmd,  // Add here
    }
}
```

### 4. Add flags (optional)

```go
func init() {
    MyCmd.Flags().BoolP("verbose", "v", false, "Verbose output")
}
```

## UI Guidelines

Use the `internal/ui` package for all terminal output:

```go
import "github.com/eth-library/dap/cli/internal/ui"

// Status messages
ui.Info("Starting...")
ui.Success("Completed", "key", value)
ui.Warn("Something might be wrong")
ui.Error("Failed")

// Sections and structure
ui.Section("Section Title")
ui.Divider()
ui.Newline()

// Key-value pairs
ui.KeyValue("name", "value")
ui.KeyValueStatus("name", "value", true)  // with checkmark
ui.KeyValueDim("name", "not set")         // dimmed

// Multi-step operations
ui.Step(1, 3, "First step...")
ui.StepDone(1, 3, "First step complete")
ui.StepFail(1, 3, "First step failed")

// Command hints
ui.CommandHint("dap dev", "start the dev server")
```

Colors auto-disable in CI environments and with `--no-color`.

## Running External Commands

Use `internal/exec` instead of `os/exec` directly:

```go
import "github.com/eth-library/dap/cli/internal/exec"

// Capture output
output, err := exec.Run("command", "arg1", "arg2")

// Pass through to terminal (for interactive commands)
err := exec.RunInteractive("command", "arg1")

// Pass through stdout/stderr only
err := exec.RunPassthrough("command", "arg1")

// Check if command exists
if exec.Which("docker") {
    // docker is available
}
```

## Error Handling

Always wrap errors with context:

```go
if err := exec.RunPassthrough("ruff", "check", "."); err != nil {
    ui.Error("Lint check failed")
    return fmt.Errorf("ruff check failed: %w", err)
}
```

## Testing

```bash
dap go test       # Run tests
dap go test -v    # Verbose output
```

## Building

```bash
dap go build      # Quick build to .cli/bin/dap
dap build         # Full rebuild (go mod tidy + gomod2nix + nix build)
```

## Dependency Management

After changing Go dependencies:

```bash
dap build --deps-only   # go mod tidy + gomod2nix
```

## Shared Constants

### Python targets

Use `PythonTargets` from `cmd/dev/dev.go`:

```go
exec.RunPassthrough("ruff", append([]string{"check"}, dev.PythonTargets...)...)
```

### Kubernetes constants

Use constants from `cmd/k8s/k8s.go`:

```go
const (
    Namespace    = "dagster"
    DagsterUIURL = "http://localhost:8080"
    K8sContext   = "docker-desktop"
)
```

## Code Style

- Follow standard Go conventions (`go fmt`, `go vet`)
- Keep functions small and focused
- Export commands as `FooCmd` (e.g., `CheckCmd`, `DevServerCmd`)
- Use `GroupID` constant from the package for command grouping

## Release Checklist

1. `dap go test` - Run tests
2. `dap go build` - Build
3. `./bin/dap --help` - Test manually
4. `dap build` - Full Nix build
