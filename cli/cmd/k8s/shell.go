package k8s

import (
	"fmt"
	"strings"

	"github.com/eth-library/dap/cli/internal/exec"
	"github.com/eth-library/dap/cli/internal/ui"
	"github.com/spf13/cobra"
)

var shellCmd = &cobra.Command{
	Use:   "shell",
	Short: "Open shell in user code pod",
	Long:  "Open an interactive bash shell in the Dagster user code pod.",
	RunE: func(cmd *cobra.Command, args []string) error {
		// Get pod name
		podName, err := exec.Run("kubectl", "get", "pods", "-n", Namespace,
			"-l", "app.kubernetes.io/name=dagster-user-deployments",
			"-o", "jsonpath={.items[0].metadata.name}")
		if err != nil || podName == "" {
			ui.Error("No user code pod found. Is Dagster deployed?")
			return fmt.Errorf("failed to find user code pod: %w", err)
		}

		podName = strings.TrimSpace(podName)
		ui.Info("Connecting to pod", "name", podName)

		return exec.RunInteractive("kubectl", "exec", "-it", "-n", Namespace, podName, "--", "/bin/bash")
	},
}

func init() {
	K8sCmd.AddCommand(shellCmd)
}
