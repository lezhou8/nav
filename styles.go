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
	Selected        lipgloss.Style
	News			lipgloss.Style
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
		Hover:           r.NewStyle().Background(lipgloss.Color("15")).Foreground(lipgloss.Color("0")),
		Path:            r.NewStyle().Foreground(lipgloss.Color("12")).Bold(true),
		DirHover:        r.NewStyle().Background(lipgloss.Color("12")).Foreground(lipgloss.Color("0")),
		SymHover:        r.NewStyle().Background(lipgloss.Color("10")).Foreground(lipgloss.Color("0")),
		PathEnd:         r.NewStyle().Bold(true),
		Filter:          r.NewStyle().Foreground(lipgloss.Color("11")),
		Selected:        r.NewStyle().Italic(true).Bold(true),
		News:			 r.NewStyle().Italic(true),
		EmptyDir:        r.NewStyle().Foreground(lipgloss.Color("8")).SetString("Empty"),
	}
}
