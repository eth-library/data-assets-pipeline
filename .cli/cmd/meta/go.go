package meta

import (
	"os"
	"path/filepath"

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
		// Find .cli directory
		cliDir, err := findCliDir()
		if err != nil {
			return err
		}

		// Change to .cli directory
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

		// Find .cli directory
		cliDir, err := findCliDir()
		if err != nil {
			return err
		}

		// Change to .cli directory
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

// findCliDir locates the .cli directory relative to the current working directory or repo root.
func findCliDir() (string, error) {
	// Check if we're already in .cli
	if _, err := os.Stat("go.mod"); err == nil {
		if _, err := os.Stat("gomod2nix.toml"); err == nil {
			return ".", nil
		}
	}

	// Check for .cli in current directory
	if _, err := os.Stat(".cli/go.mod"); err == nil {
		return ".cli", nil
	}

	// Walk up to find repo root with .cli
	dir, _ := os.Getwd()
	for {
		cliPath := filepath.Join(dir, ".cli")
		if _, err := os.Stat(filepath.Join(cliPath, "go.mod")); err == nil {
			return cliPath, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}

	return "", os.ErrNotExist
}

func init() {
	CliTestCmd.Flags().BoolP("verbose", "v", false, "Verbose test output")
	CliCmd.AddCommand(CliBuildCmd)
	CliCmd.AddCommand(CliTestCmd)
}
