package ui

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	// Primary colors
	ColorPrimary   = lipgloss.Color("14")  // Cyan
	ColorSuccess   = lipgloss.Color("10")  // Green
	ColorWarning   = lipgloss.Color("11")  // Yellow
	ColorError     = lipgloss.Color("9")   // Red
	ColorSecondary = lipgloss.Color("5")   // Magenta
	ColorDark      = lipgloss.Color("8")   // Gray

	// Title style
	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorPrimary).
			Background(lipgloss.Color("0")).
			Padding(1, 2).
			Width(50).
			Align(lipgloss.Center)

	// Menu item styles
	MenuItemStyle = lipgloss.NewStyle().
			Padding(0, 2).
			Foreground(lipgloss.Color("15"))

	SelectedItemStyle = lipgloss.NewStyle().
				Padding(0, 2).
				Foreground(ColorSuccess).
				Bold(true)

	// Button styles
	ButtonStyle = lipgloss.NewStyle().
			Padding(0, 2).
			Foreground(ColorPrimary).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorPrimary)

	SuccessStyle = lipgloss.NewStyle().
			Foreground(ColorSuccess).
			Bold(true)

	ErrorStyle = lipgloss.NewStyle().
			Foreground(ColorError).
			Bold(true)

	WarningStyle = lipgloss.NewStyle().
			Foreground(ColorWarning)

	HelpStyle = lipgloss.NewStyle().
			Foreground(ColorDark).
			Italic(true)

	// Input style
	InputStyle = lipgloss.NewStyle().
			Padding(0, 1).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorSecondary)
)

func RenderBox(title, content string) string {
	return lipgloss.NewStyle().
		Padding(1, 2).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(ColorPrimary).
		Render(
			lipgloss.JoinVertical(
				lipgloss.Left,
				TitleStyle.Copy().Width(0).Padding(0).Render(title),
				content,
			),
		)
}
