package env

import (
	"os"

	"github.com/eth-library/dap/cli/internal/ui"
	"github.com/spf13/cobra"
)

// WelcomeCmd shows the welcome banner and environment info.
var WelcomeCmd = &cobra.Command{
	Use:     "welcome",
	Short:   "Show welcome banner and environment info",
	Long:    "Display the project banner, environment versions, and quick command hints.",
	GroupID: GroupID,
	Run: func(cmd *cobra.Command, args []string) {
		showWelcome()
	},
}

func showWelcome() {
	// Banner
	ui.Banner("Data Archive Pipeline (DAP) - Orchestrator", "ETH Library Zurich")

	// Versions section
	ui.Section("Versions")
	ShowVersionsCompact()

	// Environment section with status indicators
	ui.Section("Environment")
	showEnvironmentPaths()

	// Command hints (unless DAP_QUIET is set)
	if os.Getenv("DAP_QUIET") == "" {
		ui.Section("Quick Start")
		ui.CommandHint("dap dev", "Start the Dagster development server")
		ui.CommandHint("dap test", "Run tests (pytest)")
		ui.CommandHint("dap check", "Run all quality checks (ruff, mypy, pytest)")
		ui.Newline()
		ui.Hint("Run 'dap --help' for all commands")
	}

	ui.Newline()
}

// showEnvironmentPaths displays important environment paths with status.
func showEnvironmentPaths() {
	// Nix flake - check if flake.nix exists
	if cwd, err := os.Getwd(); err == nil {
		_, flakeErr := os.Stat(cwd + "/flake.nix")
		ui.KeyValueStatus("nix flake", cwd, flakeErr == nil)
	}

	// Python venv path - check if it exists
	if venv := os.Getenv("VIRTUAL_ENV"); venv != "" {
		_, venvErr := os.Stat(venv)
		ui.KeyValueStatus("python venv", venv, venvErr == nil)
	} else {
		ui.KeyValueDim("python venv", "not set")
	}

	// DAGSTER_HOME - check if directory exists
	if dagsterHome := os.Getenv("DAGSTER_HOME"); dagsterHome != "" {
		_, homeErr := os.Stat(dagsterHome)
		ui.KeyValueStatus("DAGSTER_HOME", dagsterHome, homeErr == nil)
	} else {
		ui.KeyValueDim("DAGSTER_HOME", "not set")
	}
}
