package editmenu

import (
	"github.com/charmbracelet/lipgloss"
)

const (
	bullet   = "•"
	ellipsis = "…"
)

type Styles struct {
	Focused lipgloss.Style
	Blurred lipgloss.Style
	Help    lipgloss.Style
	Normal  lipgloss.Style
}

func NewStyles() Styles {
	selectedStyle := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder(), false, false, false, true).
		BorderForeground(lipgloss.AdaptiveColor{Light: "#F793FF", Dark: "#AD58B4"}).
		Foreground(lipgloss.AdaptiveColor{Light: "#F793FF", Dark: "#AD58B4"}).
		Padding(0, 0, 0, 2).MarginRight(1)
	normalStyle := lipgloss.NewStyle().
		Foreground(lipgloss.AdaptiveColor{Light: "#1a1a1a", Dark: "#dddddd"}).
		Padding(0, 0, 0, 2)
	helpStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Padding(0, 0, 0, 2)

	return Styles{
		Focused: selectedStyle,
		Help:    helpStyle,
		Normal:  normalStyle,
	}
}
