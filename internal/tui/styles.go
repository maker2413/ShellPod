package tui

import "github.com/charmbracelet/lipgloss"

var radioStyle = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).
	BorderForeground(lipgloss.Color("#FFFFFF")).Padding(1, 2)

func (m *model) centeredStyle(contents string) string {
	return lipgloss.Place(
		m.width, m.height, lipgloss.Center, lipgloss.Center, radioStyle.Render(contents))
}
