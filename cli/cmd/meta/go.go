package meta

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/eth-library/dap/cli/internal/exec"
	"github.com/eth-library/dap/cli/internal/ui"
	"github.com/spf13/cobra"
)

// CliCmd is the parent command for CLI development operations.
var CliCmd = &cobra.Command{
	Use:   "cli",
	Short: "CLI development commands",
	Long:  "Commands for building and testing the dap CLI.",
}

// CliBuildCmd rebuilds the dap CLI with Nix.
var CliBuildCmd = &cobra.Command{
	Use:   "build",
	Short: "Rebuild the dap CLI",
	Long:  "Runs go mod tidy, gomod2nix, and nix build to rebuild the CLI.",
	RunE: func(cmd *cobra.Command, args []string) error {
		// Find cli directory
		cliDir, err := findCliDir()
		if err != nil {
			return err
		}

		// Change to cli directory
		origDir, _ := os.Getwd()
		if err := os.Chdir(cliDir); err != nil {
			return err
		}
		defer os.Chdir(origDir)

		// Step 1: go mod tidy
		ui.Info("Running go mod tidy...")
		if err := exec.RunPassthrough("go", "mod", "tidy"); err != nil {
			ui.Error("go mod tidy failed", "error", err)
			return err
		}
		ui.Success("go mod tidy complete")

		// Step 2: gomod2nix
		if !exec.Which("gomod2nix") {
			ui.Warn("gomod2nix not found in PATH - skipping")
		} else {
			ui.Info("Running gomod2nix...")
			if err := exec.RunPassthrough("gomod2nix"); err != nil {
				ui.Error("gomod2nix failed", "error", err)
				return err
			}
			ui.Success("gomod2nix.toml updated")
		}

		// Step 3: nix build
		os.Chdir(origDir) // Back to repo root for nix build
		ui.Info("Running nix build...")
		if err := exec.RunPassthrough("nix", "build", ".#dap"); err != nil {
			ui.Error("nix build failed", "error", err)
			return err
		}
		ui.Success("Build complete")
		ui.Newline()
		ui.CommandHint("direnv reload", "to use the new binary")

		return nil
	},
}

// CliTestCmd runs tests for the Go CLI.
var CliTestCmd = &cobra.Command{
	Use:   "test",
	Short: "Run Go CLI tests",
	Long:  "Runs go test for all CLI packages.",
	RunE: func(cmd *cobra.Command, args []string) error {
		verbose, _ := cmd.Flags().GetBool("verbose")

		// Find cli directory
		cliDir, err := findCliDir()
		if err != nil {
			return err
		}

		// Change to cli directory
		origDir, _ := os.Getwd()
		if err := os.Chdir(cliDir); err != nil {
			return err
		}
		defer os.Chdir(origDir)

		ui.Info("Running CLI tests...")
		testArgs := []string{"test", "./..."}
		if verbose {
			testArgs = append(testArgs, "-v")
		}
		if err := exec.RunPassthrough("go", testArgs...); err != nil {
			ui.Error("Tests failed", "error", err)
			return err
		}
		ui.Success("All tests passed")
		return nil
	},
}

// CliLintCmd runs linting and format checks for the Go CLI.
var CliLintCmd = &cobra.Command{
	Use:   "lint",
	Short: "Lint and check formatting",
	Long:  "Runs go vet and checks gofmt formatting for all CLI packages.",
	RunE: func(cmd *cobra.Command, args []string) error {
		// Find cli directory
		cliDir, err := findCliDir()
		if err != nil {
			return err
		}

		// Change to cli directory
		origDir, _ := os.Getwd()
		if err := os.Chdir(cliDir); err != nil {
			return err
		}
		defer os.Chdir(origDir)

		ui.Info("Running go vet...")
		if err := exec.RunPassthrough("go", "vet", "./..."); err != nil {
			ui.Error("go vet failed", "error", err)
			return err
		}
		ui.Success("go vet passed")

		ui.Info("Checking formatting...")
		output, err := exec.Run("gofmt", "-l", ".")
		if err != nil {
			ui.Error("gofmt failed", "error", err)
			return err
		}
		if output != "" {
			ui.Error("Files not formatted:")
			for _, f := range strings.Split(strings.TrimSpace(output), "\n") {
				ui.ListItem(f)
			}
			return fmt.Errorf("gofmt check failed")
		}
		ui.Success("Formatting check passed")

		return nil
	},
}

// findCliDir returns the absolute path to the cli directory using the git repo root.
func findCliDir() (string, error) {
	root, err := exec.Run("git", "rev-parse", "--show-toplevel")
	if err != nil {
		return "", fmt.Errorf("not in a git repository: %w", err)
	}
	return filepath.Join(strings.TrimSpace(root), "cli"), nil
}

func init() {
	CliTestCmd.Flags().BoolP("verbose", "v", false, "Verbose test output")
	CliCmd.AddCommand(CliBuildCmd)
	CliCmd.AddCommand(CliTestCmd)
	CliCmd.AddCommand(CliLintCmd)
}
