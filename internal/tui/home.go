package tui

import tea "github.com/charmbracelet/bubbletea"

func (m *model) homeUpdate(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg.(type) {
	case tea.KeyMsg:
		m.Play()

		return m, func() tea.Msg { return changePageMsg{playerPage} }
	}

	return m, nil
}

func (m *model) homeView() string {
	return "♫ Welcome to ShellPod! ♫" +
		"\nPress any key to continue" +
		"\nOr press escape key to exit"
}
