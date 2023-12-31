package main

import "github.com/charmbracelet/bubbles/key"

type KeyMap struct {
	Up           key.Binding
	Down         key.Binding
	Right        key.Binding
	Left         key.Binding
	GoToTop      key.Binding
	GoToBot      key.Binding
	PgDn         key.Binding
	PgUp         key.Binding
	HalfPgDn     key.Binding
	HalfPgUp     key.Binding
	ToggleDots   key.Binding
	GoHome       key.Binding
	FilterOn     key.Binding
	FilterOff    key.Binding
	FilterAccept key.Binding
	Quit         key.Binding
	ForceQuit    key.Binding
	// TODO: select, copy, cut, paste, bulk rename
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
		PgDn: key.NewBinding(
			key.WithKeys("pgdown", "ctrl+f"),
			key.WithHelp("pagedown/ctrl+f", "page down"),
		),
		PgUp: key.NewBinding(
			key.WithKeys("pgup", "ctrl+b"),
			key.WithHelp("pageup/ctrl+b", "page up"),
		),
		HalfPgDn: key.NewBinding(
			key.WithKeys("ctrl+d"),
			key.WithHelp("ctrl+d", "half page down"),
		),
		HalfPgUp: key.NewBinding(
			key.WithKeys("ctrl+u"),
			key.WithHelp("ctrl+u", "half page up"),
		),
		ToggleDots: key.NewBinding(
			key.WithKeys("."),
			key.WithHelp(".", "toggle show hidden files"),
		),
		GoHome: key.NewBinding(
			key.WithKeys("~"),
			key.WithHelp("~", "go to home directory"),
		),
		FilterOn: key.NewBinding(
			key.WithKeys("/"),
			key.WithHelp("/", "search mode on"),
		),
		FilterOff: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "search mode off"),
		),
		FilterAccept: key.NewBinding(
			key.WithKeys("enter", "tab", "shift+tab", "ctrl+k", "up", "ctrl+j", "down"),
			key.WithHelp("enter", "apply filter"),
		),
		Quit: key.NewBinding(
			key.WithKeys("q"),
			key.WithHelp("q", "quit"),
		),
		ForceQuit: key.NewBinding(key.WithKeys("ctrl+c")),
	}
}
