package task

import "github.com/charmbracelet/lipgloss"

const (
	bullet   = "•"
	ellipsis = "…"
)

type subStyle struct {
	Title lipgloss.Style
	Desc  lipgloss.Style
	Date  lipgloss.Style
}

type Styles struct {
	// The Normal state.
	Normal subStyle

	// The selected item state.
	Selected subStyle

	// The dimmed state, for when the filter input is initially activated.
	Dimmed subStyle

	// Characters matching the current filter, if any.
	FilterMatch lipgloss.Style
}

func newTaskStyles() Styles {
	return Styles{
		Normal: subStyle{
			Title: lipgloss.NewStyle().
				Foreground(lipgloss.AdaptiveColor{Light: "#1a1a1a", Dark: "#dddddd"}).
				Padding(0, 0, 0, 2), //nolint:mnd
			Desc: lipgloss.NewStyle().
				Foreground(lipgloss.AdaptiveColor{Light: "#A49FA5", Dark: "#777777"}).
				Padding(0, 0, 0, 2),
			Date: lipgloss.NewStyle().
				Foreground(lipgloss.AdaptiveColor{Light: "#A49FA5", Dark: "#777777"}).
				Padding(0, 0, 0, 2),
		},
		Selected: subStyle{
			Title: lipgloss.NewStyle().
				Border(lipgloss.NormalBorder(), false, false, false, true).
				BorderForeground(lipgloss.AdaptiveColor{Light: "#F793FF", Dark: "#AD58B4"}).
				Foreground(lipgloss.AdaptiveColor{Light: "#F793FF", Dark: "#AD58B4"}).
				Padding(0, 0, 0, 2), //nolint:mnd
			Desc: lipgloss.NewStyle().
				Border(lipgloss.NormalBorder(), false, false, false, true).
				BorderForeground(lipgloss.AdaptiveColor{Light: "#F793FF", Dark: "#AD58B4"}).
				Foreground(lipgloss.AdaptiveColor{Light: "#F793FF", Dark: "#AD58B4"}).
				Padding(0, 0, 0, 2),
			Date: lipgloss.NewStyle().
				Border(lipgloss.NormalBorder(), false, false, false, true).
				BorderForeground(lipgloss.AdaptiveColor{Light: "#F793FF", Dark: "#AD58B4"}).
				Foreground(lipgloss.AdaptiveColor{Light: "#F793FF", Dark: "#AD58B4"}).
				Padding(0, 0, 0, 2),
		},
		Dimmed: subStyle{
			Title: lipgloss.NewStyle().
				Foreground(lipgloss.AdaptiveColor{Light: "#A49FA5", Dark: "#777777"}).
				Padding(0, 0, 0, 2), //nolint:mnd
			Desc: lipgloss.NewStyle().
				Foreground(lipgloss.AdaptiveColor{Light: "#A49FA5", Dark: "#777777"}),
			Date: lipgloss.NewStyle().
				Foreground(lipgloss.AdaptiveColor{Light: "#A49FA5", Dark: "#777777"}),
		},
		FilterMatch: lipgloss.NewStyle().Underline(true),
	}
}
