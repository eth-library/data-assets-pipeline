// Package main is the entry point for the dap CLI.
// dap is a developer experience tool for the Data Archive Pipeline.
package main

import (
	"os"

	"github.com/eth-library/dap/cli/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
