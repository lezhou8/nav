package main

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	filePath string
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
				case "q", "ctrl+c":
					return m, tea.Quit
			}
	}
	return m, nil
}

func (m model) View() string {
	return m.filePath
}

func main() {
	f, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	m := model{
		filePath: "...",
	}

	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
