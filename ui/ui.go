package ui

import (
	"adminmsyql/ui/navigation"
	"adminmsyql/ui/views"
	"adminmsyql/ui/views/rg"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

type KeyMap struct {
	FirstTab  key.Binding
	SecondTab key.Binding
	// more tab
	Quit key.Binding
}

var DefaultKeyMap = KeyMap{
	FirstTab: key.NewBinding(
		// TODO: MOD + 1 or F1 test
		key.WithKeys("1"),
		key.WithHelp("1", "First tab"),
	),
	SecondTab: key.NewBinding(
		key.WithKeys("2"),
		key.WithHelp("2", "Second tab"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+q"),
		key.WithHelp("q/Q", "quit"),
	),
}

type Model struct {
	keymap KeyMap
	nav    navigation.Model
	views  []views.View
}

func NewModel() Model {
	m := Model{
		keymap: DefaultKeyMap,
	}
	m.nav = navigation.NewModel()

	m.views = append(m.views, rg.NewModel())

	return m
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(tea.EnterAltScreen)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := make([]tea.Cmd, 0)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keymap.Quit):
			return m, tea.Quit
		case key.Matches(msg, m.keymap.FirstTab):
			m.nav.NthTab(1)
			return m, nil
		case key.Matches(msg, m.keymap.SecondTab):
			m.nav.NthTab(2)
			return m, nil
		}
	}

	v, cmd := m.views[m.nav.CurrentId].Update(msg)
	m.views[m.nav.CurrentId] = v
	cmds = append(cmds, cmd)

	nav, cmd := m.nav.Update(msg)
	m.nav = nav
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	s := strings.Builder{}
	s.WriteString(m.nav.View() + "\n\n")
	s.WriteString(m.views[m.nav.CurrentId].View())
	return s.String()
}
