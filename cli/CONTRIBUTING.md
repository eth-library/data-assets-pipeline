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
cli/
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
dap cli test       # Run tests
dap cli test -v    # Verbose output
```

## Linting

```bash
dap cli lint          # Check go vet and formatting
dap cli lint --fix    # Auto-fix formatting, then run go vet
```

## Building

```bash
dap cli build     # Full rebuild (go mod tidy + gomod2nix + nix build)
```

`dap cli build` runs three steps:

1. **`go mod tidy`** — syncs `go.mod` and `go.sum` with the actual imports in the code
2. **`gomod2nix`** — regenerates `gomod2nix.toml` from `go.sum`, translating Go dependency hashes into a format Nix can use (skipped if not installed)
3. **`nix build .#dap`** — builds the binary in a hermetic Nix sandbox

The Nix build chain works as follows:

- The root `flake.nix` delegates to `cli/flake.nix`, which loads the [gomod2nix](https://github.com/nix-community/gomod2nix) overlay
- `cli/package.nix` defines the build using `buildGoApplication`, which compiles the Go source with pinned dependencies from `gomod2nix.toml`
- Linker flags (`-s -w`) strip debug info for a smaller binary, and `-X ...cmd.version` embeds the version string
- The output binary is renamed from `cli` to `dap` and placed at `result/bin/dap`

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

1. `dap cli test` - Run tests
2. `dap cli lint` - Check formatting and vet
3. `dap cli build` - Full rebuild
4. `./bin/dap --help` - Test manually
