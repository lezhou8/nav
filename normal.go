package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"io"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/deckarep/golang-set"
)

func (m Model) getFileLen() int {
	if m.filterState == FilterApplied {
		return len(m.filteredFiles)
	}
	return len(m.files)
}

func getSelectedFilePaths(selection map[string]mapset.Set) []string {
	var paths []string
	for dir, fileSet := range selection {
		for f := range fileSet.Iter() {
			path := filepath.Join(dir, f.(string))
			paths = append(paths, path)
		}
	}
	return paths
}

func pathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func copyFile(src, dest string) error {
	pathAlreadyExists, err := pathExists(dest)
	if err != nil {
		return err
	}
	for pathAlreadyExists {
		dest += "_"
		pathAlreadyExists, err = pathExists(dest)
		if err != nil {
			return err
		}
	}
	srcfile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcfile.Close()

	destfile, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer destfile.Close()

	_, err = io.Copy(destfile, srcfile)
	if err != nil {
		return err
	}

	srcinfo, err := os.Stat(src)
	if err != nil {
		return err
	}
	return os.Chmod(dest, srcinfo.Mode())
}

func copyDir(src, dest string) error {
	return nil
}

func copySymlink(src, dest string) error {
	return nil
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

func (m *Model) toggleSelect() {
	var f string
	if m.filterState == Unfiltered {
		f = m.files[m.idx].Name()
	} else if m.filterState == FilterApplied {
		f = m.filteredFiles[m.idx].file.Name()
	}
	m.down()
	fileSet, ok := m.selection[m.currDir]
	if !ok {
		fileSet := mapset.NewSet()
		fileSet.Add(f)
		m.selection[m.currDir] = fileSet
		return
	}
	if !fileSet.Contains(f) {
		m.selection[m.currDir].Add(f)
		return
	}
	fileSet.Remove(f)
	if fileSet.Cardinality() == 0 {
		delete(m.selection, m.currDir)
	}
}

func (m *Model) toggleSelectAll() {
	var fs []string
	if m.filterState == Unfiltered {
		for _, f := range m.files {
			fs = append(fs, f.Name())
		}
	} else if m.filterState == FilterApplied {
		for _, f := range m.filteredFiles {
			fs = append(fs, f.file.Name())
		}
	}
	fileSet, ok := m.selection[m.currDir]
	if !ok {
		fileSet := mapset.NewSet()
		for _, f := range fs {
			fileSet.Add(f)
		}
		m.selection[m.currDir] = fileSet
		return
	}
	if fileSet.Cardinality() == len(fs) {
		delete(m.selection, m.currDir)
		return
	}
	for _, f := range fs {
		fileSet.Add(f)
	}
}

func (m *Model) yank() {
	m.copyBuffer = make([]string, 0)
	if len(m.selection) == 0 {
		var currPath string
		if m.filterState == Unfiltered {
			currPath = filepath.Join(m.currDir, m.files[m.idx].Name())
		} else if m.filterState == FilterApplied {
			currPath = filepath.Join(m.currDir, m.filteredFiles[m.idx].file.Name())
		} else {
			m.news = "Yank error"
			return
		}
		m.copyBuffer = append(m.copyBuffer, currPath)
		m.news = "Yanked 1 file"
		return
	}
	m.copyBuffer = getSelectedFilePaths(m.selection)
	var plural string
	copyBufferLen := len(m.copyBuffer)
	if copyBufferLen == 1 {
		plural = " file"
	} else {
		plural = " files"
	}
	m.news = "Yanked " + fmt.Sprintf("%d", len(m.copyBuffer)) + plural
}

func (m *Model) cut() {
	m.yank()
	cutAmount := len(m.copyBuffer)
	if cutAmount == 1 {
		m.news = "1 file ready to be cut and pasted"
	} else if 1 < cutAmount {
		m.news = fmt.Sprintf("%d", cutAmount) + " files ready to be cut and pasted"
	} else {
		m.news = "Cut error"
	}
	m.isCutting = true
}

func (m *Model) paste() {
	if len(m.copyBuffer) == 0 {
		m.news = "Nothing pasted"
		return
	}
	dirEmpty := len(m.files) == 0
	dest := m.currDir
	anyErrors := false
	for _, f := range m.copyBuffer {
		fInfo, err := os.Stat(f)
		if err != nil {
			anyErrors = true
			continue
		}
		if fInfo.Mode()&os.ModeSymlink != 0 {
			if err := copySymlink(f, filepath.Join(dest, filepath.Base(f))); err != nil {
				anyErrors = true
			}
		} else if fInfo.IsDir() {
			if err := copyDir(f, dest); err != nil {
				anyErrors = true
			}
		} else {
			if err := copyFile(f, filepath.Join(dest, filepath.Base(f))); err != nil {
				anyErrors = true
			}
		}
	}
	if dirEmpty {
		m.idx = 0
	}
	if !m.isCutting {
		return
	}
	if anyErrors {
		m.news = "Error while cutting and pasting"
		return
	}
	for _, f := range m.copyBuffer {
		fInfo, err := os.Stat(f)
		if err != nil {
			continue
		}
		if fInfo.IsDir() {
			err := os.RemoveAll(f)
			if err != nil {
				continue
			}
		} else {
			err := os.Remove(f)
			if err != nil {
				continue
			}
		}
	}
	m.isCutting = false
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
		oldNews := m.news
		m.news = ""
		switch {
		case key.Matches(msg, m.keys.Quit):
			m.news = oldNews
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
		case key.Matches(msg, m.keys.ToggleSelect):
			m.toggleSelect()
		case key.Matches(msg, m.keys.ToggleSelectAll):
			m.toggleSelectAll()
		case key.Matches(msg, m.keys.Yank):
			m.yank()
		case key.Matches(msg, m.keys.Cut):
			m.cut()
		case key.Matches(msg, m.keys.Paste):
			m.paste()
			if m.filterState == FilterApplied {
				m.filterOff()
			}
			return m, m.readDir(m.currDir)
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
