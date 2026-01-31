package tui

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

const (
	tickRate = time.Second
)

type tickMsg time.Time

func leftShiftString(s string) string {
	if len(s) <= 1 {
		return s
	}

	b := make([]byte, len(s))
	copy(b, s[1:])
	b[len(s)-1] = s[0]
	return string(b)
}

func tick() tea.Cmd {
	return tea.Tick(tickRate, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}
