package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

var (
	lastID int
	idMtx  sync.Mutex
)

func nextID() int {
	idMtx.Lock()
	defer idMtx.Unlock()
	lastID++
	return lastID
}

type readDirMsg struct {
	id    int
	files []os.DirEntry
}

type Model struct {
	files     []os.DirEntry
	currDir   string
	maxHeight int
	idx       int
	keys      KeyMap
	styles    Styles
	id        int
}

func New() Model {
	return Model{
		currDir:   ".",
		maxHeight: 0,
		idx:       0,
		keys:      DefaultKeyMap(),
		styles:    DefaultStyles(),
		id:        nextID(),
	}
}

func (m Model) readDir(path string) tea.Cmd {
	return func() tea.Msg {
		dirEntries, err := os.ReadDir(path)
		if err != nil {
			return err
		}
		// sort?
		return readDirMsg{id: m.id, files: dirEntries}
	}
}

func (m Model) Init() tea.Cmd {
	return m.readDir(m.currDir)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.maxHeight = msg.Height
	case readDirMsg:
		if msg.id != m.id {
			break
		}
		m.files = msg.files
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Quit):
			return m, tea.Quit
		case key.Matches(msg, m.keys.Up):
			m.idx--
			if m.idx < 0 {
				m.idx = 0
			}
		case key.Matches(msg, m.keys.Down):
			m.idx++
			if m.idx >= len(m.files) {
				m.idx = len(m.files) - 1
			}
		}
	}

	return m, nil
}

func (m Model) View() string {
	currPath, err := filepath.Abs(m.currDir)
	if err != nil {
		currPath = fmt.Sprintf("Error displaying absolute path: %s", err)
	}

	currPath = m.styles.Path.Render(currPath)

	if len(m.files) == 0 {
		return currPath + "\n\n" + m.styles.EmptyDir.String()
	}

	isRoot := m.currDir == "/"
	if !isRoot {
		currPath += m.styles.Path.Render("/")
	}

	var files string
	var hovered string
	for i, f := range m.files {
		info, err := f.Info()
		if err != nil {
			files += fmt.Sprintf("Error reading file info: %s\n", err)
			continue
		}
		isSymlink := info.Mode()&os.ModeSymlink != 0

		file := f.Name()
		if i == m.idx {
			hovered = m.styles.PathEnd.Render(file)
			file = m.styles.Hover.Render(file)
		}

		if f.IsDir() {
			file = m.styles.Directory.Render(file)
		} else if isSymlink {
			file = m.styles.Symlink.Render(file)
		}

		files += file + "\n"
	}

	return currPath + hovered + "\n\n" + files
}

func main() {
	if _, err := tea.NewProgram(New()).Run(); err != nil {
		log.Fatal(err)
	}
}
