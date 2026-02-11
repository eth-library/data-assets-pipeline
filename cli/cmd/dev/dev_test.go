package dev

import (
	"testing"
)

func TestCommands(t *testing.T) {
	cmds := Commands()
	if len(cmds) == 0 {
		t.Error("Commands() returned empty slice")
	}

	// Verify expected command count
	if len(cmds) != 5 {
		t.Errorf("Commands() returned %d commands, want 5", len(cmds))
	}

	// Verify all commands have GroupID set
	for _, cmd := range cmds {
		if cmd.GroupID != GroupID {
			t.Errorf("command %q has GroupID %q, want %q", cmd.Name(), cmd.GroupID, GroupID)
		}
	}
}

func TestPythonTargets(t *testing.T) {
	if len(PythonTargets) == 0 {
		t.Error("PythonTargets is empty")
	}

	// Should contain the main packages
	hasMain := false
	hasTests := false
	for _, target := range PythonTargets {
		if target == "da_pipeline" {
			hasMain = true
		}
		if target == "da_pipeline_tests" {
			hasTests = true
		}
	}

	if !hasMain {
		t.Error("PythonTargets missing da_pipeline")
	}
	if !hasTests {
		t.Error("PythonTargets missing da_pipeline_tests")
	}
}

func TestGroupID(t *testing.T) {
	if GroupID != "development" {
		t.Errorf("GroupID = %q, want %q", GroupID, "development")
	}
}

func TestCheckCmd(t *testing.T) {
	if CheckCmd.Use != "check" {
		t.Errorf("CheckCmd.Use = %q, want 'check'", CheckCmd.Use)
	}
	if CheckCmd.Short == "" {
		t.Error("CheckCmd.Short is empty")
	}
	if CheckCmd.Long == "" {
		t.Error("CheckCmd.Long is empty")
	}
	if CheckCmd.GroupID != GroupID {
		t.Errorf("CheckCmd.GroupID = %q, want %q", CheckCmd.GroupID, GroupID)
	}
	if CheckCmd.RunE == nil {
		t.Error("CheckCmd.RunE is nil")
	}
}

func TestDevServerCmd(t *testing.T) {
	if DevServerCmd.Use != "dev" {
		t.Errorf("DevServerCmd.Use = %q, want 'dev'", DevServerCmd.Use)
	}
	if DevServerCmd.Short == "" {
		t.Error("DevServerCmd.Short is empty")
	}
	if DevServerCmd.GroupID != GroupID {
		t.Errorf("DevServerCmd.GroupID = %q, want %q", DevServerCmd.GroupID, GroupID)
	}
	if DevServerCmd.RunE == nil {
		t.Error("DevServerCmd.RunE is nil")
	}
}

func TestLintCmd(t *testing.T) {
	if LintCmd.Use != "lint" {
		t.Errorf("LintCmd.Use = %q, want 'lint'", LintCmd.Use)
	}
	if LintCmd.Short == "" {
		t.Error("LintCmd.Short is empty")
	}
	if LintCmd.GroupID != GroupID {
		t.Errorf("LintCmd.GroupID = %q, want %q", LintCmd.GroupID, GroupID)
	}
	if LintCmd.RunE == nil {
		t.Error("LintCmd.RunE is nil")
	}

	// Verify --fix flag exists
	flag := LintCmd.Flags().Lookup("fix")
	if flag == nil {
		t.Error("LintCmd should have --fix flag")
	}
}

func TestTestCmd(t *testing.T) {
	if TestCmd.Use != "test [pytest-args...]" {
		t.Errorf("TestCmd.Use = %q, want 'test [pytest-args...]'", TestCmd.Use)
	}
	if TestCmd.Short == "" {
		t.Error("TestCmd.Short is empty")
	}
	if TestCmd.GroupID != GroupID {
		t.Errorf("TestCmd.GroupID = %q, want %q", TestCmd.GroupID, GroupID)
	}
	if TestCmd.RunE == nil {
		t.Error("TestCmd.RunE is nil")
	}
}

func TestTypecheckCmd(t *testing.T) {
	if TypecheckCmd.Use != "typecheck" {
		t.Errorf("TypecheckCmd.Use = %q, want 'typecheck'", TypecheckCmd.Use)
	}
	if TypecheckCmd.Short == "" {
		t.Error("TypecheckCmd.Short is empty")
	}
	if TypecheckCmd.GroupID != GroupID {
		t.Errorf("TypecheckCmd.GroupID = %q, want %q", TypecheckCmd.GroupID, GroupID)
	}
	if TypecheckCmd.RunE == nil {
		t.Error("TypecheckCmd.RunE is nil")
	}
}

func TestCommandNames(t *testing.T) {
	cmds := Commands()

	expectedNames := map[string]bool{
		"check":     false,
		"dev":       false,
		"lint":      false,
		"test":      false,
		"typecheck": false,
	}

	for _, cmd := range cmds {
		if _, ok := expectedNames[cmd.Name()]; ok {
			expectedNames[cmd.Name()] = true
		}
	}

	for name, found := range expectedNames {
		if !found {
			t.Errorf("Expected command %q not found in Commands()", name)
		}
	}
}
