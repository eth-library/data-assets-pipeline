// Package exec provides utilities for running external commands.
package exec

import (
	"bytes"
	"os"
	"os/exec"
	"strings"
)

// Run executes a command and returns its combined output.
// Returns empty string if the command fails.
func Run(name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(stdout.String()), nil
}

// RunInteractive runs a command with stdin/stdout/stderr connected to the terminal.
func RunInteractive(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// RunPassthrough runs a command, passing through all output to the terminal.
// Returns the exit code.
func RunPassthrough(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// Which checks if a command exists in PATH.
func Which(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}
