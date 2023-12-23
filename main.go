package main

import (
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	// "github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/bubbles/list"
)

/* style */
var style = lipgloss.NewStyle()
// var italicsStyle = lipgloss.NewStyle().Italic(true)

/* custom item */

type File struct {
	name string
}

func (f File) Title() string {
	return f.name
}

func (f File) Description() string {
	return ""
}

func (f File) FilterValue() string {
	return f.name
}

/* main model */

type Model struct {
	list list.Model
	err error
}

func New() *Model {
	return &Model {}
}

func (m *Model) initList(w, h int) {
	dirEntries, err := os.ReadDir(".")
	if err != nil {
		m.err = err
		return
	}

	var items []list.Item
	for _, dirEntry := range dirEntries {
		items = append(items, File{ dirEntry.Name() })
	}

	m.list = list.New(items, list.NewDefaultDelegate(), w, h)

	m.list.Title = "..."
	m.list.SetShowHelp(false)
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.initList(msg.Width, msg.Height)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	return style.Render(m.list.View())
}

func main() {
	f, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	m := New()

	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
