package dagster

import (
	"github.com/eth-library/dap/cli/internal/exec"
	"github.com/spf13/cobra"
)

// RunCmd runs the Dagster ingest_sip_job.
var RunCmd = &cobra.Command{
	Use:     "run [flags]",
	Short:   "Run Dagster job",
	Long:    "Run the ingest_sip_job. Any additional flags are passed to dagster.",
	GroupID: GroupID,
	// Allow passing flags through to dagster
	DisableFlagParsing: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		dagsterArgs := []string{
			"job", "execute",
			"-m", "da_pipeline.definitions",
			"-j", "ingest_sip_job",
		}
		dagsterArgs = append(dagsterArgs, args...)
		return exec.RunInteractive("dagster", dagsterArgs...)
	},
}
