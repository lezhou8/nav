package main

import (
	"github.com/sahilm/fuzzy"
	"os"
	"sort"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/bubbles/key"
)

type FilterState int

const (
	Unfiltered FilterState = iota
	Filtering
	FilterApplied
)

type filteredItem struct {
	item    os.DirEntry
	matches []int
}

type filteredItems []filteredItem

func (f filteredItems) items() []os.DirEntry {
	agg := make([]os.DirEntry, len(f))
	for i, v := range f {
		agg[i] = v.item
	}
	return agg
}

type FilterMatchesMsg []filteredItem

type FilterFunc func(string, []string) []Rank

type Rank struct {
	Index          int
	MatchedIndexes []int
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

func (m Model) filterMode(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.FilterOff):
			m.filterOff()
		case key.Matches(msg, m.keys.FilterAccept):
			m.filterAccept()
		}
	}
	m.filterInput, cmd = m.filterInput.Update(msg)
	return m, cmd
}

func (m *Model) filterAccept() {
	m.filterState = FilterApplied
	// m.updateKeyBindings()
}

func (m *Model) filterOn() {
	m.filterState = Filtering
	m.filterInput.Focus()
	// m.updateKeyBindings()
}

func (m *Model) filterOff() {
	m.filterState = Unfiltered
	// m.updateKeyBindings()
}
