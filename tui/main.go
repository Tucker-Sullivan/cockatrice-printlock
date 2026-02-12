package tui

import (
	"fmt"
	"strings"

	"github.com/Tucker-Sullivan/cockatrice-printlock/internal/cockatrice/cardsdb"
	"github.com/Tucker-Sullivan/cockatrice-printlock/internal/config"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	selected         bool
	width            int
	height           int
	numScreens       int
	numHiddenScreens int
	screen           int
	screens          []submodel
	selectedDeck     string
	selectedSets     []string
	config           *config.Config
	db               *cardsdb.CardsDB
}

type submodel interface {
	Title() string
	Hidden() bool
	Init() tea.Cmd
	Update(tea.Msg) (tea.Model, tea.Cmd)
	View() string
}

type deckSelectedMsg struct {
	path string
}

type setsSelectedMsg struct {
	sets []string
}

func renderScreen(content string, width, height int) string {
	if width <= 0 || height <= 0 {
		return content
	}
	return lipgloss.Place(
		width,
		height,
		lipgloss.Center,
		lipgloss.Center,
		content,
		lipgloss.WithWhitespaceBackground(lipgloss.Color(colorBase)),
		lipgloss.WithWhitespaceChars(" "),
	)
}

func (m model) updateSubmodel(msg tea.Msg) (tea.Model, tea.Cmd) {
	updated, cmd := m.screens[m.screen].Update(msg)
	if screen, ok := updated.(submodel); ok {
		m.screens[m.screen] = screen
	}
	return m, cmd
}

func InitModel() (model, error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		return model{}, err
	}

	db, err := cardsdb.ParseCardsDB(cfg.CardsXMLPath)
	if err != nil {
		return model{}, err
	}

	deckPicker, err := initDeckPickerModel(cfg.DeckDir)
	if err != nil {
		return model{}, err
	}

	screens := []submodel{
		initComingSoonModel(),
		deckPicker,
	}

	hiddenScreens := []submodel{
		initSetSelectionModel(db),
		initApplyConfirmationModel(),
	}

	return model{
		numScreens:       len(screens),
		numHiddenScreens: len(hiddenScreens),
		screens:          append(screens, hiddenScreens...),
		config:           cfg,
		db:               db,
	}, nil
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		if m.selected {
			return m.updateSubmodel(msg)
		}
	case deckSelectedMsg:
		m.selectedDeck = msg.path
		m.screen = 2 // set selection screen
		return m, tea.Batch(
			m.screens[m.screen].Init(),
			func() tea.Msg { return tea.WindowSizeMsg{Width: m.width, Height: m.height} },
		)
	case setsSelectedMsg:
		m.selectedSets = msg.sets
		m.screen = 3 // apply confirmation screen
		if confirmation, ok := m.screens[m.screen].(applyConfirmationModel); ok {
			confirmation.deckPath = m.selectedDeck
			confirmation.sets = m.selectedSets
			m.screens[m.screen] = confirmation
		}
		return m, tea.Batch(
			m.screens[m.screen].Init(),
			func() tea.Msg { return tea.WindowSizeMsg{Width: m.width, Height: m.height} },
		)
	case tea.KeyMsg:
		if msg.Type == tea.KeyCtrlC {
			return m, tea.Quit
		}
		if m.selected {
			if msg.String() == "q" {
				m.selected = false
				return m, nil
			} else {
				return m.updateSubmodel(msg)
			}
		}

		switch msg.String() {
		case "q":
			return m, tea.Quit
		case "up", "k":
			if m.screen > 0 {
				m.screen--
			}
		case "down", "j":

			if m.screen < m.numScreens-1 {
				m.screen++
			}
		case "enter", " ":
			m.selected = true
			return m, tea.Batch(
				m.screens[m.screen].Init(),
				func() tea.Msg {
					return tea.WindowSizeMsg{Width: m.width, Height: m.height}
				},
			)
		}
	default:
		if m.selected {
			return m.updateSubmodel(msg)
		}
	}

	return m, nil
}

func (m model) View() string {
	var s strings.Builder

	if m.selected {
		s.WriteString(styleTitle.Render(m.screens[m.screen].Title()) + "\n\n")
		s.WriteString(m.screens[m.screen].View())
	} else {
		s.WriteString(styleTitle.Render("Printlock") + "\n\n")
		s.WriteString(styleBase.Render("Select a process:") + "\n\n")
		for i, item := range m.screens {
			if item.Hidden() {
				continue
			}
			cursor := " "
			if i == m.screen {
				cursor = ">"
			}
			fmt.Fprintf(&s, "%s %s\n", cursor, item.Title())
		}
		s.WriteString("\n" + styleHelp.Render("Use Up/Down or j/k to move, Enter/Space to select.") + "\n")
		s.WriteString(styleHelp.Render("Ctrl+C to quit."))
	}

	return renderScreen(stylePanel.Render(s.String()), m.width, m.height)
}
