package meta

import (
	"os"
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
	// Save current directory
	origDir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(origDir)

	// Test from within cli directory (where go.mod and gomod2nix.toml exist)
	if _, err := os.Stat("go.mod"); err == nil {
		if _, err := os.Stat("gomod2nix.toml"); err == nil {
			dir, err := findCliDir()
			if err != nil {
				t.Errorf("findCliDir() from cli failed: %v", err)
			}
			if dir != "." {
				t.Errorf("findCliDir() from cli = %q, want %q", dir, ".")
			}
		}
	}
}

func TestFindCliDirFromParent(t *testing.T) {
	// Save current directory
	origDir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(origDir)

	// Go up to parent directory
	parent := filepath.Dir(origDir)
	if err := os.Chdir(parent); err != nil {
		t.Skip("cannot change to parent directory")
	}

	// Check if cli exists in parent (may not exist in Nix build environment)
	if _, err := os.Stat("cli"); os.IsNotExist(err) {
		t.Skip("not in expected directory structure (Nix build environment)")
	}

	dir, err := findCliDir()
	if err != nil {
		t.Errorf("findCliDir() from parent failed: %v", err)
		return
	}
	// Should find cli (either relative or absolute path ending in cli)
	if dir != "cli" && filepath.Base(dir) != "cli" {
		t.Errorf("findCliDir() from parent = %q, want cli or path ending in cli", dir)
	}
}
