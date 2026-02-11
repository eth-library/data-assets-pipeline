package dev

import (
	"fmt"

	"github.com/eth-library/dap/cli/internal/exec"
	"github.com/eth-library/dap/cli/internal/ui"
	"github.com/spf13/cobra"
)

// TypecheckCmd runs mypy type checking.
var TypecheckCmd = &cobra.Command{
	Use:     "typecheck",
	Short:   "Run type checking",
	Long:    "Run mypy to check Python type annotations.",
	GroupID: GroupID,
	RunE: func(cmd *cobra.Command, args []string) error {
		ui.Info("Type checking with mypy...")
		if err := exec.RunPassthrough("mypy", PythonTargets...); err != nil {
			ui.Error("Type check failed")
			return fmt.Errorf("mypy type check failed: %w", err)
		}
		ui.Success("No type errors")
		return nil
	},
}
