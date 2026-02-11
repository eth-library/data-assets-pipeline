package k8s

import (
	"github.com/eth-library/dap/cli/internal/exec"
	"github.com/spf13/cobra"
)

var logsCmd = &cobra.Command{
	Use:   "logs [flags]",
	Short: "Stream logs from user code pod",
	Long:  "Stream logs from the Dagster user code pod. Additional flags are passed to kubectl.",
	// Allow passing flags through to kubectl
	DisableFlagParsing: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		kubectlArgs := []string{
			"logs", "-n", Namespace,
			"-l", "app.kubernetes.io/name=dagster-user-deployments",
			"--tail=100", "-f",
		}
		kubectlArgs = append(kubectlArgs, args...)
		return exec.RunInteractive("kubectl", kubectlArgs...)
	},
}

func init() {
	K8sCmd.AddCommand(logsCmd)
}
