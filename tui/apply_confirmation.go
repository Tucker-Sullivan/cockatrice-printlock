package tui

import (
	"errors"
	"fmt"
	"strings"

	"github.com/Tucker-Sullivan/cockatrice-printlock/internal/cockatrice/cardsdb"
	"github.com/Tucker-Sullivan/cockatrice-printlock/internal/cockatrice/deck"
	"github.com/Tucker-Sullivan/cockatrice-printlock/internal/config"
	"github.com/Tucker-Sullivan/cockatrice-printlock/internal/printlock"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type applyConfirmationModel struct {
	deckPath string
	sets     []string
	cfg      *config.Config
	db       *cardsdb.CardsDB
	err      error
	viewport viewport.Model
}

type applyFinishedMsg struct {
	err error
}

func initApplyConfirmationModel(cfg *config.Config, db *cardsdb.CardsDB) applyConfirmationModel {
	return applyConfirmationModel{
		cfg: cfg,
		db:  db,
	}
}

func (m applyConfirmationModel) Title() string { return "Apply Confirmation" }
func (m applyConfirmationModel) Hidden() bool  { return true }

func (m applyConfirmationModel) CustomInit() submodel {
	return m
}

func (m applyConfirmationModel) Init() tea.Cmd {
	return func() tea.Msg {
		opt := printlock.ApplyOptions{
			PreferHigherCollectionNumber: m.cfg.PreferHigherCollectionNumber,
			SetPriority:                  m.sets,
		}
		d, err := deck.LoadDeck(m.deckPath)
		if err != nil {
			return applyFinishedMsg{err: err}
		}

		d, errs := printlock.ApplyPrintings(m.db, d, opt)

		// write deck
		if err := d.WriteDeckFile(m.deckPath); err != nil {
			return applyFinishedMsg{err: err}
		}

		if len(errs) > 0 {
			var b strings.Builder
			for _, e := range errs {
				b.WriteString(e.Error())
				b.WriteByte('\n')
			}
			return applyFinishedMsg{err: errors.New(b.String())}
		}

		return applyFinishedMsg{err: nil}
	}
}

func (m applyConfirmationModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case applyFinishedMsg:
		m.err = msg.err
		if m.err != nil {
			m.viewport.SetContent(m.err.Error())
		}
		return m, nil
	case tea.WindowSizeMsg:
		m.viewport = viewport.New(msg.Width, msg.Height-8)
		if m.err != nil {
			m.viewport.SetContent(m.err.Error())
		}
	}

	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m applyConfirmationModel) View() string {
	out := styleTitle.Render("Deck: "+m.deckPath) + "\n"
	out += styleHelp.Render("Sets: "+strings.Join(m.sets, ",")) + "\n\n"
	if m.err != nil {
		out += fmt.Sprintf("%s\n%s\n%s\n\n", m.headerView(), styleError.Render(m.viewport.View()), m.footerView())
	} else {
		out += styleSuccess.Render("apply completed successfully")
	}
	out += styleHelp.Render("Press Esc or q to return to the main menu.") + "\n"
	out += styleHelp.Render("Ctrl+C to quit.")
	return out
}

func (m applyConfirmationModel) headerView() string {
	title := styleTitle.Render("Apply completed with issues")
	line := strings.Repeat("─", max(0, m.viewport.Width-lipgloss.Width(title)))
	return lipgloss.JoinHorizontal(lipgloss.Center, title, line)
}

func (m applyConfirmationModel) footerView() string {
	info := styleAccent.Render(fmt.Sprintf("%3.f%%", m.viewport.ScrollPercent()*100))
	line := strings.Repeat("─", max(0, m.viewport.Width-lipgloss.Width(info)))
	return lipgloss.JoinHorizontal(lipgloss.Center, line, info)
}
