package k8s

import (
	"fmt"

	"github.com/eth-library/dap/cli/internal/exec"
	"github.com/eth-library/dap/cli/internal/ui"
	"github.com/spf13/cobra"
)

var restartCmd = &cobra.Command{
	Use:   "restart",
	Short: "Rebuild and restart user code pod",
	Long:  "Rebuild the Docker image and restart the user code deployment.",
	RunE: func(cmd *cobra.Command, args []string) error {
		// Check Kubernetes connectivity
		if err := checkK8s(); err != nil {
			return err
		}

		// Build Docker image
		ui.Info("Building Docker image...")
		if err := exec.RunPassthrough("docker", "build", "-t", Image, "-q", "."); err != nil {
			ui.Error("Docker build failed")
			return fmt.Errorf("docker build failed: %w", err)
		}
		ui.Success("Image built", "tag", Image)

		// Restart deployment
		ui.Info("Restarting user code deployment...")
		if err := exec.RunPassthrough("kubectl", "rollout", "restart", "deployment",
			"-n", Namespace, "-l", "app.kubernetes.io/name=dagster-user-deployments"); err != nil {
			ui.Error("Restart failed")
			return fmt.Errorf("deployment restart failed: %w", err)
		}

		// Wait for rollout
		ui.Info("Waiting for rollout...")
		if err := exec.RunPassthrough("kubectl", "rollout", "status", "deployment",
			"-n", Namespace, "-l", "app.kubernetes.io/name=dagster-user-deployments", "--timeout="+RolloutTimeout); err != nil {
			ui.Error("Rollout failed")
			return fmt.Errorf("rollout status check failed: %w", err)
		}

		ui.Success("Restart complete")
		return nil
	},
}

func init() {
	K8sCmd.AddCommand(restartCmd)
}
