package main

import "github.com/charmbracelet/lipgloss"

type Styles struct {
	Directory lipgloss.Style
	Symlink   lipgloss.Style
	Hover     lipgloss.Style
	Path      lipgloss.Style
	PathEnd   lipgloss.Style
	EmptyDir  lipgloss.Style
}

func DefaultStyles() Styles {
	return DefaultStylesWithRenderer(lipgloss.DefaultRenderer())
}

func DefaultStylesWithRenderer(r *lipgloss.Renderer) Styles {
	return Styles{
		Directory: r.NewStyle().Foreground(lipgloss.Color("12")),
		Symlink:   r.NewStyle().Foreground(lipgloss.Color("10")),
		Hover:     r.NewStyle().Bold(true),
		Path:      r.NewStyle().Foreground(lipgloss.Color("12")).Bold(true),
		PathEnd:   r.NewStyle().Bold(true),
		EmptyDir:  r.NewStyle().Foreground(lipgloss.Color("8")).SetString("Empty directory"),
	}
}
