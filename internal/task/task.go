package task

import (
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/lipgloss"
	"github.com/google/uuid"
)

//
// Constants
//

const (
	ellipsis = "…"
)

//
// Messages
//

type EnterEditMsg struct{}
type ToggleDoneMsg struct{}
type DeleteMsg struct{}

//
// Styles
//

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

	StatusMessage lipgloss.Style
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
		StatusMessage: lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "#04B575", Dark: "#04B575"}),
	}
}

var statusMessageStyle = lipgloss.NewStyle().
	Foreground(lipgloss.AdaptiveColor{Light: "#04B575", Dark: "#04B575"}).
	Render

//
// Keymap
//

type TaskKeyMap struct {
	EditItem   key.Binding
	ToggleDone key.Binding
	RemoveItem key.Binding
}

func newTaskKeyMap() *TaskKeyMap {
	return &TaskKeyMap{
		EditItem: key.NewBinding(
			key.WithKeys("e"),
			key.WithHelp("e", "edit item"),
		),
		ToggleDone: key.NewBinding(
			key.WithKeys(" "),
			key.WithHelp("space", "toggle done"),
		),
		RemoveItem: key.NewBinding(
			key.WithKeys("r"),
			key.WithHelp("r", "remove item"),
		),
	}
}

func (t TaskKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		t.ToggleDone,
		t.EditItem,
		t.RemoveItem,
	}
}

func (t TaskKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{
			t.ToggleDone,
			t.EditItem,
			t.RemoveItem,
		},
	}
}

//
// Task model
//

type Task struct {
	taskID   uuid.UUID
	TitleStr string
	DescStr  string
	DueDate  time.Time
	Done     bool
}

func (t Task) FilterValue() string { return t.TitleStr }
func (t Task) Title() string       { return t.TitleStr }
func (t Task) Description() string { return t.DescStr }

func (t *Task) GetID() uuid.UUID {
	return t.taskID
}

func (t *Task) SetID(id uuid.UUID) {
	t.taskID = id
}

func (t Task) IsEmpty() bool {
	return t.TitleStr == "" &&
		t.DescStr == "" &&
		t.DueDate.IsZero() &&
		!t.Done
}

func New() Task {
	return Task{
		taskID:   uuid.New(),
		TitleStr: "",
		DescStr:  "",
		// leave Duedate null here
		Done: false,
	}
}
