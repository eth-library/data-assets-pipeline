package dev

import (
	"fmt"

	"github.com/eth-library/dap/cli/internal/exec"
	"github.com/eth-library/dap/cli/internal/ui"
	"github.com/spf13/cobra"
)

var lintFix bool

// LintCmd checks code style and formatting.
var LintCmd = &cobra.Command{
	Use:     "lint",
	Short:   "Check code style and formatting",
	Long:    "Run ruff to check code style and formatting. Use --fix to auto-fix issues.",
	GroupID: GroupID,
	RunE: func(cmd *cobra.Command, args []string) error {
		if lintFix {
			// Fix mode
			ui.Info("Fixing lint issues...")
			if err := exec.RunPassthrough("ruff", append([]string{"check", "--fix"}, PythonTargets...)...); err != nil {
				return fmt.Errorf("ruff check --fix failed: %w", err)
			}
			ui.Info("Formatting code...")
			if err := exec.RunPassthrough("ruff", append([]string{"format"}, PythonTargets...)...); err != nil {
				return fmt.Errorf("ruff format failed: %w", err)
			}
			ui.Success("Code fixed and formatted")
		} else {
			// Check mode
			ui.Info("Checking code style...")
			if err := exec.RunPassthrough("ruff", append([]string{"check"}, PythonTargets...)...); err != nil {
				ui.Error("Lint check failed")
				return fmt.Errorf("ruff check failed: %w", err)
			}
			ui.Info("Checking formatting...")
			if err := exec.RunPassthrough("ruff", append([]string{"format", "--check"}, PythonTargets...)...); err != nil {
				ui.Error("Format check failed")
				return fmt.Errorf("ruff format check failed: %w", err)
			}
			ui.Success("All lint checks passed")
		}
		return nil
	},
}

func init() {
	LintCmd.Flags().BoolVar(&lintFix, "fix", false, "Auto-fix issues")
}
