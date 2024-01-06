package main

import (
	"os"
	"sort"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/sahilm/fuzzy"
)

type FilterState int

const (
	Unfiltered FilterState = iota
	Filtering
	FilterApplied
)

type filteredFile struct {
	file    os.DirEntry
	matches []int
}

type filteredFiles []filteredFile

type FilterMatchesMsg filteredFiles

type FilterFunc func(string, []string) []Rank

type Rank struct {
	Index          int
	MatchedIndexes []int
}

func (ff filteredFiles) filteredFilesAsDirEntries() []os.DirEntry {
	de := make([]os.DirEntry, len(ff))
	for i, f := range ff {
		de[i] = f.file
	}
	return de
}

func DefaultFilter(term string, targets []string) []Rank {
	ranks := fuzzy.Find(term, targets)
	sort.Stable(ranks)
	result := make([]Rank, len(ranks))
	for i, r := range ranks {
		result[i] = Rank{
			Index:          r.Index,
			MatchedIndexes: r.MatchedIndexes,
		}
	}
	return result
}

func filterFiles(m Model) tea.Cmd {
	return func() tea.Msg {
		if m.filterInput.Value() == "" {
			return FilterMatchesMsg(m.filesAsFilterFiles())
		}
		fs := m.files
		targets := make([]string, len(fs))

		for i, f := range fs {
			targets[i] = f.Name()
		}

		filterMatches := filteredFiles{}
		for _, r := range m.filter(m.filterInput.Value(), targets) {
			filterMatches = append(filterMatches, filteredFile{
				file:    fs[r.Index],
				matches: r.MatchedIndexes,
			})
		}

		return FilterMatchesMsg(filterMatches)
	}
}

func (m Model) filesAsFilterFiles() filteredFiles {
	ff := make(filteredFiles, len(m.files))
	for i, f := range m.files {
		ff[i] = filteredFile{
			file: f,
		}
	}
	return ff
}

func (m *Model) filterAccept() {
	m.filterState = FilterApplied
	m.idx = 0
}

func (m *Model) filterOn() {
	m.filterState = Filtering
	if m.filterInput.Value() == "" {
		m.filteredFiles = m.filesAsFilterFiles()
	}
	m.filterInput.Focus()
}

func (m *Model) filterOff() {
	m.filterState = Unfiltered
	m.filterInput.Reset()
	m.filteredFiles = nil
	m.news = ""
}

func (m Model) filterMode(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.FilterOff):
			m.filterOff()
		case key.Matches(msg, m.keys.FilterAccept):
			m.filterAccept()
		}
	}
	newFilterInputModel, inputCmd := m.filterInput.Update(msg)
	filterChanged := m.filterInput.Value() != newFilterInputModel.Value()
	m.filterInput = newFilterInputModel
	cmds = append(cmds, inputCmd)
	if filterChanged {
		cmds = append(cmds, filterFiles(m))
	}
	return m, tea.Batch(cmds...)
}
