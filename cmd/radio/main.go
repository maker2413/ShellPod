package main

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/maker2413/shellpod/internal/config"
	"github.com/maker2413/shellpod/internal/radio"
	"github.com/maker2413/shellpod/internal/tui"
)

func main() {
	config, err := config.GetConfig()
	if err != nil {
		log.Fatal(err)
	}

	if config.Debug {
		f, err := tea.LogToFile("debug.log", "debug")
		if err != nil {
			log.Fatal(err)
		}
		defer func() {
			err = f.Close()
			if err != nil {
				log.Fatal(err)
			}
		}()
	}

	radio, err := radio.NewRadio(config)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		err = radio.Resp.Body.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()
	defer func() {
		err = radio.Streamer.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	m, err := tui.NewModel(radio)
	if err != nil {
		log.Fatal(err)
	}

	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
