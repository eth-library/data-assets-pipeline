package k8s

import (
	"fmt"
	"os"
	"os/exec"

	dapexec "github.com/eth-library/dap/cli/internal/exec"
	"github.com/eth-library/dap/cli/internal/ui"
	"github.com/spf13/cobra"
)

var upCmd = &cobra.Command{
	Use:   "up",
	Short: "Build and deploy to local Kubernetes",
	Long:  "Build the Docker image and deploy to local Kubernetes cluster (localhost:8080).",
	RunE: func(cmd *cobra.Command, args []string) error {
		// Check Kubernetes connectivity
		if err := checkK8s(); err != nil {
			return err
		}

		// Build Docker image
		ui.Info("Building Docker image...")
		if err := dapexec.RunPassthrough("docker", "build", "-t", Image, "-q", "."); err != nil {
			ui.Error("Docker build failed")
			return fmt.Errorf("docker build failed: %w", err)
		}
		ui.Success("Image built", "tag", Image)

		// Create namespace
		ui.Info("Creating namespace...")
		nsCmd := exec.Command("kubectl", "create", "namespace", Namespace, "--dry-run=client", "-o", "yaml")
		applyCmd := exec.Command("kubectl", "apply", "-f", "-")
		applyCmd.Stdin, _ = nsCmd.StdoutPipe()
		applyCmd.Stdout = nil
		applyCmd.Stderr = nil
		nsCmd.Start()
		applyCmd.Run()
		nsCmd.Wait()
		ui.Success("Namespace ready", "name", Namespace)

		// Create PostgreSQL secret if it doesn't exist
		if _, err := dapexec.Run("kubectl", "get", "secret", PGSecretName, "-n", Namespace); err != nil {
			ui.Info("Creating PostgreSQL secret...")
			password, _ := dapexec.Run("openssl", "rand", "-base64", "24")
			if err := dapexec.RunPassthrough("kubectl", "create", "secret", "generic", PGSecretName,
				"--from-literal=postgresql-password="+password,
				"-n", Namespace); err != nil {
				ui.Error("Failed to create PostgreSQL secret")
				return fmt.Errorf("failed to create postgresql secret: %w", err)
			}
			ui.Success("PostgreSQL secret created")
		}

		// Add Helm repo
		ui.Info("Setting up Helm repository...")
		dapexec.Run("helm", "repo", "add", "dagster", "https://dagster-io.github.io/helm")
		dapexec.Run("helm", "repo", "update", "dagster")
		ui.Success("Helm repo ready")

		// Apply PVC
		ui.Info("Applying persistent volume claim...")
		if err := dapexec.RunPassthrough("kubectl", "apply", "-n", Namespace, "-f", "helm/pvc.yaml"); err != nil {
			ui.Error("Failed to apply PVC")
			return fmt.Errorf("failed to apply pvc: %w", err)
		}
		ui.Success("PVC applied")

		// Create test data ConfigMap if test data exists
		if _, err := os.Stat("da_pipeline_tests/test_data"); err == nil {
			ui.Info("Creating test data ConfigMap...")
			configMapCmd := exec.Command("kubectl", "create", "configmap", "test-data-xml",
				"--from-file=da_pipeline_tests/test_data/",
				"-n", Namespace, "--dry-run=client", "-o", "yaml")
			applyCmd := exec.Command("kubectl", "apply", "-f", "-")
			applyCmd.Stdin, _ = configMapCmd.StdoutPipe()
			configMapCmd.Start()
			applyCmd.Run()
			configMapCmd.Wait()
			ui.Success("Test data ConfigMap created")
		}

		// Deploy with Helm
		ui.Info("Deploying with Helm...", "version", HelmVersion)
		if err := dapexec.RunPassthrough("helm", "upgrade", "--install", Release, HelmChart,
			"-f", "helm/values.yaml",
			"-f", "helm/values-local.yaml",
			"-n", Namespace,
			"--version", HelmVersion,
			"--skip-schema-validation"); err != nil {
			ui.Error("Helm deployment failed")
			return fmt.Errorf("helm deployment failed: %w", err)
		}
		ui.Success("Dagster deployed")

		// Show status
		ui.Newline()
		statusCmd.RunE(cmd, args)
		ui.Newline()

		ui.Success("Dagster deployed")
		ui.Info("UI available at", "url", DagsterUIURL)
		ui.Hint("Service exposed via LoadBalancer - no port-forward needed")
		return nil
	},
}

func checkK8s() error {
	if _, err := dapexec.Run("kubectl", "cluster-info"); err != nil {
		ui.Error("Kubernetes not available. Enable it in Docker Desktop.")
		return fmt.Errorf("kubernetes not available")
	}
	dapexec.Run("kubectl", "config", "use-context", K8sContext)
	ui.Success("Kubernetes cluster connected")
	return nil
}

func init() {
	K8sCmd.AddCommand(upCmd)
}
