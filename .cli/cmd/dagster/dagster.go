// Package dagster contains Dagster-related commands.
package dagster

import "github.com/spf13/cobra"

// GroupID for dagster commands (matches root.go GroupDagster)
const GroupID = "dagster"

// Commands returns all dagster commands to be registered with the root command.
func Commands() []*cobra.Command {
	return []*cobra.Command{
		MaterializeCmd,
		RunCmd,
	}
}
