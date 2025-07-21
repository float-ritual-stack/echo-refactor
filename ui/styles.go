package ui

import (
	"strings"
	
	"github.com/charmbracelet/lipgloss"
)

var (
	// Colors
	primaryColor   = lipgloss.Color("#00ff00") // Green
	secondaryColor = lipgloss.Color("#888888") // Gray
	accentColor    = lipgloss.Color("#00ffff") // Cyan
	errorColor     = lipgloss.Color("#ff0000") // Red
	bgColor        = lipgloss.Color("#000000") // Black
	fgColor        = lipgloss.Color("#ffffff") // White
	dimColor       = lipgloss.Color("#666666") // Dim gray
	
	// Header styles
	HeaderStyle = lipgloss.NewStyle().
		Foreground(primaryColor).
		Bold(true).
		Padding(0, 1)
	
	// Entry styles
	EntryStyle = lipgloss.NewStyle().
		Foreground(fgColor).
		PaddingLeft(2)
	
	SelectedEntryStyle = lipgloss.NewStyle().
		Border(lipgloss.NormalBorder(), false, false, false, true).
		BorderForeground(primaryColor).
		PaddingLeft(1)
	
	// Type badge styles
	TypeBadgeStyle = lipgloss.NewStyle().
		Foreground(primaryColor).
		Bold(true)
	
	// Timestamp styles
	TimestampStyle = lipgloss.NewStyle().
		Foreground(dimColor).
		MarginRight(1)
	
	// Input styles
	InputStyle = lipgloss.NewStyle().
		Border(lipgloss.NormalBorder(), true, false, false, false).
		BorderForeground(secondaryColor).
		Padding(0, 1)
	
	FocusedInputStyle = lipgloss.NewStyle().
		Border(lipgloss.NormalBorder(), true, false, false, false).
		BorderForeground(primaryColor).
		Padding(0, 1)
	
	// Status bar styles
	StatusBarStyle = lipgloss.NewStyle().
		Foreground(secondaryColor).
		Padding(0, 1)
	
	// Mode indicator styles
	ModeStyle = lipgloss.NewStyle().
		Foreground(accentColor).
		Bold(true)
	
	// Help text styles
	HelpStyle = lipgloss.NewStyle().
		Foreground(dimColor)
	
	// Error styles
	ErrorStyle = lipgloss.NewStyle().
		Foreground(errorColor).
		Bold(true)
	
	// Viewport styles
	ViewportStyle = lipgloss.NewStyle().
		PaddingTop(1).
		PaddingBottom(1)
	
	// Insertion indicator
	InsertionIndicatorStyle = lipgloss.NewStyle().
		Foreground(primaryColor).
		Bold(true)
)

// GetTypeColor returns the color for a given entry type
func GetTypeColor(entryType string) lipgloss.Color {
	typeColors := map[string]lipgloss.Color{
		"ctx":       primaryColor,
		"bridge":    lipgloss.Color("#ff00ff"),
		"mode":      accentColor,
		"project":   lipgloss.Color("#ffff00"),
		"dispatch":  lipgloss.Color("#ff8800"),
		"todo":      lipgloss.Color("#ff0088"),
		"meeting":   lipgloss.Color("#8800ff"),
		"claude":    lipgloss.Color("#0088ff"),
		"error":     errorColor,
		"fix":       lipgloss.Color("#00ff88"),
		"idea":      lipgloss.Color("#ff88ff"),
		"highlight": lipgloss.Color("#ffff88"),
		"note":      lipgloss.Color("#88ffff"),
		"log":       secondaryColor,
	}
	
	if color, ok := typeColors[entryType]; ok {
		return color
	}
	return fgColor
}

// RenderTypePrefix renders a type prefix with appropriate styling
func RenderTypePrefix(entryType string) string {
	if entryType == "log" || entryType == "text" || entryType == "" {
		return ""
	}
	
	color := GetTypeColor(entryType)
	style := lipgloss.NewStyle().
		Foreground(color).
		Bold(true)
	
	return style.Render(entryType+"::") + " "
}

// RenderIndent creates indentation string
func RenderIndent(level int) string {
	return strings.Repeat("  ", level)
}