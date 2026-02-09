package ui

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

// captureStderr captures stderr output from a function
func captureStderr(fn func()) string {
	// Save original stderr
	oldStderr := os.Stderr

	// Create a pipe
	r, w, _ := os.Pipe()
	os.Stderr = w

	// Run the function
	fn()

	// Close writer and restore stderr
	w.Close()
	os.Stderr = oldStderr

	// Read captured output
	var buf bytes.Buffer
	buf.ReadFrom(r)
	return buf.String()
}

func TestLogger(t *testing.T) {
	// Logger should be initialized
	if Logger == nil {
		t.Error("Logger is nil")
	}
}

func TestBanner(t *testing.T) {
	output := captureStderr(func() {
		Banner("Test Title", "Test Subtitle")
	})

	if output == "" {
		t.Error("Banner() produced no output")
	}
	// In CI mode, should contain simple text
	if !strings.Contains(output, "Test Title") && !strings.Contains(output, "DAP") {
		t.Error("Banner() should contain title or DAP text")
	}
}

func TestTitle(t *testing.T) {
	output := captureStderr(func() {
		Title("Test Title")
	})

	if !strings.Contains(output, "Test Title") {
		t.Errorf("Title() output should contain 'Test Title', got: %s", output)
	}
}

func TestSubtitle(t *testing.T) {
	output := captureStderr(func() {
		Subtitle("Test Subtitle")
	})

	if !strings.Contains(output, "Test Subtitle") {
		t.Errorf("Subtitle() output should contain 'Test Subtitle', got: %s", output)
	}
}

func TestSection(t *testing.T) {
	output := captureStderr(func() {
		Section("Test Section")
	})

	if !strings.Contains(output, "Test Section") {
		t.Errorf("Section() output should contain 'Test Section', got: %s", output)
	}
}

func TestDivider(t *testing.T) {
	output := captureStderr(func() {
		Divider()
	})

	if output == "" {
		t.Error("Divider() produced no output")
	}
}

func TestSuccess(t *testing.T) {
	output := captureStderr(func() {
		Success("Operation completed")
	})

	if !strings.Contains(output, "Operation completed") {
		t.Errorf("Success() should contain message, got: %s", output)
	}
}

func TestSuccessWithKeyvals(t *testing.T) {
	output := captureStderr(func() {
		Success("Done", "count", 5, "status", "ok")
	})

	if !strings.Contains(output, "Done") {
		t.Error("Success() should contain message")
	}
	if !strings.Contains(output, "count") {
		t.Error("Success() should contain key 'count'")
	}
}

func TestError(t *testing.T) {
	output := captureStderr(func() {
		Error("Something failed")
	})

	if !strings.Contains(output, "Something failed") {
		t.Errorf("Error() should contain message, got: %s", output)
	}
}

func TestWarn(t *testing.T) {
	output := captureStderr(func() {
		Warn("Warning message")
	})

	if !strings.Contains(output, "Warning message") {
		t.Errorf("Warn() should contain message, got: %s", output)
	}
}

func TestInfo(t *testing.T) {
	output := captureStderr(func() {
		Info("Info message")
	})

	if !strings.Contains(output, "Info message") {
		t.Errorf("Info() should contain message, got: %s", output)
	}
}

func TestKeyValue(t *testing.T) {
	output := captureStderr(func() {
		KeyValue("Version", "1.0.0")
	})

	if !strings.Contains(output, "Version") {
		t.Error("KeyValue() should contain key")
	}
	if !strings.Contains(output, "1.0.0") {
		t.Error("KeyValue() should contain value")
	}
}

func TestKeyValueStatus(t *testing.T) {
	t.Run("ok status", func(t *testing.T) {
		output := captureStderr(func() {
			KeyValueStatus("Python", "3.12.0", true)
		})

		if !strings.Contains(output, "Python") {
			t.Error("KeyValueStatus() should contain key")
		}
		if !strings.Contains(output, "3.12.0") {
			t.Error("KeyValueStatus() should contain value")
		}
	})

	t.Run("not ok status", func(t *testing.T) {
		output := captureStderr(func() {
			KeyValueStatus("Python", "not found", false)
		})

		if !strings.Contains(output, "Python") {
			t.Error("KeyValueStatus() should contain key")
		}
	})
}

func TestKeyValueDim(t *testing.T) {
	output := captureStderr(func() {
		KeyValueDim("Optional", "not set")
	})

	if !strings.Contains(output, "Optional") {
		t.Error("KeyValueDim() should contain key")
	}
	if !strings.Contains(output, "not set") {
		t.Error("KeyValueDim() should contain value")
	}
}

func TestStep(t *testing.T) {
	output := captureStderr(func() {
		Step(1, 3, "Installing dependencies")
	})

	if !strings.Contains(output, "1") && !strings.Contains(output, "3") {
		t.Error("Step() should contain step numbers")
	}
	if !strings.Contains(output, "Installing dependencies") {
		t.Error("Step() should contain description")
	}
}

func TestStepDone(t *testing.T) {
	output := captureStderr(func() {
		StepDone(1, 3, "Installed dependencies")
	})

	if !strings.Contains(output, "Installed dependencies") {
		t.Error("StepDone() should contain description")
	}
}

func TestStepFail(t *testing.T) {
	output := captureStderr(func() {
		StepFail(2, 3, "Failed to compile")
	})

	if !strings.Contains(output, "Failed to compile") {
		t.Error("StepFail() should contain description")
	}
}

func TestTaskStart(t *testing.T) {
	output := captureStderr(func() {
		TaskStart("Running tests")
	})

	if !strings.Contains(output, "Running tests") {
		t.Error("TaskStart() should contain task name")
	}
}

func TestTaskDone(t *testing.T) {
	output := captureStderr(func() {
		TaskDone("Tests passed")
	})

	if !strings.Contains(output, "Tests passed") {
		t.Error("TaskDone() should contain task name")
	}
}

func TestTaskFail(t *testing.T) {
	output := captureStderr(func() {
		TaskFail("Tests failed")
	})

	if !strings.Contains(output, "Tests failed") {
		t.Error("TaskFail() should contain task name")
	}
}

func TestCommandHint(t *testing.T) {
	output := captureStderr(func() {
		CommandHint("dap dev", "Start development server")
	})

	if !strings.Contains(output, "dap dev") {
		t.Error("CommandHint() should contain command")
	}
	if !strings.Contains(output, "Start development server") {
		t.Error("CommandHint() should contain description")
	}
}

func TestHint(t *testing.T) {
	output := captureStderr(func() {
		Hint("Run dap --help for more info")
	})

	if !strings.Contains(output, "Run dap --help") {
		t.Error("Hint() should contain hint text")
	}
}

func TestListItem(t *testing.T) {
	output := captureStderr(func() {
		ListItem("First item")
	})

	if !strings.Contains(output, "First item") {
		t.Error("ListItem() should contain item text")
	}
}

func TestListItemStatus(t *testing.T) {
	t.Run("ok status", func(t *testing.T) {
		output := captureStderr(func() {
			ListItemStatus("Passed test", true)
		})

		if !strings.Contains(output, "Passed test") {
			t.Error("ListItemStatus() should contain item text")
		}
	})

	t.Run("not ok status", func(t *testing.T) {
		output := captureStderr(func() {
			ListItemStatus("Failed test", false)
		})

		if !strings.Contains(output, "Failed test") {
			t.Error("ListItemStatus() should contain item text")
		}
	})
}

func TestSuccessBox(t *testing.T) {
	output := captureStderr(func() {
		SuccessBox("Success!", "Operation completed successfully")
	})

	if output == "" {
		t.Error("SuccessBox() produced no output")
	}
	if !strings.Contains(output, "Success!") {
		t.Error("SuccessBox() should contain title")
	}
}

func TestErrorBox(t *testing.T) {
	output := captureStderr(func() {
		ErrorBox("Error!", "Something went wrong")
	})

	if output == "" {
		t.Error("ErrorBox() produced no output")
	}
	if !strings.Contains(output, "Error!") {
		t.Error("ErrorBox() should contain title")
	}
}

func TestNewline(t *testing.T) {
	output := captureStderr(func() {
		Newline()
	})

	if output != "\n" {
		t.Errorf("Newline() should produce a single newline, got: %q", output)
	}
}

func TestPrintStatusLineWithOddKeyvals(t *testing.T) {
	// Test with odd number of keyvals (incomplete pair)
	output := captureStderr(func() {
		Success("Test", "key1", "value1", "orphan")
	})

	// Should still contain the message and complete pairs
	if !strings.Contains(output, "Test") {
		t.Error("Should contain message")
	}
	if !strings.Contains(output, "key1") {
		t.Error("Should contain key1")
	}
}
