package dagster

import (
	"github.com/eth-library/dap/cli/internal/exec"
	"github.com/spf13/cobra"
)

// MaterializeCmd materializes all Dagster assets.
var MaterializeCmd = &cobra.Command{
	Use:     "materialize [flags]",
	Short:   "Materialize Dagster assets",
	Long:    "Materialize all Dagster assets. Any additional flags are passed to dagster.",
	GroupID: GroupID,
	// Allow passing flags through to dagster
	DisableFlagParsing: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		dagsterArgs := []string{
			"asset", "materialize",
			"-m", "da_pipeline.definitions",
			"--select", "*",
		}
		dagsterArgs = append(dagsterArgs, args...)
		return exec.RunInteractive("dagster", dagsterArgs...)
	},
}
