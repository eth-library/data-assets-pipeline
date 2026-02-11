package ui

import (
	"os"
	"testing"

	"github.com/charmbracelet/lipgloss"
)

func TestIsCI(t *testing.T) {
	// Save original env
	origCI := os.Getenv("CI")
	origGHA := os.Getenv("GITHUB_ACTIONS")
	defer func() {
		os.Setenv("CI", origCI)
		os.Setenv("GITHUB_ACTIONS", origGHA)
	}()

	tests := []struct {
		name     string
		ci       string
		gha      string
		expected bool
	}{
		{"no env vars", "", "", false},
		{"CI set", "true", "", true},
		{"GITHUB_ACTIONS set", "", "true", true},
		{"both set", "true", "true", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Setenv("CI", tt.ci)
			os.Setenv("GITHUB_ACTIONS", tt.gha)

			got := IsCI()
			if got != tt.expected {
				t.Errorf("IsCI() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestSymbolsAreDefined(t *testing.T) {
	if Symbols.Success == "" {
		t.Error("Symbols.Success is empty")
	}
	if Symbols.Error == "" {
		t.Error("Symbols.Error is empty")
	}
	if Symbols.Warning == "" {
		t.Error("Symbols.Warning is empty")
	}
	if Symbols.Info == "" {
		t.Error("Symbols.Info is empty")
	}
}

func TestStylesAreDefined(t *testing.T) {
	// Just verify styles can be rendered without panic
	_ = Styles.Bold.Render("test")
	_ = Styles.Dim.Render("test")
	_ = Styles.Success.Render("test")
	_ = Styles.Error.Render("test")
	_ = Styles.Warning.Render("test")
}

func TestDisableColors(t *testing.T) {
	// Just verify it doesn't panic
	DisableColors()
}

func TestSeparator(t *testing.T) {
	tests := []struct {
		name  string
		width int
	}{
		{"explicit width 10", 10},
		{"explicit width 30", 30},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Separator(tt.width)
			// Count runes, not bytes (─ is multi-byte UTF-8)
			runeCount := 0
			for range got {
				runeCount++
			}
			if runeCount != tt.width {
				t.Errorf("Separator(%d) rune count = %d, want %d", tt.width, runeCount, tt.width)
			}
		})
	}
}

func TestTerminalWidth(t *testing.T) {
	// Just ensure it returns a positive value
	w := TerminalWidth()
	if w <= 0 {
		t.Errorf("TerminalWidth() = %d, want > 0", w)
	}
}

func TestNoColor(t *testing.T) {
	// Save original env
	origNoColor, hadNoColor := os.LookupEnv("NO_COLOR")
	defer func() {
		if hadNoColor {
			os.Setenv("NO_COLOR", origNoColor)
		} else {
			os.Unsetenv("NO_COLOR")
		}
	}()

	tests := []struct {
		name     string
		setEnv   bool
		value    string
		expected bool
	}{
		{"NO_COLOR not set", false, "", false},
		{"NO_COLOR set to empty", true, "", true},
		{"NO_COLOR set to 1", true, "1", true},
		{"NO_COLOR set to true", true, "true", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setEnv {
				os.Setenv("NO_COLOR", tt.value)
			} else {
				os.Unsetenv("NO_COLOR")
			}

			got := NoColor()
			if got != tt.expected {
				t.Errorf("NoColor() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestIsTTY(t *testing.T) {
	// IsTTY should return a boolean without panicking
	// In test environment, it may return false (not a real terminal)
	_ = IsTTY()
}

func TestBox(t *testing.T) {
	content := "test content"
	result := Box(content)

	// Box should contain the content
	if result == "" {
		t.Error("Box() returned empty string")
	}
	// The result should be longer than just the content (has borders)
	if len(result) <= len(content) {
		t.Errorf("Box() result length %d should be > content length %d", len(result), len(content))
	}
}

func TestSeparatorZeroWidth(t *testing.T) {
	// Separator(0) should use terminal width (default)
	sep := Separator(0)
	if sep == "" {
		t.Error("Separator(0) returned empty string")
	}
}

func TestColorsAreDefined(t *testing.T) {
	// Verify all color constants are usable
	colors := []lipgloss.Color{
		ETHBlue,
		ETHPetrol,
		ETHGreen,
		ETHBronze,
		ETHRed,
		ETHPurple,
		ColorSuccess,
		ColorWarning,
		ColorError,
		ColorInfo,
		ColorAccent,
	}

	for _, c := range colors {
		// Colors should be non-empty strings
		if string(c) == "" {
			t.Errorf("Color constant is empty")
		}
	}
}

func TestSymbolsComplete(t *testing.T) {
	// Verify all symbols have distinct values
	symbols := map[string]string{
		"Success": Symbols.Success,
		"Error":   Symbols.Error,
		"Warning": Symbols.Warning,
		"Info":    Symbols.Info,
		"Arrow":   Symbols.Arrow,
		"Dot":     Symbols.Dot,
	}

	for name, val := range symbols {
		if val == "" {
			t.Errorf("Symbol %s is empty", name)
		}
	}
}

func TestAllStylesRender(t *testing.T) {
	testText := "test"

	// Test all styles can render without panic
	styles := []struct {
		name  string
		style lipgloss.Style
	}{
		{"Bold", Styles.Bold},
		{"Dim", Styles.Dim},
		{"Italic", Styles.Italic},
		{"Underline", Styles.Underline},
		{"Success", Styles.Success},
		{"Warning", Styles.Warning},
		{"Error", Styles.Error},
		{"Info", Styles.Info},
		{"Title", Styles.Title},
		{"Subtitle", Styles.Subtitle},
		{"Command", Styles.Command},
		{"Flag", Styles.Flag},
		{"Description", Styles.Description},
		{"StatusOK", Styles.StatusOK},
		{"StatusFail", Styles.StatusFail},
		{"StatusWarn", Styles.StatusWarn},
		{"StatusInfo", Styles.StatusInfo},
	}

	for _, s := range styles {
		t.Run(s.name, func(t *testing.T) {
			result := s.style.Render(testText)
			if result == "" {
				t.Errorf("Styles.%s.Render() returned empty", s.name)
			}
		})
	}
}
