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

const HeightBuffer int = 5

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

type stack []int

func (s *stack) push(i int) {
	*s = append(*s, i)
}

func (s *stack) pop() int {
	if len(*s) == 0 {
		return -1
	}
	i := (*s)[len(*s)-1]
	*s = (*s)[:len(*s)-1]
	return i
}

type Model struct {
	files      []os.DirEntry
	currDir    string
	maxHeight  int
	idx        int
	keys       KeyMap
	styles     Styles
	min        int
	max        int
	pageDist   int
	halfDist   int
	showHidden bool
	lastFile   string
	stack      stack
	newIdx     int
	lastIdx    int
	id         int
}

func New() Model {
	return Model{
		currDir:    ".",
		maxHeight:  0,
		idx:        0,
		keys:       DefaultKeyMap(),
		styles:     DefaultStyles(),
		min:        0,
		max:        0,
		pageDist:   37,
		halfDist:   18,
		showHidden: false,
		lastFile:   "",
		stack:      make(stack, 0),
		newIdx:     -1,
		lastIdx:    -1,
		id:         nextID(),
	}
}

func (m Model) readDir(path string) tea.Cmd {
	return func() tea.Msg {
		dirEntries, err := os.ReadDir(path)
		if err != nil {
			return err
		}
		if m.showHidden {
			return readDirMsg{id: m.id, files: dirEntries}
		}
		var filtered []os.DirEntry
		for _, f := range dirEntries {
			if f.Name()[0] != '.' {
				filtered = append(filtered, f)
			}
		}
		return readDirMsg{id: m.id, files: filtered}
	}
}

func (m Model) Init() tea.Cmd {
	return m.readDir(m.currDir)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.maxHeight = msg.Height - HeightBuffer
		m.max = m.maxHeight
	case readDirMsg:
		if msg.id != m.id {
			break
		}
		m.files = msg.files
		m.max = m.maxHeight
		if len(m.files)-1 < m.idx {
			m.idx = len(m.files) - 1
		}
		if m.lastFile != "" {
			for i, f := range m.files {
				if f.Name() == m.lastFile {
					m.idx = i
					break
				}
			}
			m.lastFile = ""
		} else if m.newIdx != -1 {
			m.idx = m.newIdx
			m.newIdx = -1
		}
		m.lastIdx = m.idx
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Quit):
			return m, tea.Quit
		case key.Matches(msg, m.keys.ForceQuit):
			return m, tea.Quit
		case key.Matches(msg, m.keys.Up):
			m.idx--
			if m.idx < 0 {
				m.idx = 0
			}
			if m.idx < m.min {
				m.min--
				m.max--
			}
		case key.Matches(msg, m.keys.Down):
			m.idx++
			if m.idx >= len(m.files) {
				m.idx = len(m.files) - 1
			}
			if m.idx > m.max {
				m.min++
				m.max++
			}
		case key.Matches(msg, m.keys.GoToTop):
			m.idx = 0
			m.min = 0
			m.max = m.maxHeight
		case key.Matches(msg, m.keys.GoToBot):
			m.idx = len(m.files) - 1
			m.min = len(m.files) - m.maxHeight
			m.max = len(m.files) - 1
		case key.Matches(msg, m.keys.HalfPgDn):
			m.idx += m.halfDist
			if m.idx >= len(m.files) {
				m.idx = len(m.files) - 1
			}
			if m.idx > m.max {
				diff := m.idx - m.max
				m.min += diff
				m.max += diff
			}
		case key.Matches(msg, m.keys.HalfPgUp):
			m.idx -= m.halfDist
			if m.idx < 0 {
				m.idx = 0
			}
			if m.idx < m.min {
				diff := m.min - m.idx
				m.min -= diff
				m.max -= diff
			}
		case key.Matches(msg, m.keys.PgDn):
			m.idx += m.pageDist
			if m.idx >= len(m.files) {
				m.idx = len(m.files) - 1
			}
			if m.idx > m.max {
				diff := m.idx - m.max
				m.min += diff
				m.max += diff
			}
		case key.Matches(msg, m.keys.PgUp):
			m.idx -= m.pageDist
			if m.idx < 0 {
				m.idx = 0
			}
			if m.idx < m.min {
				diff := m.min - m.idx
				m.min -= diff
				m.max -= diff
			}
		case key.Matches(msg, m.keys.Left):
			if m.currDir == "/" {
				break
			}
			m.stack.push(m.idx)
			m.lastFile = filepath.Base(m.currDir)
			newDir, err := filepath.Abs(m.currDir)
			if err != nil {
				log.Fatal(err)
			}
			newDir = filepath.Dir(newDir)
			m.currDir = newDir
			m.min = 0
			m.max = m.maxHeight
			return m, m.readDir(m.currDir)
		case key.Matches(msg, m.keys.Right):
			info, err := m.files[m.idx].Info()
			if err != nil {
				break
			}
			isSymlink := info.Mode()&os.ModeSymlink != 0
			if len(m.files) == 0 || (!m.files[m.idx].IsDir() && !isSymlink) {
				break
			}
			m.lastFile = ""
			if m.idx == m.lastIdx {
				m.newIdx = m.stack.pop()
			}
			m.currDir = filepath.Join(m.currDir, m.files[m.idx].Name())
			if isSymlink {
				target, err := filepath.EvalSymlinks(m.currDir)
				if err != nil {
					break
				}
				m.currDir = filepath.Dir(target)
				m.stack = make(stack, 0)
				m.newIdx = -1
			}
			m.idx = 0
			m.min = 0
			m.max = m.maxHeight
			return m, m.readDir(m.currDir)
		case key.Matches(msg, m.keys.ToggleDots):
			var hiddenCount int
			if m.showHidden {
				for _, f := range m.files {
					if f.Name()[0] != '.' {
						break
					}
					hiddenCount++
				}
			} else {
				dirEntries, err := os.ReadDir(m.currDir)
				if err != nil {
					break
				}
				for _, f := range dirEntries {
					if f.Name()[0] != '.' {
						break
					}
					hiddenCount++
				}
			}
			m.showHidden = !m.showHidden
			if m.showHidden {
				m.idx += hiddenCount
			} else {
				if m.idx < hiddenCount {
					m.idx = 0
				} else {
					m.idx -= hiddenCount
				}
			}
			return m, m.readDir(m.currDir)
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

	var files, hovered string
	for i, f := range m.files {
		if i < m.min {
			continue
		}
		if i > m.max {
			break
		}

		info, err := f.Info()
		if err != nil {
			files += fmt.Sprintf("Error reading file info: %s\n", err)
			continue
		}
		isSymlink := info.Mode()&os.ModeSymlink != 0

		file := f.Name()
		switch {
		case i == m.idx && f.IsDir():
			hovered = m.styles.PathEnd.Render(file)
			file = m.styles.DirHover.Render(file)
		case i == m.idx && isSymlink:
			hovered = m.styles.PathEnd.Render(file)
			target, err := filepath.EvalSymlinks(filepath.Join(m.currDir, file))
			if err != nil {
				file = m.styles.SymHover.Render(file + " -> ... " + fmt.Sprintf("%s", err))
				break
			}
			file = m.styles.SymHover.Render(file + " -> " + target)
		case f.IsDir():
			file = m.styles.Directory.Render(file)
		case isSymlink:
			target, err := filepath.EvalSymlinks(filepath.Join(m.currDir, file))
			if err != nil {
				file = m.styles.Symlink.Render(file + " -> ... " + fmt.Sprintf("%s", err))
				break
			}
			file = m.styles.Symlink.Render(file + " -> " + target)
		case i == m.idx:
			hovered = m.styles.PathEnd.Render(file)
			file = m.styles.Hover.Render(file)
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
