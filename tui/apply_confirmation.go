package tui

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type applyConfirmationModel struct {
	deckPath string
	sets     []string
}

func initApplyConfirmationModel() applyConfirmationModel {
	return applyConfirmationModel{}
}

func (m applyConfirmationModel) Title() string { return "Apply Confirmation" }
func (m applyConfirmationModel) Hidden() bool  { return true }

func (m applyConfirmationModel) Init() tea.Cmd {
	return nil
}

func (m applyConfirmationModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m applyConfirmationModel) View() string {
	out := styleTitle.Render("Deck: "+m.deckPath) + "\n"
	out += styleHelp.Render("Sets: "+strings.Join(m.sets, ",")) + "\n\n"
	out += styleHelp.Render("Press Esc or q to return to the main menu.") + "\n"
	out += styleHelp.Render("Ctrl+C to quit.")
	return out
}
