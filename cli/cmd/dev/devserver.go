package dev

import (
	"github.com/eth-library/dap/cli/internal/exec"
	"github.com/eth-library/dap/cli/internal/ui"
	"github.com/spf13/cobra"
)

// DevServerCmd starts the Dagster development server.
var DevServerCmd = &cobra.Command{
	Use:   "dev",
	Short: "Start Dagster development server",
	Long: `Start the Dagster development server with hot-reload.

Runs 'dagster dev' which starts the webserver at http://localhost:3000.
Code changes are automatically reloaded without server restart.

Press Ctrl+C to stop the server.`,
	GroupID: GroupID,
	RunE: func(cmd *cobra.Command, args []string) error {
		ui.TaskStart("Starting Dagster dev server...")
		ui.KeyValue("url", "http://localhost:3000")
		ui.Newline()
		return exec.RunInteractive("dagster", "dev")
	},
}
