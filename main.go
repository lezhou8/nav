package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

const (
	HeightBuffer int    = 5
	XDGCacheDir  string = "$XDG_CACHE_HOME"
	CacheSubDir  string = "nav"
	CacheFile    string = ".nav_d"
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

func isDirAccessible(path string) bool {
	f, err := os.Open(path)
	if err != nil {
		return false
	}
	defer f.Close()
	return true
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

type Model struct {
	files         []os.DirEntry
	currDir       string
	maxHeight     int
	idx           int
	keys          KeyMap
	styles        Styles
	min           int
	max           int
	pageDist      int
	halfDist      int
	showHidden    bool
	lastFile      string
	cursorSave    map[string]int
	filter        FilterFunc
	filterState   FilterState
	filteredItems filteredItems
	filterInput   textinput.Model
	id            int
}

func New() Model {
	dir, err := filepath.Abs(".")
	if err != nil {
		log.Fatal(err)
	}
	filterInput := textinput.New()
	filterInput.Prompt = "/"
	return Model{
		currDir:     dir,
		maxHeight:   0,
		idx:         0,
		keys:        DefaultKeyMap(),
		styles:      DefaultStyles(),
		min:         0,
		max:         0,
		pageDist:    37,
		halfDist:    18,
		showHidden:  false,
		lastFile:    "",
		cursorSave:  make(map[string]int),
		filter:      DefaultFilter,
		filterState: Unfiltered,
		filterInput: filterInput,
		id:          nextID(),
	}
}

// func (m *Model) updateKeyBindings() {
// 	switch m.filterState {
// 	case Filtering:
// 		m.keys.Up.SetEnabled(false)
// 		m.keys.Down.SetEnabled(false)
// 		m.keys.GoToTop.SetEnabled(false)
// 		m.keys.GoToBot.SetEnabled(false)
// 		m.keys.HalfPgDn.SetEnabled(false)
// 		m.keys.HalfPgUp.SetEnabled(false)
// 		m.keys.PgDn.SetEnabled(false)
// 		m.keys.PgUp.SetEnabled(false)
// 		m.keys.FilterOn.SetEnabled(false)
// 		m.keys.Left.SetEnabled(false)
// 		m.keys.Right.SetEnabled(false)
// 		m.keys.ToggleDots.SetEnabled(false)
// 		m.keys.GoHome.SetEnabled(false)
// 	default:
// 		m.keys.Up.SetEnabled(true)
// 		m.keys.Down.SetEnabled(true)
// 		m.keys.GoToTop.SetEnabled(true)
// 		m.keys.GoToBot.SetEnabled(true)
// 		m.keys.HalfPgDn.SetEnabled(true)
// 		m.keys.HalfPgUp.SetEnabled(true)
// 		m.keys.PgDn.SetEnabled(true)
// 		m.keys.PgUp.SetEnabled(true)
// 		m.keys.FilterOn.SetEnabled(true)
// 		m.keys.Left.SetEnabled(true)
// 		m.keys.Right.SetEnabled(true)
// 		m.keys.ToggleDots.SetEnabled(true)
// 		m.keys.GoHome.SetEnabled(true)
// 	}
// }

func (m *Model) refreshFiles() {
	if len(m.files)-1 < m.idx {
		m.idx = len(m.files) - 1
	}
	if m.lastFile == "" {
		return
	}
	for i, f := range m.files {
		if f.Name() == m.lastFile {
			m.idx = i
			break
		}
	}
	m.lastFile = ""
}

func (m Model) Init() tea.Cmd {
	return m.readDir(m.currDir)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if key.Matches(msg, m.keys.ForceQuit) {
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		m.maxHeight = msg.Height - HeightBuffer
		m.max = m.maxHeight
	case readDirMsg:
		if msg.id != m.id {
			break
		}
		m.files = msg.files
		m.max = m.maxHeight
		m.refreshFiles()
	}

	if m.filterState == Filtering {
		return m.filterMode(msg)
	}
	return m.normalMode(msg)
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
			file = m.styles.DirHover.Render(file + "/")
		case i == m.idx && isSymlink:
			hovered = m.styles.PathEnd.Render(file)
			target, err := filepath.EvalSymlinks(filepath.Join(m.currDir, file))
			if err != nil {
				file = m.styles.SymHover.Render(file + " -> ... " + fmt.Sprintf("%s", err))
				break
			}
			file = m.styles.SymHover.Render(file + " -> " + target)
		case f.IsDir():
			file = m.styles.Directory.Render(file + "/")
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
	filterBar := "\n\n"
	if m.filterState == Filtering || m.filterState == FilterApplied {
		filterBar = "\n" + m.styles.Filter.Render(m.filterInput.View()) + "\n"
	}
	return currPath + hovered + filterBar + files
}

func main() {
	if _, err := tea.NewProgram(New()).Run(); err != nil {
		log.Fatal(err)
	}
}
