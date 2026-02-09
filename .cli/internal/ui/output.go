package ui

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	"golang.org/x/term"
)

// Logger is the global logger instance for dap CLI.
var Logger *log.Logger

func init() {
	// Configure the logger with a clean, minimal style
	Logger = log.NewWithOptions(os.Stderr, log.Options{
		ReportTimestamp: true,
		TimeFormat:      time.Kitchen,
		Prefix:          "",
	})

	// Adjust log level based on environment
	if IsCI() {
		Logger.SetReportTimestamp(false)
	}
}

// Layout constants
const (
	defaultSeparatorWidth = 80
)

// TerminalWidth returns the current terminal width, or a default if unavailable.
func TerminalWidth() int {
	if w, _, err := term.GetSize(int(os.Stderr.Fd())); err == nil && w > 0 {
		return w
	}
	return defaultSeparatorWidth
}

// Separator returns a horizontal line of the specified width.
// If width is 0, it uses the full terminal width.
func Separator(width int) string {
	if width <= 0 {
		width = TerminalWidth()
	}
	return strings.Repeat("─", width)
}

// ═══════════════════════════════════════════════════════════════════════════
// Banner & Headers
// ═══════════════════════════════════════════════════════════════════════════

// Banner prints a styled banner with ASCII art logo centered above text.
// In CI/CD environments, prints a simple text banner instead.
func Banner(title, subtitle string) {
	// Simple banner for CI/CD
	if IsCI() || !IsTTY() {
		fmt.Fprintln(os.Stderr)
		fmt.Fprintln(os.Stderr, "  DAP-O | "+title)
		fmt.Fprintln(os.Stderr, "          "+subtitle)
		return
	}

	// Fancy banner for TTY
	logoStyle := lipgloss.NewStyle().
		Foreground(ETHBlue).
		Bold(true)

	titleStyle := lipgloss.NewStyle().
		Foreground(ColorMuted)

	subtitleStyle := lipgloss.NewStyle().
		Foreground(ETHPetrol)

	logoLines := []string{
		"██████╗  █████╗ ██████╗          ██████╗",
		"██╔══██╗██╔══██╗██╔══██╗        ██╔═══██╗",
		"██║  ██║███████║██████╔╝ █████╗ ██║   ██║",
		"██║  ██║██╔══██║██╔═══╝  ╚════╝ ██║   ██║",
		"██████╔╝██║  ██║██║             ╚██████╔╝",
		"╚═════╝ ╚═╝  ╚═╝╚═╝              ╚═════╝",
	}

	// Calculate padding to center logo above title
	logoWidth := len([]rune(logoLines[0]))
	titleWidth := len([]rune(title))
	padding := "  " // Base padding
	if titleWidth > logoWidth {
		// Center logo over title
		extra := (titleWidth - logoWidth) / 2
		padding = "  " + strings.Repeat(" ", extra)
	}

	fmt.Fprintln(os.Stderr)
	for _, line := range logoLines {
		fmt.Fprintln(os.Stderr, padding+logoStyle.Render(line))
	}
	fmt.Fprintln(os.Stderr, "  "+titleStyle.Render(title))
	fmt.Fprintln(os.Stderr, "  "+subtitleStyle.Render(subtitle))
}

// Title prints a styled title/header.
func Title(text string) {
	style := lipgloss.NewStyle().
		Bold(true).
		Foreground(ETHBlue).
		MarginTop(1)
	fmt.Fprintln(os.Stderr, style.Render(text))
}

// Subtitle prints a styled subtitle.
func Subtitle(text string) {
	fmt.Fprintln(os.Stderr, Styles.Subtitle.Render(text))
}

// Section prints a section header.
func Section(title string) {
	style := lipgloss.NewStyle().
		Foreground(ETHPetrol).
		Bold(true).
		MarginTop(1)

	fmt.Fprintln(os.Stderr, style.Render(title))
}

// Divider prints a horizontal line.
func Divider() {
	fmt.Fprintln(os.Stderr, Styles.Dim.Render(Separator(0)))
}

// ═══════════════════════════════════════════════════════════════════════════
// Status Messages
// ═══════════════════════════════════════════════════════════════════════════

// Success prints a success message with a checkmark.
func Success(msg string, keyvals ...interface{}) {
	icon := Styles.StatusOK.Render(Symbols.Success)
	text := Styles.Success.Render(msg)
	printStatusLine(icon, text, keyvals...)
}

// Error prints an error message with an X mark.
func Error(msg string, keyvals ...interface{}) {
	icon := Styles.StatusFail.Render(Symbols.Error)
	text := Styles.Error.Render(msg)
	printStatusLine(icon, text, keyvals...)
}

// Warn prints a warning message with an exclamation mark.
func Warn(msg string, keyvals ...interface{}) {
	icon := Styles.StatusWarn.Render(Symbols.Warning)
	text := Styles.Warning.Render(msg)
	printStatusLine(icon, text, keyvals...)
}

// Info prints an info message with an arrow.
func Info(msg string, keyvals ...interface{}) {
	icon := Styles.StatusInfo.Render(Symbols.Info)
	printStatusLine(icon, msg, keyvals...)
}

func printStatusLine(icon, msg string, keyvals ...interface{}) {
	fmt.Fprintf(os.Stderr, "%s %s", icon, msg)
	if len(keyvals) > 0 {
		fmt.Fprintf(os.Stderr, " ")
		for i := 0; i < len(keyvals); i += 2 {
			if i+1 < len(keyvals) {
				fmt.Fprintf(os.Stderr, "%s=%v ", Styles.Dim.Render(fmt.Sprint(keyvals[i])), keyvals[i+1])
			}
		}
	}
	fmt.Fprintln(os.Stderr)
}

// ═══════════════════════════════════════════════════════════════════════════
// Key-Value Display
// ═══════════════════════════════════════════════════════════════════════════

// KeyValue prints a key-value pair with consistent formatting.
func KeyValue(key, value string) {
	keyStyle := Styles.Dim.Width(14)
	fmt.Fprintf(os.Stderr, "  %s %s\n", keyStyle.Render(key), value)
}

// KeyValueStatus prints a key-value pair with a status indicator.
func KeyValueStatus(key, value string, ok bool) {
	keyStyle := Styles.Dim.Width(14)
	var status string
	if ok {
		status = Styles.StatusOK.Render(Symbols.Success)
	} else {
		status = Styles.StatusWarn.Render(Symbols.Warning)
	}
	fmt.Fprintf(os.Stderr, "  %s %s %s\n", keyStyle.Render(key), value, status)
}

// KeyValueDim prints a key-value pair with dimmed value (for "not set" etc).
func KeyValueDim(key, value string) {
	keyStyle := Styles.Dim.Width(14)
	valueStyle := Styles.Dim
	fmt.Fprintf(os.Stderr, "  %s %s\n", keyStyle.Render(key), valueStyle.Render(value))
}

// ═══════════════════════════════════════════════════════════════════════════
// Task & Step Indicators
// ═══════════════════════════════════════════════════════════════════════════

// Step prints a step indicator for multi-step operations.
func Step(current, total int, description string) {
	stepStyle := lipgloss.NewStyle().
		Foreground(ETHBlue).
		Bold(true)

	descStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#CCCCCC"))

	step := stepStyle.Render(fmt.Sprintf("[%d/%d]", current, total))
	desc := descStyle.Render(description)

	fmt.Fprintf(os.Stderr, "%s %s\n", step, desc)
}

// StepDone prints a completed step.
func StepDone(current, total int, description string) {
	stepStyle := lipgloss.NewStyle().
		Foreground(ColorSuccess).
		Bold(true)

	step := stepStyle.Render(fmt.Sprintf("[%d/%d]", current, total))
	icon := Styles.StatusOK.Render(Symbols.Success)

	fmt.Fprintf(os.Stderr, "%s %s %s\n", step, icon, description)
}

// StepFail prints a failed step.
func StepFail(current, total int, description string) {
	stepStyle := lipgloss.NewStyle().
		Foreground(ColorError).
		Bold(true)

	step := stepStyle.Render(fmt.Sprintf("[%d/%d]", current, total))
	icon := Styles.StatusFail.Render(Symbols.Error)

	fmt.Fprintf(os.Stderr, "%s %s %s\n", step, icon, description)
}

// TaskStart prints a task starting message with arrow.
func TaskStart(task string) {
	icon := lipgloss.NewStyle().Foreground(ETHPetrol).Render("▸")
	fmt.Fprintf(os.Stderr, "%s %s\n", icon, task)
}

// TaskDone prints a task completed message.
func TaskDone(task string) {
	icon := Styles.StatusOK.Render(Symbols.Success)
	fmt.Fprintf(os.Stderr, "%s %s\n", icon, Styles.Success.Render(task))
}

// TaskFail prints a task failed message.
func TaskFail(task string) {
	icon := Styles.StatusFail.Render(Symbols.Error)
	fmt.Fprintf(os.Stderr, "%s %s\n", icon, Styles.Error.Render(task))
}

// ═══════════════════════════════════════════════════════════════════════════
// Command Hints
// ═══════════════════════════════════════════════════════════════════════════

// CommandHint prints a command hint for the user.
func CommandHint(cmd, description string) {
	cmdStyle := Styles.Command.Width(14)
	fmt.Fprintf(os.Stderr, "  %s %s\n", cmdStyle.Render(cmd), Styles.Description.Render(description))
}

// Hint prints a dimmed hint/tip message.
func Hint(text string) {
	fmt.Fprintf(os.Stderr, "  %s\n", Styles.Dim.Render(text))
}

// ═══════════════════════════════════════════════════════════════════════════
// Tables & Lists
// ═══════════════════════════════════════════════════════════════════════════

// ListItem prints a bulleted list item.
func ListItem(text string) {
	bullet := Styles.Dim.Render("•")
	fmt.Fprintf(os.Stderr, "  %s %s\n", bullet, text)
}

// ListItemStatus prints a list item with status icon.
func ListItemStatus(text string, ok bool) {
	var icon string
	if ok {
		icon = Styles.StatusOK.Render(Symbols.Success)
	} else {
		icon = Styles.StatusFail.Render(Symbols.Error)
	}
	fmt.Fprintf(os.Stderr, "  %s %s\n", icon, text)
}

// ═══════════════════════════════════════════════════════════════════════════
// Boxes & Panels
// ═══════════════════════════════════════════════════════════════════════════

// Box renders content in a bordered box.
func Box(content string) string {
	style := lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(ColorMuted).
		Padding(0, 1)
	return style.Render(content)
}

// SuccessBox renders a success message in a styled box.
func SuccessBox(title, message string) {
	style := lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(ColorSuccess).
		Padding(0, 2).
		Width(50)

	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(ColorSuccess)

	content := titleStyle.Render(Symbols.Success+" "+title) + "\n" + message
	fmt.Fprintln(os.Stderr, style.Render(content))
}

// ErrorBox renders an error message in a styled box.
func ErrorBox(title, message string) {
	style := lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(ColorError).
		Padding(0, 2).
		Width(50)

	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(ColorError)

	content := titleStyle.Render(Symbols.Error+" "+title) + "\n" + message
	fmt.Fprintln(os.Stderr, style.Render(content))
}

// Newline prints an empty line.
func Newline() {
	fmt.Fprintln(os.Stderr)
}
