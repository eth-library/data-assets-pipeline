// Package cmd contains all CLI commands for dap.
package cmd

import (
	"github.com/eth-library/dap/cli/cmd/dagster"
	"github.com/eth-library/dap/cli/cmd/dev"
	"github.com/eth-library/dap/cli/cmd/env"
	"github.com/eth-library/dap/cli/cmd/k8s"
	"github.com/eth-library/dap/cli/cmd/meta"
	"github.com/eth-library/dap/cli/internal/ui"
	"github.com/spf13/cobra"
)

// Command group IDs for organized --help output
const (
	GroupDevelopment = "development"
	GroupEnvironment = "environment"
	GroupDagster     = "dagster"
	GroupKubernetes  = "kubernetes"
	GroupGoCLI       = "gocli"
)

var rootCmd = &cobra.Command{
	Use:   "dap",
	Short: "Developer tools for the Data Archive Pipeline (DAP) Orchestrator",
	Long: `dap is the CLI for the Data Archive Pipeline (DAP) Orchestrator.

A Dagster-based orchestrator for processing digital assets following
the OAIS reference model. This tool provides commands for local development,
testing, code quality, and Kubernetes deployment.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if noColor, _ := cmd.Flags().GetBool("no-color"); noColor {
			ui.DisableColors()
		}
	},
}

func init() {
	// Disable default completion command
	rootCmd.CompletionOptions.DisableDefaultCmd = true

	// Define command groups for organized help output
	rootCmd.AddGroup(
		&cobra.Group{ID: GroupDevelopment, Title: "Development:"},
		&cobra.Group{ID: GroupEnvironment, Title: "Environment:"},
		&cobra.Group{ID: GroupDagster, Title: "Dagster:"},
		&cobra.Group{ID: GroupKubernetes, Title: "Kubernetes:"},
		&cobra.Group{ID: GroupGoCLI, Title: "CLI Development:"},
	)

	// Global flags
	rootCmd.PersistentFlags().Bool("no-color", false, "Disable colored output")

	// Register commands from dev package
	for _, cmd := range dev.Commands() {
		rootCmd.AddCommand(cmd)
	}

	// Register commands from env package
	for _, cmd := range env.Commands() {
		rootCmd.AddCommand(cmd)
	}

	// Register commands from dagster package
	for _, cmd := range dagster.Commands() {
		rootCmd.AddCommand(cmd)
	}

	// Add k8s subcommand (nested, not flat)
	k8s.K8sCmd.GroupID = GroupKubernetes
	rootCmd.AddCommand(k8s.K8sCmd)

	// Register meta/maintenance commands
	for _, cmd := range meta.Commands() {
		rootCmd.AddCommand(cmd)
	}

	// Add go subcommand for CLI development
	meta.CliCmd.GroupID = GroupGoCLI
	rootCmd.AddCommand(meta.CliCmd)
}

// Execute runs the root command.
func Execute() error {
	return rootCmd.Execute()
}
