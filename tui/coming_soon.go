package tui

import (
	tea "github.com/charmbracelet/bubbletea"
)

type comingSoonModel struct{}

func initComingSoonModel() comingSoonModel {
	return comingSoonModel{}
}

func (m comingSoonModel) Title() string { return "Comming Soon" }
func (m comingSoonModel) Hidden() bool  { return false }

func (m comingSoonModel) Init() tea.Cmd {
	return nil
}

func (m comingSoonModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m comingSoonModel) View() string {
	out := styleTitle.Render("Apply printings") + "\n\n"
	out += styleBase.Render("Coming soon.") + "\n\n"
	out += styleHelp.Render("Press Esc or q to return to the main menu.") + "\n"
	out += styleHelp.Render("Ctrl+C to quit.")
	return out
}
