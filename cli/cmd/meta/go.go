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
		if err := inCliDir(func() error {
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

			return nil
		}); err != nil {
			return err
		}

		// Step 3: nix build (must run from repo root where flake.nix lives)
		root, err := repoRoot()
		if err != nil {
			return err
		}
		if err := os.Chdir(root); err != nil {
			return err
		}
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

		return inCliDir(func() error {
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
		})
	},
}

var cliLintFix bool

// CliLintCmd runs linting and format checks for the Go CLI.
var CliLintCmd = &cobra.Command{
	Use:   "lint",
	Short: "Lint and check formatting",
	Long:  "Runs go vet and checks gofmt formatting for all CLI packages. Use --fix to auto-format.",
	RunE: func(cmd *cobra.Command, args []string) error {
		return inCliDir(func() error {
			if cliLintFix {
				return lintFixCli()
			}
			return lintCheckCli()
		})
	},
}

// lintFixCli formats code and runs go vet.
func lintFixCli() error {
	ui.Info("Fixing formatting...")
	if err := exec.RunPassthrough("gofmt", "-w", "."); err != nil {
		ui.Error("gofmt -w failed", "error", err)
		return err
	}
	ui.Success("Formatting fixed")

	return goVet()
}

// lintCheckCli runs go vet and checks formatting without modifying files.
func lintCheckCli() error {
	if err := goVet(); err != nil {
		return err
	}

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
}

// goVet runs go vet on all packages.
func goVet() error {
	ui.Info("Running go vet...")
	if err := exec.RunPassthrough("go", "vet", "./..."); err != nil {
		ui.Error("go vet failed", "error", err)
		return err
	}
	ui.Success("go vet passed")
	return nil
}

// inCliDir runs fn inside the cli directory, restoring the original directory afterwards.
func inCliDir(fn func() error) error {
	cliDir, err := findCliDir()
	if err != nil {
		return err
	}
	origDir, _ := os.Getwd()
	if err := os.Chdir(cliDir); err != nil {
		return err
	}
	defer os.Chdir(origDir)
	return fn()
}

// repoRoot returns the absolute path to the git repository root.
func repoRoot() (string, error) {
	root, err := exec.Run("git", "rev-parse", "--show-toplevel")
	if err != nil {
		return "", fmt.Errorf("not in a git repository: %w", err)
	}
	return strings.TrimSpace(root), nil
}

// findCliDir returns the absolute path to the cli directory using the git repo root.
func findCliDir() (string, error) {
	root, err := repoRoot()
	if err != nil {
		return "", err
	}
	return filepath.Join(root, "cli"), nil
}

func init() {
	CliTestCmd.Flags().BoolP("verbose", "v", false, "Verbose test output")
	CliLintCmd.Flags().BoolVar(&cliLintFix, "fix", false, "Auto-fix formatting")
	CliCmd.AddCommand(CliBuildCmd)
	CliCmd.AddCommand(CliTestCmd)
	CliCmd.AddCommand(CliLintCmd)
}
