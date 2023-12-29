package main

import "github.com/charmbracelet/bubbles/key"

type KeyMap struct {
	Up        key.Binding
	Down      key.Binding
	Right     key.Binding
	Left      key.Binding
	GoToTop   key.Binding
	GoToBot   key.Binding
	Quit      key.Binding
	ForceQuit key.Binding
	// TODO: select, copy, cut, paste, search, bulk rename, PgDn, PgUp
}

func DefaultKeyMap() KeyMap {
	return KeyMap{
		Up: key.NewBinding(
			key.WithKeys("k", "up"),
			key.WithHelp("^/k", "up"),
		),
		Down: key.NewBinding(
			key.WithKeys("down", "j"),
			key.WithHelp("v/j", "down"),
		),
		Right: key.NewBinding(
			key.WithKeys("right", "l"),
			key.WithHelp(">/l", "open directory"),
		),
		Left: key.NewBinding(
			key.WithKeys("left", "h"),
			key.WithHelp("</h", "go to parent directory"),
		),
		GoToTop: key.NewBinding(
			key.WithKeys("home", "g"),
			key.WithHelp("g/home", "go to top"),
		),
		GoToBot: key.NewBinding(
			key.WithKeys("end", "G"),
			key.WithHelp("G/end", "go to bottom"),
		),
		Quit: key.NewBinding(
			key.WithKeys("q"),
			key.WithHelp("q", "quit"),
		),
		ForceQuit: key.NewBinding(key.WithKeys("ctrl+c")),
	}
}
