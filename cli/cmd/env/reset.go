package env

import (
	"fmt"

	"github.com/eth-library/dap/cli/internal/exec"
	"github.com/eth-library/dap/cli/internal/ui"
	"github.com/spf13/cobra"
)

// ResetCmd cleans and reinstalls dependencies.
var ResetCmd = &cobra.Command{
	Use:     "reset",
	Short:   "Clean and reinstall dependencies",
	Long:    "Remove the virtual environment, caches, and reinstall all dependencies.",
	GroupID: GroupID,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Run clean
		ui.Info("Cleaning environment...")
		if err := CleanCmd.RunE(cmd, args); err != nil {
			return fmt.Errorf("clean failed: %w", err)
		}

		ui.Newline()

		// Reinstall dependencies
		ui.Info("Reinstalling dependencies...")
		if err := exec.RunPassthrough("uv", "sync", "--extra", "dev"); err != nil {
			ui.Error("Failed to sync dependencies")
			return fmt.Errorf("uv sync failed: %w", err)
		}

		ui.Success("Environment reset complete")
		return nil
	},
}
