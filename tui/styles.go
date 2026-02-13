package tui

import (
	catppuccin "github.com/catppuccin/go"
	"github.com/charmbracelet/lipgloss"
)

var flavor = catppuccin.Macchiato

var (
	styleBase    = lipgloss.NewStyle().Foreground(lipgloss.Color(flavor.Base().Hex))
	styleTitle   = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color(flavor.Mauve().Hex))
	styleAccent  = lipgloss.NewStyle().Foreground(lipgloss.Color(flavor.Teal().Hex))
	styleSuccess = lipgloss.NewStyle().Foreground(lipgloss.Color(flavor.Green().Hex))
	styleError   = lipgloss.NewStyle().Foreground(lipgloss.Color(flavor.Red().Hex))
	styleHelp    = lipgloss.NewStyle().Foreground(lipgloss.Color(flavor.Overlay1().Hex))
	stylePanel   = lipgloss.NewStyle().Background(lipgloss.Color(flavor.Base().Hex)).Foreground(lipgloss.Color(flavor.Text().Hex)).Padding(1, 2)
)
