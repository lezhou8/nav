package main

import (
	"fmt"
	"io"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	// "github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
)

/* rendering */

type itemDelegate struct{}

func (d itemDelegate) Height() int                               { return 1 }
func (d itemDelegate) Spacing() int                              { return 0 }
func (d itemDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }

func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(File)
	if !ok {
		return
	}

	str := string(i.name)
	style := lipgloss.NewStyle().Padding(0, 2)

	if i.fileType == dir {
		style = style.Foreground(lipgloss.Color("12"))
	} else if i.fileType == symlink {
		style = style.Foreground(lipgloss.Color("32"))
	}

	if index == m.Index() {
		style = style.Foreground(lipgloss.Color("11")).Bold(true)
	}
	if i.isSelected {
		style = style.Italic(true)
	}

	fn := style.Render
	fmt.Fprintf(w, fn(str))
}

/* custom item */

type fileType uint

const (
	regular fileType = iota
	dir
	symlink
)

type File struct {
	name       string
	fileType   fileType
	isSelected bool
}

func (f File) FilterValue() string {
	return f.name
}

/* main model */

type Model struct {
	list list.Model
	err  error
}

func New() *Model {
	var m Model
	m.initList()
	return &m
}

func (m *Model) initList() {
	dirEntries, err := os.ReadDir(".")
	if err != nil {
		m.err = err
		return
	}

	var items []list.Item
	for _, dirEntry := range dirEntries {
		var f File
		f.name = dirEntry.Name()
		if dirEntry.IsDir() {
			f.fileType = dir
		} else if (dirEntry.Type() & os.ModeSymlink) != 0 {
			f.fileType = symlink
		} else {
			f.fileType = regular
		}
		f.isSelected = false
		items = append(items, f)
	}

	m.list = list.New(items, itemDelegate{}, 0, 0)
	m.list.SetShowHelp(false)
	m.list.SetShowStatusBar(false)
	m.list.Styles.Title = lipgloss.NewStyle().Bold(true)

	currDir, err := os.Getwd()
	if err != nil {
		m.list.Title = fmt.Sprintf("%s", err)
		return
	}
	m.list.Title = string(currDir)
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
		m.list.SetHeight(msg.Height)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	return m.list.View()
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
