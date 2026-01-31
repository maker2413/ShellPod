package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gopxl/beep/speaker"
)

func (m *model) playerUpdate(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tickMsg:
		m.titleMutex.Lock()
		if len(m.displayedTitle)-len(titlePadding) > m.maxDisplayTitleSize {
			m.displayedTitle = leftShiftString(m.displayedTitle)
		}
		m.titleMutex.Unlock()

		if m.titleUpdate() {
			return m, tea.Batch(
				tick(),
				tea.SetWindowTitle("♫ "+m.stationName+" ~ "+m.currentTitle+" ♫"),
			)
		}

		return m, tick()
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			speaker.Close()
			return m, tea.Quit
		case "w":
			speaker.Lock()
			m.volumeMutex.Lock()
			m.volume.Volume += 0.1
			m.volumeMutex.Unlock()
			speaker.Unlock()
		case "s":
			speaker.Lock()
			m.volumeMutex.Lock()
			m.volume.Volume -= 0.1
			m.volumeMutex.Unlock()
			speaker.Unlock()
		case "m", " ":
			speaker.Lock()
			m.volumeMutex.Lock()
			m.volume.Silent = !m.volume.Silent
			m.volumeMutex.Unlock()
			speaker.Unlock()
		}
	}

	return m, nil
}

func (m *model) playerView() string {
	m.titleMutex.Lock()
	title := m.displayedTitle
	m.titleMutex.Unlock()
	if len(title) > m.maxDisplayTitleSize {
		title = title[:m.maxDisplayTitleSize]
	}

	output := "Station: " + m.stationName +
		"\nSong: " + title +
		"\n\nPress Space or M to mute" +
		"\nUse W and S to control Volume" +
		"\nTo exit press esc, or Ctrl+c"

	style := lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#FFFFFF")).Padding(1, 2)

	return lipgloss.Place(
		m.width, m.height, lipgloss.Center, lipgloss.Center, style.Render(output))
}

func (m *model) titleUpdate() bool {
	select {
	case title := <-m.titleChan:
		if len(title) > 0 {
			m.titleMutex.Lock()
			m.currentTitle = title
			m.displayedTitle = m.currentTitle + titlePadding
			m.titleMutex.Unlock()
		}

		return true
	default:
		return false
	}
}

func (m *model) Play() {
	speaker.Play(m.volume)
}

func (m *model) Stop() {
	speaker.Clear()
}
