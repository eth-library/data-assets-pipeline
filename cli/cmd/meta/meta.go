// Package meta contains CLI maintenance commands (hidden from regular users).
package meta

import "github.com/spf13/cobra"

// Commands returns hidden CLI maintenance commands.
func Commands() []*cobra.Command {
	return []*cobra.Command{}
}
