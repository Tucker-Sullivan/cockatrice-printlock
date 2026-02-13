package tui

import (
	"fmt"
	"io"
	"sort"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/Tucker-Sullivan/cockatrice-printlock/internal/cockatrice/cardsdb"
	"github.com/Tucker-Sullivan/cockatrice-printlock/internal/config"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
)

type item struct {
	title, desc string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title + "\t" + i.desc }

type setDelegate struct {
	list.DefaultDelegate
	selected map[string]bool
}

type setsSelectionModel struct {
	globalSetsHavePriority bool
	width                  int
	height                 int
	list                   list.Model
	selected               map[string]bool
	selectedOrder          []string
	globalSelectedOrder    []string
	globalSetsMessage      string
	cfg                    *config.Config
	db                     *cardsdb.CardsDB
}

func filterByCodeOrDesc(term string, targets []string) []list.Rank {
	term = strings.TrimSpace(term)
	useCode := isAllUpper(term)
	filterTargets := make([]string, len(targets))
	for i, target := range targets {
		code, desc := splitFilterTarget(target)
		if useCode {
			filterTargets[i] = code
		} else {
			filterTargets[i] = desc
		}
	}
	return list.DefaultFilter(term, filterTargets)
}

func splitFilterTarget(value string) (string, string) {
	parts := strings.SplitN(value, "\t", 2)
	if len(parts) == 2 {
		return parts[0], parts[1]
	}
	return value, value
}

func isAllUpper(term string) bool {
	hasLetter := false
	for _, r := range term {
		if unicode.IsLetter(r) {
			hasLetter = true
			if !unicode.IsUpper(r) {
				return false
			}
		}
	}
	return hasLetter
}

func newSetDelegate(selected map[string]bool) setDelegate {
	return setDelegate{
		DefaultDelegate: list.NewDefaultDelegate(),
		selected:        selected,
	}
}

func (d setDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	it, ok := listItem.(item)
	if !ok {
		return
	}
	if m.Width() <= 0 {
		return
	}
	checked := " "
	if d.selected[it.title] {
		checked = "x"
	}

	title := fmt.Sprintf("[%s] %s", checked, it.title)
	desc := it.desc
	prefixLen := utf8.RuneCountInString("[ ] ")

	textWidth := m.Width() - d.Styles.NormalTitle.GetPaddingLeft() - d.Styles.NormalTitle.GetPaddingRight()
	title = ansi.Truncate(title, textWidth, "...")
	if d.ShowDescription {
		var lines []string
		for i, line := range strings.Split(desc, "\n") {
			if i >= d.Height()-1 {
				break
			}
			lines = append(lines, ansi.Truncate(line, textWidth, "..."))
		}
		desc = strings.Join(lines, "\n")
	}

	isSelected := index == m.Index()
	emptyFilter := m.FilterState() == list.Filtering && m.FilterValue() == ""
	isFiltered := m.FilterState() == list.Filtering || m.FilterState() == list.FilterApplied

	var matchedRunes []int
	if isFiltered && isAllUpper(strings.TrimSpace(m.FilterValue())) {
		matchedRunes = m.MatchesForItem(index)
		if len(matchedRunes) > 0 {
			shifted := make([]int, len(matchedRunes))
			for i, r := range matchedRunes {
				shifted[i] = r + prefixLen
			}
			matchedRunes = shifted
		}
	}

	if emptyFilter {
		title = d.Styles.DimmedTitle.Render(title)
		desc = d.Styles.DimmedDesc.Render(desc)
	} else if isSelected && m.FilterState() != list.Filtering {
		if len(matchedRunes) > 0 {
			unmatched := d.Styles.SelectedTitle.Inline(true)
			matched := unmatched.Inherit(d.Styles.FilterMatch)
			title = lipgloss.StyleRunes(title, matchedRunes, matched, unmatched)
		}
		title = d.Styles.SelectedTitle.Render(title)
		desc = d.Styles.SelectedDesc.Render(desc)
	} else {
		if len(matchedRunes) > 0 {
			unmatched := d.Styles.NormalTitle.Inline(true)
			matched := unmatched.Inherit(d.Styles.FilterMatch)
			title = lipgloss.StyleRunes(title, matchedRunes, matched, unmatched)
		}
		title = d.Styles.NormalTitle.Render(title)
		desc = d.Styles.NormalDesc.Render(desc)
	}

	if d.ShowDescription {
		fmt.Fprintf(w, "%s\n%s", title, desc) //nolint:errcheck
		return
	}
	fmt.Fprintf(w, "%s", title) //nolint:errcheck
}

func selectedLabel(selected []string) string {
	if len(selected) == 0 {
		return "Selected sets: (none)"
	}
	return "Selected sets: " + strings.Join(selected, ", ")
}

func (s setsSelectionModel) combineSets() []string {
	if s.globalSetsHavePriority {
		return append(s.globalSelectedOrder, s.selectedOrder...)
	} else {
		return append(s.selectedOrder, s.globalSelectedOrder...)
	}
}

func initSetSelectionModel(db *cardsdb.CardsDB, cfg *config.Config) setsSelectionModel {
	setItems := make([]item, 0, len(db.SetsByCode))
	for _, v := range db.SetsByCode {
		setItems = append(setItems, item{title: v.Code, desc: v.LongName})
	}
	sort.Slice(setItems, func(i, j int) bool {
		return setItems[i].title < setItems[j].title
	})
	items := make([]list.Item, 0, len(setItems))
	for _, it := range setItems {
		items = append(items, it)
	}

	selectedMap := map[string]bool{}
	for _, item := range cfg.GlobalSetPriority {
		selectedMap[item] = true
	}

	m1 := "(global sets "
	m2 := "don't "
	m3 := "have priority)"
	if cfg.GlobalSetsHavePriority {
		m2 = ""
	}

	d := newSetDelegate(selectedMap)
	m := setsSelectionModel{
		globalSetsHavePriority: cfg.GlobalSetsHavePriority,
		list:                   list.New(items, d, 1, 1),
		selected:               selectedMap,
		selectedOrder:          make([]string, 0),
		globalSelectedOrder:    cfg.GlobalSetPriority,
		globalSetsMessage:      m1 + m2 + m3,
		cfg:                    cfg,
		db:                     db,
	}

	m.list.ResetFilter()
	m.list.ResetSelected()
	m.list.SetShowPagination(false)
	m.list.Title = "Sets"
	m.list.Filter = filterByCodeOrDesc
	return m
}

func (s setsSelectionModel) Title() string { return "Set Selection" }
func (s setsSelectionModel) Hidden() bool  { return true }

func (s setsSelectionModel) CustomInit() submodel {
	return initSetSelectionModel(s.db, s.cfg)
}

func (s setsSelectionModel) Init() tea.Cmd {
	return nil
}

func (s setsSelectionModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		s.width = msg.Width
		s.height = msg.Height
		s.list.SetSize(min(s.width, 50), max(s.height-10, 1))
	case tea.KeyMsg:
		switch msg.String() {
		case " ":
			if s.list.SettingFilter() {
				break
			}
			if it, ok := s.list.SelectedItem().(item); ok {
				if s.selected == nil {
					s.selected = map[string]bool{}
				}
				if !s.selected[it.title] {
					s.selected[it.title] = true
					s.selectedOrder = append(s.selectedOrder, it.title)
				} else {
					s.selected[it.title] = false
					for i, v := range s.selectedOrder {
						if v == it.title {
							s.selectedOrder = append(s.selectedOrder[:i], s.selectedOrder[i+1:]...)
						}
					}
				}
			}
		case "enter":
			if s.list.FilterState() != list.Filtering {
				return s, func() tea.Msg { return setsSelectedMsg{sets: s.combineSets()} }
			}
		}
	}
	var cmd tea.Cmd
	s.list, cmd = s.list.Update(msg)
	return s, cmd
}

func (s setsSelectionModel) View() string {
	var out strings.Builder
	out.WriteString(styleAccent.Render(selectedLabel(s.combineSets())) + "\n")
	out.WriteString(styleAccent.Render(s.globalSetsMessage) + "\n")
	out.WriteString(stylePanel.Render(ansi.Strip(s.list.View())))
	return out.String()
}
