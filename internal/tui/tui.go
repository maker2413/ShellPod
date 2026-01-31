package tui

import (
	"log"
	"sync"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/gopxl/beep"
	"github.com/gopxl/beep/effects"
	"github.com/gopxl/beep/speaker"
	"github.com/maker2413/shellpod/internal/stack"
)

const (
	titlePadding = "        "
)

type page int

const (
	homePage page = iota
	playerPage
)

type changePageMsg struct {
	next page
}

type model struct {
	sampleRate          beep.SampleRate
	streamer            beep.StreamSeeker
	volume              *effects.Volume
	volumeMutex         sync.Mutex
	stationName         string
	titleChan           <-chan string
	titleMutex          sync.Mutex
	currentTitle        string
	displayedTitle      string
	maxDisplayTitleSize int
	currentPage         page
	pageHistory         stack.Stack[page]
	width               int
	height              int
}

func NewModel(
	sampleRate beep.SampleRate,
	streamer beep.StreamSeeker,
	stationName string,
	titleChan <-chan string,
	maxDisplayedTitleSize int,
) (tea.Model, error) {
	volume := &effects.Volume{Streamer: streamer, Base: 2, Volume: -2.0}

	if maxDisplayedTitleSize <= 0 {
		maxDisplayedTitleSize = len(titlePadding)
	}

	return &model{sampleRate: sampleRate,
		streamer:            streamer,
		volume:              volume,
		stationName:         stationName,
		titleChan:           titleChan,
		currentTitle:        "",
		displayedTitle:      "",
		maxDisplayTitleSize: maxDisplayedTitleSize,
		currentPage:         homePage,
		pageHistory:         stack.Stack[page]{},
	}, nil
}

func (m *model) Init() tea.Cmd {
	return tick()
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			speaker.Close()
			return m, tea.Quit
		case "esc":
			m.Stop()
			if m.pageHistory.IsEmpty() {
				speaker.Close()
				return m, tea.Quit
			} else {
				previousPage, err := m.pageHistory.Pop()
				if err != nil {
					log.Fatal(err)
				}
				m.currentPage = previousPage
				return m, nil
			}
		}
	case changePageMsg:
		m.pageHistory.Push(m.currentPage)
		m.currentPage = msg.next

		return m, tick()
	}

	switch m.currentPage {
	case homePage:
		return m.homeUpdate(msg)
	case playerPage:
		return m.playerUpdate(msg)
	}

	return m, nil
}

func (m *model) View() string {
	switch m.currentPage {
	case homePage:
		return m.homeView()
	case playerPage:
		return m.playerView()
	default:
		return m.homeView()
	}
}
