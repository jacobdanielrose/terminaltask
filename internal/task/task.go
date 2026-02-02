// Package task defines the Task type and related styles, key bindings,
// and messages used to represent and render tasks in the list view.
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

// EnterEditMsg signals that the user wants to edit the currently
// selected task in the list.
type EnterEditMsg struct{}

// ToggleDoneMsg signals that the user toggled the completion state
// of the currently selected task.
type ToggleDoneMsg struct{}

// DeleteMsg signals that the user requested deletion of the
// currently selected task.
type DeleteMsg struct{}

//
// Styles
//

// subStyle groups styles for the individual parts of a task: title,
// description, and due date.
type subStyle struct {
	Title lipgloss.Style
	Desc  lipgloss.Style
	Date  lipgloss.Style
}

// Styles contains all styles used to render tasks in the list,
// including normal, selected, dimmed, and filter match highlighting.
type Styles struct {
	// The Normal state.
	Normal subStyle

	// The selected item state.
	Selected subStyle

	// The dimmed state, for when the filter input is initially activated.
	Dimmed subStyle

	// Characters matching the current filter, if any.
	FilterMatch lipgloss.Style

	// StatusMessage styles status text shown for task operations.
	StatusMessage lipgloss.Style
}

// newTaskStyles constructs the default Styles used for rendering
// tasks in the list.
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

//
// Keymap
//

// TaskKeyMap defines key bindings for actions on a single task in
// the list, such as editing, toggling completion, and removing.
type TaskKeyMap struct {
	EditItem   key.Binding
	ToggleDone key.Binding
	RemoveItem key.Binding
}

// newTaskKeyMap constructs the default key bindings used for tasks
// in the list.
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

// ShortHelp implements the help.KeyMap interface for condensed help.
func (t TaskKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		t.ToggleDone,
		t.EditItem,
		t.RemoveItem,
	}
}

// FullHelp implements the help.KeyMap interface for the full help
// view for task-related key bindings.
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

// Task represents a single task, including ID, title, description,
// due date, and completion status.
type Task struct {
	taskID   uuid.UUID
	TitleStr string
	DescStr  string
	DueDate  time.Time
	Done     bool
}

// FilterValue implements list.Item and is used by the list filter.
func (t Task) FilterValue() string { return t.TitleStr }

// Title returns the task title.
func (t Task) Title() string { return t.TitleStr }

// Description returns the task description.
func (t Task) Description() string { return t.DescStr }

// GetID returns the task's unique identifier.
func (t Task) GetID() uuid.UUID {
	return t.taskID
}

// SetID sets the task's unique identifier.
func (t *Task) SetID(id uuid.UUID) {
	t.taskID = id
}

// IsEmpty reports whether the task has no title, description, due
// date, and is not marked as done.
func (t Task) IsEmpty() bool {
	return t.TitleStr == "" &&
		t.DescStr == "" &&
		t.DueDate.IsZero() &&
		!t.Done
}

// New constructs a new, empty Task with a generated ID.
func New() Task {
	return Task{
		taskID:   uuid.New(),
		TitleStr: "",
		DescStr:  "",
		// leave Duedate null here
		Done: false,
	}
}

// NewWithOptions constructs a Task with the provided values and a
// generated ID.
func NewWithOptions(title, desc string, duedate time.Time, done bool) Task {
	return Task{
		taskID:   uuid.New(),
		TitleStr: title,
		DescStr:  desc,
		DueDate:  duedate,
		Done:     done,
	}
}
