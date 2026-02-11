// Package k8s contains Kubernetes-related commands.
package k8s

import (
	"github.com/spf13/cobra"
)

// Configuration for Kubernetes deployment
const (
	Namespace    = "dagster"
	Release      = "dagster"
	Image        = "da-pipeline:local"
	HelmChart    = "dagster/dagster"
	HelmVersion  = "1.10.14"
	PGSecretName = "dagster-postgresql"

	// Network configuration
	DagsterUIURL = "http://localhost:8080"
	K8sContext   = "docker-desktop"

	// Timeouts
	RolloutTimeout = "120s"
)

// K8sCmd is the parent command for Kubernetes operations.
var K8sCmd = &cobra.Command{
	Use:   "k8s",
	Short: "Kubernetes operations",
	Long:  "Commands for deploying and managing the pipeline on Kubernetes.",
}

func init() {
	// Subcommands are added in their respective files
}
