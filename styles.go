package main

import "github.com/charmbracelet/lipgloss"

type Styles struct {
	Regular         lipgloss.Style
	Directory       lipgloss.Style
	InaccessibleDir lipgloss.Style
	Symlink         lipgloss.Style
	Hover           lipgloss.Style
	Path            lipgloss.Style
	DirHover        lipgloss.Style
	SymHover        lipgloss.Style
	PathEnd         lipgloss.Style
	Filter          lipgloss.Style
	EmptyDir        lipgloss.Style
}

func DefaultStyles() Styles {
	return DefaultStylesWithRenderer(lipgloss.DefaultRenderer())
}

func DefaultStylesWithRenderer(r *lipgloss.Renderer) Styles {
	return Styles{
		Regular:         r.NewStyle(),
		Directory:       r.NewStyle().Foreground(lipgloss.Color("12")),
		InaccessibleDir: r.NewStyle().Foreground(lipgloss.Color("9")),
		Symlink:         r.NewStyle().Foreground(lipgloss.Color("10")),
		Hover:           r.NewStyle().Underline(true),
		Path:            r.NewStyle().Foreground(lipgloss.Color("12")).Bold(true),
		DirHover:        r.NewStyle().Underline(true).Foreground(lipgloss.Color("12")),
		SymHover:        r.NewStyle().Underline(true).Foreground(lipgloss.Color("10")),
		PathEnd:         r.NewStyle().Bold(true),
		Filter:          r.NewStyle().Foreground(lipgloss.Color("11")),
		EmptyDir:        r.NewStyle().Foreground(lipgloss.Color("8")).SetString("Empty"),
	}
}