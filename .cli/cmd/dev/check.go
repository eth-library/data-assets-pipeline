package dev

import (
	"fmt"

	"github.com/eth-library/dap/cli/internal/exec"
	"github.com/eth-library/dap/cli/internal/ui"
	"github.com/spf13/cobra"
)

// CheckCmd runs all quality checks (lint, typecheck, test).
var CheckCmd = &cobra.Command{
	Use:   "check",
	Short: "Run all quality checks",
	Long: `Run all quality checks in sequence. Fails fast on first error.

Steps:
  1. Lint      ruff check + ruff format --check
  2. Typecheck mypy da_pipeline
  3. Test      pytest da_pipeline_tests`,
	GroupID: GroupID,
	RunE: func(cmd *cobra.Command, args []string) error {
		ui.Section("Quality Checks")

		totalSteps := 3

		// Step 1: Lint
		ui.Step(1, totalSteps, "Checking code style...")
		if err := exec.RunPassthrough("ruff", append([]string{"check"}, PythonTargets...)...); err != nil {
			ui.StepFail(1, totalSteps, "Lint check failed")
			return fmt.Errorf("ruff check failed: %w", err)
		}
		if err := exec.RunPassthrough("ruff", append([]string{"format", "--check"}, PythonTargets...)...); err != nil {
			ui.StepFail(1, totalSteps, "Format check failed")
			return fmt.Errorf("ruff format check failed: %w", err)
		}
		ui.StepDone(1, totalSteps, "Lint passed")

		// Step 2: Typecheck
		ui.Step(2, totalSteps, "Type checking...")
		if err := exec.RunPassthrough("mypy", PythonTargets...); err != nil {
			ui.StepFail(2, totalSteps, "Type check failed")
			return fmt.Errorf("mypy type check failed: %w", err)
		}
		ui.StepDone(2, totalSteps, "Typecheck passed")

		// Step 3: Test
		ui.Step(3, totalSteps, "Running tests...")
		if err := exec.RunPassthrough("pytest", "da_pipeline_tests"); err != nil {
			ui.StepFail(3, totalSteps, "Tests failed")
			return fmt.Errorf("pytest failed: %w", err)
		}
		ui.StepDone(3, totalSteps, "Tests passed")

		ui.Newline()
		ui.Success("All checks passed")
		return nil
	},
}
