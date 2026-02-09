package k8s

import (
	"fmt"

	"github.com/eth-library/dap/cli/internal/exec"
	"github.com/eth-library/dap/cli/internal/ui"
	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show pods and services",
	Long:  "Display the current status of Kubernetes pods and services.",
	RunE: func(cmd *cobra.Command, args []string) error {
		// Check if namespace exists
		if _, err := exec.Run("kubectl", "get", "namespace", Namespace); err != nil {
			ui.Warn("Deployment not running")
			ui.CommandHint("dap k8s up", "deploy to Kubernetes")
			return nil
		}

		// Show pods
		ui.Section("Pods")
		pods, err := exec.Run("kubectl", "get", "pods", "-n", Namespace, "-o", "wide")
		if err != nil || pods == "" || pods == "No resources found in "+Namespace+" namespace." {
			ui.KeyValueDim("status", "no pods running")
		} else {
			fmt.Println(pods)
		}

		// Show services
		ui.Section("Services")
		svcs, err := exec.Run("kubectl", "get", "svc", "-n", Namespace)
		if err != nil || svcs == "" {
			ui.KeyValueDim("status", "no services running")
		} else {
			fmt.Println(svcs)
		}

		return nil
	},
}

func init() {
	K8sCmd.AddCommand(statusCmd)
}
