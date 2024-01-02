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
	"github.com/deckarep/golang-set"
)

const (
	HeightBuffer int    = 7
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
	filteredFiles filteredFiles
	filterInput   textinput.Model
	selection     map[string]mapset.Set
	copyBuffer    []string
	isCutting     bool
	news          string
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
		selection:   make(map[string]mapset.Set),
		copyBuffer:  make([]string, 0),
		isCutting:   false,
		news:        "",
		id:          nextID(),
	}
}

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
	case FilterMatchesMsg:
		m.filteredFiles = filteredFiles(msg)
		return m, nil
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

	files := ""
	hovered := ""
	if m.filterState == Unfiltered || m.filterState == FilterApplied {
		var filesIterate []os.DirEntry
		if m.filterState == Unfiltered {
			filesIterate = m.files
		} else {
			filesIterate = m.filteredFiles.filteredFilesAsDirEntries()
		}
		for i, f := range filesIterate {
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
					file = m.styles.SymHover.Render(file + " -> " + fmt.Sprintf("%s", err))
					break
				}
				file = m.styles.SymHover.Render(file + " -> " + target)
			case f.IsDir():
				file = m.styles.Directory.Render(file + "/")
			case isSymlink:
				target, err := filepath.EvalSymlinks(filepath.Join(m.currDir, file))
				if err != nil {
					file = m.styles.Symlink.Render(file + " -> " + fmt.Sprintf("%s", err))
					break
				}
				file = m.styles.Symlink.Render(file + " -> " + target)
			case i == m.idx:
				hovered = m.styles.PathEnd.Render(file)
				file = m.styles.Hover.Render(file)
			}
			fileSet, ok := m.selection[m.currDir]
			if ok && fileSet.Contains(f.Name()) {
				file = m.styles.Selected.Render(file)
			}
			files += file + "\n"
		}
	} else {
		for _, f := range m.filteredFiles {
			files += m.styles.Filter.Render(f.file.Name()) + "\n"
		}
	}
	filterBar := "\n\n"
	if m.filterState == Filtering || m.filterState == FilterApplied {
		filterBar = "\n" + m.styles.Filter.Render(m.filterInput.View()) + "\n"
	}
	return currPath + hovered + filterBar + files + m.news + "\n"
}

func main() {
	if _, err := tea.NewProgram(New()).Run(); err != nil {
		log.Fatal(err)
	}
}
