package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) getFileLen() int {
	if m.filterState == FilterApplied {
		return len(m.filteredFiles)
	}
	return len(m.files)
}

func (m Model) quitRoutine() {
	cacheDir := os.Getenv(XDGCacheDir)
	if cacheDir == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			log.Fatal(err)
		}
		cacheDir = filepath.Join(homeDir, ".cache")
	}
	subDir := filepath.Join(cacheDir, CacheSubDir)
	err := os.MkdirAll(subDir, 0755)
	if err != nil {
		log.Fatal(err)
	}
	fp := filepath.Join(subDir, CacheFile)
	f, err := os.Create(fp)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	data := []byte(m.currDir + "\n")
	_, err = f.Write(data)
	if err != nil {
		log.Fatal(err)
	}
}

func (m *Model) up() {
	m.idx--
	if m.idx < 0 {
		m.idx = 0
	}
	if m.idx < m.min {
		m.min--
		m.max--
	}
}

func (m *Model) down() {
	fileLen := m.getFileLen()
	m.idx++
	if m.idx >= fileLen {
		m.idx = fileLen - 1
	}
	if m.idx > m.max {
		m.min++
		m.max++
	}
}

func (m *Model) goToTop() {
	m.idx = 0
	m.min = 0
	m.max = m.maxHeight
}

func (m *Model) goToBot() {
	fileLen := m.getFileLen()
	m.idx = fileLen - 1
	m.min = fileLen - m.maxHeight
	m.max = fileLen - 1
}

func (m *Model) halfPgDn() {
	fileLen := m.getFileLen()
	m.idx += m.halfDist
	if m.idx >= fileLen {
		m.idx = fileLen - 1
	}
	if m.idx > m.max {
		diff := m.idx - m.max
		m.min += diff
		m.max += diff
	}
}

func (m *Model) halfPgUp() {
	m.idx -= m.halfDist
	if m.idx < 0 {
		m.idx = 0
	}
	if m.idx < m.min {
		diff := m.min - m.idx
		m.min -= diff
		m.max -= diff
	}
}

func (m *Model) pgDn() {
	fileLen := m.getFileLen()
	m.idx += m.pageDist
	if m.idx >= fileLen {
		m.idx = fileLen - 1
	}
	if m.idx > m.max {
		diff := m.idx - m.max
		m.min += diff
		m.max += diff
	}
}

func (m *Model) pgUp() {
	m.idx -= m.pageDist
	if m.idx < 0 {
		m.idx = 0
	}
	if m.idx < m.min {
		diff := m.min - m.idx
		m.min -= diff
		m.max -= diff
	}
}

func (m Model) left() (tea.Model, tea.Cmd) {
	if m.currDir == "/" {
		return m, nil
	}
	m.lastFile = filepath.Base(m.currDir)
	m.cursorSave[m.currDir] = m.idx
	newDir, err := filepath.Abs(m.currDir)
	if err != nil {
		log.Fatal(err)
	}
	newDir = filepath.Dir(newDir)
	m.currDir = newDir
	m.min = 0
	m.max = m.maxHeight
	m.filterOff()
	return m, m.readDir(m.currDir)
}

func (m Model) right() (tea.Model, tea.Cmd) {
	info, err := m.files[m.idx].Info()
	if err != nil {
		return m, nil
	}
	isSymlink := info.Mode()&os.ModeSymlink != 0
	if len(m.files) == 0 || (!m.files[m.idx].IsDir() && !isSymlink) {
		return m, nil
	}
	oldDir := m.currDir
	newPath := filepath.Join(m.currDir, m.files[m.idx].Name())
	if !isDirAccessible(newPath) {
		return m, nil
	}
	if isSymlink {
		target, err := filepath.EvalSymlinks(newPath)
		if err != nil {
			return m, nil
		}
		targetInfo, err := os.Stat(target)
		if err != nil {
			return m, nil
		}
		if !targetInfo.IsDir() {
			return m, nil
		}
		m.currDir = target
	} else {
		m.currDir = newPath
	}
	m.cursorSave[oldDir] = m.idx
	m.lastFile = ""
	if val, ok := m.cursorSave[m.currDir]; ok {
		m.idx = val
	} else {
		m.idx = 0
	}
	m.min = 0
	m.max = m.maxHeight
	m.filterOff()
	return m, m.readDir(m.currDir)
}

func (m Model) toggleDots() (tea.Model, tea.Cmd) {
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
			return m, nil
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

func (m Model) goHome() (tea.Model, tea.Cmd) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return m, nil
	}
	m.cursorSave[m.currDir] = m.idx
	m.currDir = homeDir
	m.lastFile = ""
	if val, ok := m.cursorSave[m.currDir]; ok {
		m.idx = val
	} else {
		m.idx = 0
	}
	m.min = 0
	m.max = m.maxHeight
	m.filterOff()
	return m, m.readDir(m.currDir)
}

func (m Model) normalMode(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Quit):
			m.quitRoutine()
			return m, tea.Quit
		case key.Matches(msg, m.keys.Up):
			m.up()
		case key.Matches(msg, m.keys.Down):
			m.down()
		case key.Matches(msg, m.keys.GoToTop):
			m.goToTop()
		case key.Matches(msg, m.keys.GoToBot):
			m.goToBot()
		case key.Matches(msg, m.keys.HalfPgDn):
			m.halfPgDn()
		case key.Matches(msg, m.keys.HalfPgUp):
			m.halfPgUp()
		case key.Matches(msg, m.keys.PgDn):
			m.pgDn()
		case key.Matches(msg, m.keys.PgUp):
			m.pgUp()
		case key.Matches(msg, m.keys.FilterOn) && 0 < len(m.files):
			m.filterOn()
		case key.Matches(msg, m.keys.FilterOff):
			m.filterOff()
		case key.Matches(msg, m.keys.Left):
			return m.left()
		case key.Matches(msg, m.keys.Right):
			return m.right()
		case key.Matches(msg, m.keys.ToggleDots):
			return m.toggleDots()
		case key.Matches(msg, m.keys.GoHome):
			return m.goHome()
		}
	}
	return m, nil
}
