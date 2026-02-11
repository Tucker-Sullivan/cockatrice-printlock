package tui

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type swatch struct {
	name  string
	hex   string
	light bool
}

const (
	swatchWidth  = 18
	swatchHeight = 5
	swatchGap    = 2
)

func StyleSwatchesView(width int) string {
	swatches := []swatch{
		{"Rosewater", colorRosewater, true},
		{"Flamingo", colorFlamingo, true},
		{"Pink", colorPink, true},
		{"Mauve", colorMauve, true},
		{"Red", colorRed, true},
		{"Maroon", colorMaroon, true},
		{"Peach", colorPeach, true},
		{"Yellow", colorYellow, true},
		{"Green", colorGreen, true},
		{"Teal", colorTeal, true},
		{"Sky", colorSky, true},
		{"Sapphire", colorSapphire, true},
		{"Blue", colorBlue, true},
		{"Lavender", colorLavender, true},
		{"Text", colorText, true},
		{"Subtext1", colorSubtext1, true},
		{"Subtext0", colorSubtext0, true},
		{"Overlay2", colorOverlay2, true},
		{"Overlay1", colorOverlay1, true},
		{"Overlay0", colorOverlay0, false},
		{"Surface2", colorSurface2, false},
		{"Surface1", colorSurface1, false},
		{"Surface0", colorSurface0, false},
		{"Base", colorBase, false},
		{"Mantle", colorMantle, false},
		{"Crust", colorCrust, false},
	}

	cols := max(swatchesColumns(width), 1)

	rows := make([]string, 0, (len(swatches)+cols-1)/cols)
	for i := 0; i < len(swatches); i += cols {
		end := min(i+cols, len(swatches))
		row := make([]string, 0, end-i)
		for _, s := range swatches[i:end] {
			row = append(row, renderSwatch(s))
		}
		rows = append(rows, lipgloss.JoinHorizontal(lipgloss.Top, row...))
	}
	return strings.Join(rows, "\n")
}

func renderSwatch(s swatch) string {
	fg := colorText
	if s.light {
		fg = colorCrust
	}
	style := lipgloss.NewStyle().
		Background(lipgloss.Color(s.hex)).
		Foreground(lipgloss.Color(fg)).
		Width(swatchWidth).
		Height(swatchHeight).
		Align(lipgloss.Center).
		AlignVertical(lipgloss.Center)
	return style.MarginRight(swatchGap).Render(s.hex)
}

func swatchesColumns(width int) int {
	if width <= 0 {
		return 1
	}
	cellWidth := swatchWidth + swatchGap
	return width / cellWidth
}

type swatchesModel struct {
	width  int
	height int
}

func NewSwatchesModel() swatchesModel {
	return swatchesModel{}
}

func (m swatchesModel) Init() tea.Cmd {
	return nil
}

func (m swatchesModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if size, ok := msg.(tea.WindowSizeMsg); ok {
		m.width = size.Width
		m.height = size.Height
	}
	if key, ok := msg.(tea.KeyMsg); ok {
		switch key.String() {
		case "ctrl+c", "q", "esc":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m swatchesModel) View() string {
	out := styleTitle.Render("Catppuccin Swatches") + "\n\n"
	out += StyleSwatchesView(m.width)
	out += "\n\n" + styleHelp.Render("Press Esc, q, or Ctrl+C to exit.")
	return renderScreen(stylePanel.Render(out), m.width, m.height)
}
