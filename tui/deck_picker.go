package tui

import (
	"errors"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/filepicker"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/x/ansi"
)

type filePickerModel struct {
	width        int
	height       int
	selectedFile string
	title        string
	err          error
	filepicker   filepicker.Model
}

type clearErrorMsg struct{}

func initDeckPickerModel(directory string) filePickerModel {
	fp := filepicker.New()
	fp.AutoHeight = false
	fp.AllowedTypes = []string{".cod"}
	fp.CurrentDirectory = directory
	fp.KeyMap.Open.SetKeys("enter", " ", "l", "right")
	fp.KeyMap.Select.SetKeys("enter", " ")

	// fp.Styles.Cursor = styleFilePicker
	// fp.Styles.Selected = styleFilePickerSelected
	// fp.Styles.File = styleFilePicker
	return filePickerModel{
		selectedFile: "",
		title:        "Apply Printings",
		filepicker:   fp,
	}
}

func clearErrorAfter(t time.Duration) tea.Cmd {
	return tea.Tick(t, func(_ time.Time) tea.Msg {
		return clearErrorMsg{}
	})
}

func (m filePickerModel) Title() string { return m.title }
func (m filePickerModel) Hidden() bool  { return false }

func (m filePickerModel) CustomInit() submodel {
	return initDeckPickerModel(m.filepicker.CurrentDirectory)
}

func (m filePickerModel) Init() tea.Cmd {
	return m.filepicker.Init()
}

func (m filePickerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if size, ok := msg.(tea.WindowSizeMsg); ok {
		m.width = size.Width
		m.height = size.Height
		m.filepicker.SetHeight(max(m.height-14, 1))
	}

	var cmd tea.Cmd
	m.filepicker, cmd = m.filepicker.Update(msg)

	if didSelect, path := m.filepicker.DidSelectDisabledFile(msg); didSelect {
		m.err = errors.New(path + " is not valid.")
		m.selectedFile = ""
		return m, tea.Batch(cmd, clearErrorAfter(2*time.Second))
	}

	if didSelect, path := m.filepicker.DidSelectFile(msg); didSelect {
		m.selectedFile = path
		return m, tea.Batch(cmd, func() tea.Msg { return deckSelectedMsg{path: m.selectedFile} })
	}

	return m, cmd
}

func (m filePickerModel) View() string {
	var s strings.Builder
	if m.err != nil {
		s.WriteString(styleError.Render(m.err.Error()))
	} else if m.selectedFile == "" {
		s.WriteString("Pick a file:")
	}

	out := s.String() + "\n\n"
	out += stylePanel.Render(ansi.Strip(m.filepicker.View())) + "\n\n"
	out += styleHelp.Render("Press Esc or q to return to the main menu.") + "\n"
	out += styleHelp.Render("Ctrl+C to quit.")
	return out
}
