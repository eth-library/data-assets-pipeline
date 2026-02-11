// Package dev contains development-related commands.
package dev

import "github.com/spf13/cobra"

// GroupID for development commands (matches root.go GroupDevelopment)
const GroupID = "development"

// PythonTargets defines the Python packages to check/lint/format.
var PythonTargets = []string{"da_pipeline", "da_pipeline_tests"}

// Commands returns all development commands to be registered with the root command.
func Commands() []*cobra.Command {
	return []*cobra.Command{
		CheckCmd,
		DevServerCmd,
		LintCmd,
		TestCmd,
		TypecheckCmd,
	}
}
