package dagster

import (
	"testing"
)

func TestCommands(t *testing.T) {
	cmds := Commands()
	if len(cmds) != 2 {
		t.Errorf("Commands() returned %d commands, want 2", len(cmds))
	}

	// Verify expected commands
	names := make(map[string]bool)
	for _, cmd := range cmds {
		names[cmd.Use] = true
	}

	if !names["materialize [flags]"] {
		t.Error("missing materialize command")
	}
	if !names["run [flags]"] {
		t.Error("missing run command")
	}
}

func TestGroupID(t *testing.T) {
	if GroupID != "dagster" {
		t.Errorf("GroupID = %q, want %q", GroupID, "dagster")
	}
}

func TestCommandsHaveDisabledFlagParsing(t *testing.T) {
	// Dagster commands should pass flags through to dagster
	for _, cmd := range Commands() {
		if !cmd.DisableFlagParsing {
			t.Errorf("command %q should have DisableFlagParsing=true", cmd.Use)
		}
	}
}
