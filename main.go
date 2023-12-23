package main

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
				case "q", "ctrl+c":
					return m, tea.Quit
				case "up", "k":
				case "down", "j":
			}
	}
	return m, nil
}

func (m Model) View() string {
	return ""
}

func main() {
	f, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	m := Model {
	}

	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
