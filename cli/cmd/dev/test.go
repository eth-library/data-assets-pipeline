package dev

import (
	"github.com/eth-library/dap/cli/internal/exec"
	"github.com/spf13/cobra"
)

// TestCmd runs the pytest test suite.
var TestCmd = &cobra.Command{
	Use:   "test [pytest-args...]",
	Short: "Run pytest tests",
	Long: `Run the pytest test suite on da_pipeline_tests/.

All arguments are passed directly to pytest.

Examples:
  dap test                Run all tests
  dap test -v             Verbose output
  dap test -k "test_foo"  Run tests matching pattern
  dap test --lf           Re-run last failed tests`,
	GroupID:               GroupID,
	DisableFlagParsing:    true,
	DisableFlagsInUseLine: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Handle --help/-h manually since DisableFlagParsing is true
		for _, arg := range args {
			if arg == "--help" || arg == "-h" {
				return cmd.Help()
			}
		}
		pytestArgs := []string{"da_pipeline_tests"}
		pytestArgs = append(pytestArgs, args...)
		return exec.RunPassthrough("pytest", pytestArgs...)
	},
}
