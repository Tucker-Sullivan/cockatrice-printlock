package tui

import (
	catppuccin "github.com/catppuccin/go"
	"github.com/charmbracelet/lipgloss"
)

var flavor = catppuccin.Macchiato

var (
	styleBase             = lipgloss.NewStyle().Foreground(lipgloss.Color(flavor.Base().Hex))
	styleTitle            = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color(flavor.Mauve().Hex))
	styleMuted            = lipgloss.NewStyle().Foreground(lipgloss.Color(flavor.Overlay1().Hex))
	styleAccent           = lipgloss.NewStyle().Foreground(lipgloss.Color(flavor.Mauve().Hex))
	styleSuccess          = lipgloss.NewStyle().Foreground(lipgloss.Color(flavor.Green().Hex))
	styleError            = lipgloss.NewStyle().Foreground(lipgloss.Color(flavor.Red().Hex))
	styleMenuItem         = lipgloss.NewStyle().Foreground(lipgloss.Color(flavor.Mauve().Hex))
	styleMenuItemSelected = lipgloss.NewStyle().Foreground(lipgloss.Color(flavor.Base().Hex)).Background(lipgloss.Color(flavor.Teal().Hex)).Bold(true)
	styleHelp             = lipgloss.NewStyle().Foreground(lipgloss.Color(flavor.Overlay1().Hex))
	styleBackground       = lipgloss.NewStyle().Background(lipgloss.Color(flavor.Base().Hex))
	stylePanel            = lipgloss.NewStyle().Background(lipgloss.Color(flavor.Base().Hex)).Foreground(lipgloss.Color(flavor.Text().Hex)).Padding(1, 2)
)
