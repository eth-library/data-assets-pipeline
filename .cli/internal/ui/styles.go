// Package ui provides terminal UI components and styling for dap CLI.
//
// Design Philosophy:
// - Minimal and professional aesthetic inspired by gh, gum, and docker CLI
// - ETH Zurich brand colors where they provide good terminal readability
// - Graceful degradation in CI environments (no colors, no animations)
// - Consistent visual language across all commands
package ui

import (
	"os"

	"github.com/charmbracelet/lipgloss"
	"github.com/mattn/go-isatty"
	"github.com/muesli/termenv"
)

// ETH Zurich brand colors (adjusted for terminal readability)
// See: https://ethz.ch/staffnet/en/service/communication/corporate-design/colours.html
var (
	// Primary colors
	ETHBlue   = lipgloss.Color("#215CAF") // Headers, primary accent
	ETHPetrol = lipgloss.Color("#007894") // Secondary accent, info
	ETHGreen  = lipgloss.Color("#627313") // Success (may need brightening on dark bg)
	ETHBronze = lipgloss.Color("#8E6713") // Warnings
	ETHRed    = lipgloss.Color("#B7352D") // Errors
	ETHPurple = lipgloss.Color("#A7117A") // Special highlights

	// Semantic aliases for better terminal contrast
	ColorSuccess = lipgloss.Color("#22c55e") // Brighter green for dark terminals
	ColorWarning = ETHBronze
	ColorError   = ETHRed
	ColorInfo    = ETHPetrol
	ColorAccent  = ETHBlue

	// Adaptive colors for light/dark terminal accessibility
	ColorMuted = lipgloss.AdaptiveColor{Light: "#555555", Dark: "#a0a0a0"}
)

// Styles defines reusable lipgloss styles for consistent output.
var Styles = struct {
	// Text styles
	Bold      lipgloss.Style
	Dim       lipgloss.Style
	Italic    lipgloss.Style
	Underline lipgloss.Style

	// Semantic styles
	Success lipgloss.Style
	Warning lipgloss.Style
	Error   lipgloss.Style
	Info    lipgloss.Style

	// Component styles
	Title       lipgloss.Style
	Subtitle    lipgloss.Style
	Command     lipgloss.Style
	Flag        lipgloss.Style
	Description lipgloss.Style

	// Status indicators
	StatusOK   lipgloss.Style
	StatusFail lipgloss.Style
	StatusWarn lipgloss.Style
	StatusInfo lipgloss.Style
}{
	// Text styles
	Bold:      lipgloss.NewStyle().Bold(true),
	Dim:       lipgloss.NewStyle().Foreground(ColorMuted),
	Italic:    lipgloss.NewStyle().Italic(true),
	Underline: lipgloss.NewStyle().Underline(true),

	// Semantic styles
	Success: lipgloss.NewStyle().Foreground(ColorSuccess),
	Warning: lipgloss.NewStyle().Foreground(ColorWarning),
	Error:   lipgloss.NewStyle().Foreground(ColorError),
	Info:    lipgloss.NewStyle().Foreground(ColorInfo),

	// Component styles
	Title:       lipgloss.NewStyle().Bold(true).Foreground(ETHBlue),
	Subtitle:    lipgloss.NewStyle().Foreground(ColorMuted),
	Command:     lipgloss.NewStyle().Bold(true).Foreground(ETHPetrol),
	Flag:        lipgloss.NewStyle().Foreground(ETHPurple),
	Description: lipgloss.NewStyle().Foreground(ColorMuted),

	// Status indicators with symbols
	StatusOK:   lipgloss.NewStyle().Foreground(ColorSuccess),
	StatusFail: lipgloss.NewStyle().Foreground(ColorError),
	StatusWarn: lipgloss.NewStyle().Foreground(ColorWarning),
	StatusInfo: lipgloss.NewStyle().Foreground(ColorInfo),
}

// Symbols for status indicators
var Symbols = struct {
	Success string
	Error   string
	Warning string
	Info    string
	Arrow   string
	Dot     string
}{
	Success: "✓",
	Error:   "✗",
	Warning: "!",
	Info:    "→",
	Arrow:   "→",
	Dot:     "•",
}

// IsCI returns true if running in a CI environment.
func IsCI() bool {
	return os.Getenv("CI") != "" || os.Getenv("GITHUB_ACTIONS") != ""
}

// IsTTY returns true if stderr is a terminal.
// Note: CLI output goes to stderr, so we check stderr not stdout.
func IsTTY() bool {
	return isatty.IsTerminal(os.Stderr.Fd()) || isatty.IsCygwinTerminal(os.Stderr.Fd())
}

// DisableColors disables all color output (for CI or --no-color flag).
func DisableColors() {
	lipgloss.SetColorProfile(termenv.Ascii)
}

// NoColor returns true if NO_COLOR env var is set (any value).
// See: https://no-color.org/
func NoColor() bool {
	_, exists := os.LookupEnv("NO_COLOR")
	return exists
}

func init() {
	// Auto-disable colors per no-color.org standard and for CI/non-TTY
	if NoColor() || IsCI() || !IsTTY() {
		DisableColors()
	}
}
