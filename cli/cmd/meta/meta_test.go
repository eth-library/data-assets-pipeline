package meta

import (
	"os/exec"
	"path/filepath"
	"testing"
)

func TestCommands(t *testing.T) {
	cmds := Commands()
	if len(cmds) != 0 {
		t.Errorf("Commands() returned %d commands, want 0", len(cmds))
	}
}

func TestCliCmdIsVisible(t *testing.T) {
	if CliCmd.Hidden {
		t.Error("CliCmd should be visible")
	}
}

func TestCliCmdHasSubcommands(t *testing.T) {
	if !CliCmd.HasSubCommands() {
		t.Error("CliCmd should have subcommands (build, test)")
	}
}

func TestFindCliDir(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not available (Nix build sandbox)")
	}
	dir, err := findCliDir()
	if err != nil {
		t.Fatalf("findCliDir() failed: %v", err)
	}
	if filepath.Base(dir) != "cli" {
		t.Errorf("findCliDir() = %q, want path ending in cli", dir)
	}
	if !filepath.IsAbs(dir) {
		t.Errorf("findCliDir() = %q, want absolute path", dir)
	}
}
