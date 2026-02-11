package tui

import "github.com/charmbracelet/lipgloss"

const (
	colorRosewater = "#f5e0dc"
	colorFlamingo  = "#f2cdcd"
	colorPink      = "#f5c2e7"
	colorMauve     = "#cba6f7"
	colorRed       = "#f38ba8"
	colorMaroon    = "#eba0ac"
	colorPeach     = "#fab387"
	colorYellow    = "#f9e2af"
	colorGreen     = "#a6e3a1"
	colorTeal      = "#94e2d5"
	colorSky       = "#89dceb"
	colorSapphire  = "#74c7ec"
	colorBlue      = "#89b4fa"
	colorLavender  = "#b4befe"

	colorText     = "#cdd6f4"
	colorSubtext1 = "#bac2de"
	colorSubtext0 = "#a6adc8"
	colorOverlay2 = "#9399b2"
	colorOverlay1 = "#7f849c"
	colorOverlay0 = "#6c7086"
	colorSurface2 = "#585b70"
	colorSurface1 = "#45475a"
	colorSurface0 = "#313244"
	colorBase     = "#1e1e2e"
	colorMantle   = "#181825"
	colorCrust    = "#11111b"
)

var (
	styleBase             = lipgloss.NewStyle().Foreground(lipgloss.Color(colorText))
	styleTitle            = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color(colorLavender))
	styleMuted            = lipgloss.NewStyle().Foreground(lipgloss.Color(colorOverlay0))
	styleAccent           = lipgloss.NewStyle().Foreground(lipgloss.Color(colorSky))
	styleSuccess          = lipgloss.NewStyle().Foreground(lipgloss.Color(colorGreen))
	styleError            = lipgloss.NewStyle().Foreground(lipgloss.Color(colorRed))
	styleMenuItem         = lipgloss.NewStyle().Foreground(lipgloss.Color(colorText))
	styleMenuItemSelected = lipgloss.NewStyle().Foreground(lipgloss.Color(colorText)).Background(lipgloss.Color(colorSurface0)).Bold(true)
	styleHelp             = lipgloss.NewStyle().Foreground(lipgloss.Color(colorSubtext0))
	styleBackground       = lipgloss.NewStyle().Background(lipgloss.Color(colorBase))
	stylePanel            = lipgloss.NewStyle().Background(lipgloss.Color(colorSurface0)).Foreground(lipgloss.Color(colorText)).Padding(1, 2)
	stylePaddingWidth     = 10
	stylePaddingHeight    = 10
)
