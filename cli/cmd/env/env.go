// Package env contains environment-related commands.
package env

import "github.com/spf13/cobra"

// GroupID for environment commands (matches root.go GroupEnvironment)
const GroupID = "environment"

// Commands returns all environment commands to be registered with the root command.
func Commands() []*cobra.Command {
	return []*cobra.Command{
		CleanCmd,
		ResetCmd,
		VersionsCmd,
		WelcomeCmd,
	}
}
