package cmd

import (
	"bytes"
	"testing"
)

func TestExecute(t *testing.T) {
	// Execute with --help should not error
	rootCmd.SetArgs([]string{"--help"})

	// Capture output
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)

	err := rootCmd.Execute()
	if err != nil {
		t.Errorf("Execute() with --help returned error: %v", err)
	}
}

func TestRootCommandGroups(t *testing.T) {
	groups := rootCmd.Groups()

	expectedGroups := []string{
		GroupDevelopment,
		GroupEnvironment,
		GroupDagster,
		GroupKubernetes,
		GroupGoCLI,
	}

	if len(groups) != len(expectedGroups) {
		t.Errorf("Expected %d groups, got %d", len(expectedGroups), len(groups))
	}

	// Verify all expected groups exist
	groupIDs := make(map[string]bool)
	for _, g := range groups {
		groupIDs[g.ID] = true
	}

	for _, expected := range expectedGroups {
		if !groupIDs[expected] {
			t.Errorf("Expected group %q not found", expected)
		}
	}
}

func TestRootCommandHasSubcommands(t *testing.T) {
	commands := rootCmd.Commands()

	if len(commands) == 0 {
		t.Error("Root command has no subcommands")
	}

	// Verify some expected commands exist
	expectedCommands := []string{"dev", "test", "lint", "check", "k8s", "cli"}
	commandNames := make(map[string]bool)
	for _, cmd := range commands {
		commandNames[cmd.Name()] = true
	}

	for _, expected := range expectedCommands {
		if !commandNames[expected] {
			t.Errorf("Expected command %q not found", expected)
		}
	}
}

func TestNoColorFlag(t *testing.T) {
	flag := rootCmd.PersistentFlags().Lookup("no-color")
	if flag == nil {
		t.Error("--no-color flag not defined")
	}

	if flag.DefValue != "false" {
		t.Errorf("--no-color default should be 'false', got %q", flag.DefValue)
	}
}

func TestGroupConstants(t *testing.T) {
	// Verify group constants have expected values
	tests := []struct {
		name     string
		constant string
		expected string
	}{
		{"GroupDevelopment", GroupDevelopment, "development"},
		{"GroupEnvironment", GroupEnvironment, "environment"},
		{"GroupDagster", GroupDagster, "dagster"},
		{"GroupKubernetes", GroupKubernetes, "kubernetes"},
		{"GroupGoCLI", GroupGoCLI, "gocli"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.constant != tt.expected {
				t.Errorf("%s = %q, want %q", tt.name, tt.constant, tt.expected)
			}
		})
	}
}

func TestRootCommandUsage(t *testing.T) {
	if rootCmd.Use != "dap" {
		t.Errorf("Root command Use = %q, want 'dap'", rootCmd.Use)
	}

	if rootCmd.Short == "" {
		t.Error("Root command Short description is empty")
	}

	if rootCmd.Long == "" {
		t.Error("Root command Long description is empty")
	}
}

func TestCompletionDisabled(t *testing.T) {
	if !rootCmd.CompletionOptions.DisableDefaultCmd {
		t.Error("Completion should be disabled")
	}
}
