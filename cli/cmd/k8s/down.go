package k8s

import (
	"github.com/eth-library/dap/cli/internal/exec"
	"github.com/eth-library/dap/cli/internal/ui"
	"github.com/spf13/cobra"
)

var downCmd = &cobra.Command{
	Use:   "down",
	Short: "Tear down Kubernetes deployment",
	Long:  "Remove the Dagster deployment from Kubernetes.",
	RunE: func(cmd *cobra.Command, args []string) error {
		ui.Info("Tearing down Dagster deployment...")

		// Uninstall Helm release
		exec.Run("helm", "uninstall", Release, "-n", Namespace, "--wait=false")
		ui.Success("Helm release uninstalled")

		// Clean up jobs
		exec.Run("kubectl", "delete", "jobs", "-n", Namespace, "-l", "dagster/run-id", "--timeout=10s")

		// Clean up pods
		exec.Run("kubectl", "delete", "pods", "-n", Namespace, "-l", "dagster/run-id",
			"--grace-period=0", "--force", "--timeout=10s")

		// Clean up PVC
		exec.Run("kubectl", "delete", "pvc", "dagster-storage", "-n", Namespace, "--timeout=10s")

		// Clean up ConfigMap
		exec.Run("kubectl", "delete", "configmap", "test-data-xml", "-n", Namespace, "--timeout=10s")

		ui.Success("Teardown complete")
		return nil
	},
}

func init() {
	K8sCmd.AddCommand(downCmd)
}
