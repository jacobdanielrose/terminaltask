package editmenu

import (
	"github.com/charmbracelet/lipgloss"
)

type Styles struct {
	TitleBar lipgloss.Style
	Title    lipgloss.Style

	HelpStyle lipgloss.Style

	Focused lipgloss.Style
	Blurred lipgloss.Style
	Normal  lipgloss.Style

	StatusMessage lipgloss.Style
}

func DefaultStyles() (s Styles) {
	s.TitleBar = lipgloss.NewStyle().Padding(0, 0, 1, 2) //nolint:mnd

	s.Title = lipgloss.NewStyle().
		Background(lipgloss.Color("62")).
		Foreground(lipgloss.Color("230")).
		Padding(0, 1)

	s.HelpStyle = lipgloss.NewStyle().Padding(1, 0, 0, 2) //nolint:mnd

	s.Focused = lipgloss.NewStyle().
		Border(lipgloss.NormalBorder(), false, false, false, true).
		BorderForeground(lipgloss.AdaptiveColor{Light: "#F793FF", Dark: "#AD58B4"}).
		Foreground(lipgloss.AdaptiveColor{Light: "#F793FF", Dark: "#AD58B4"}).
		Padding(0, 0, 0, 2).MarginRight(1)
	s.Normal = lipgloss.NewStyle().
		Foreground(lipgloss.AdaptiveColor{Light: "#1a1a1a", Dark: "#dddddd"}).
		Padding(0, 0, 0, 2)

	s.StatusMessage = lipgloss.NewStyle().
		Foreground(lipgloss.AdaptiveColor{Light: "#04B575", Dark: "#04B575"})

	return s
}
