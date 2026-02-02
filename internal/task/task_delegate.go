package task

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	"github.com/charmbracelet/x/ansi"
)

const (
	defaultDelegateHeight  = 3
	defaultDelegateSpacing = 1
	dateFormat             = "2006-01-02"
)

//
// Delegate
//

// TaskDelegate implements list.ItemDelegate for Task items, handling
// rendering and key bindings for actions on tasks in the list.
type TaskDelegate struct {
	Styles  Styles
	height  int
	spacing int
	keymap  *TaskKeyMap
}

// NewTaskDelegate constructs a TaskDelegate with default styles and keymap.
func NewTaskDelegate() TaskDelegate {
	return TaskDelegate{
		Styles:  newTaskStyles(),
		keymap:  newTaskKeyMap(),
		height:  defaultDelegateHeight,
		spacing: defaultDelegateSpacing,
	}
}

// Height returns the number of lines used to render each task.
func (t TaskDelegate) Height() int { return t.height }

// Spacing returns the number of blank lines between rendered tasks.
func (t TaskDelegate) Spacing() int { return t.spacing }

// Update handles key events for the list delegate
func (t TaskDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, t.keymap.ToggleDone):
			return func() tea.Msg {
				return ToggleDoneMsg{}
			}

		case key.Matches(msg, t.keymap.EditItem):
			return func() tea.Msg {
				return EnterEditMsg{}
			}

		case key.Matches(msg, t.keymap.RemoveItem):
			if len(m.Items()) == 0 {
				t.keymap.RemoveItem.SetEnabled(false)
			}
			return func() tea.Msg {
				return DeleteMsg{}
			}
		}
	}

	return nil
}

// Render draws a Task into the list, applying styles for selection,
// filtering, and completion state.
func (t TaskDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	var (
		title, desc, date string
		titleStyle        lipgloss.Style
		descStyle         lipgloss.Style
		dateStyle         lipgloss.Style
		done              bool
		matchedRunes      []int
		s                 = &t.Styles
	)

	if i, ok := item.(Task); !ok {
		log.Error("Could not render task",
			"task.title", i.Title(),
			"task.desc", i.Description(),
			"task.done", i.Done,
		)
		return
	} else {
		title = i.Title()
		desc = i.Description()
		done = i.Done
		date = i.DueDate.Format(dateFormat)
	}

	if m.Width() <= 0 {
		// short-circuit
		log.Error("Model width is less than zero!")
		return
	}

	// Prevent text from exceeding list width
	textwidth := m.Width() - s.Normal.Title.GetPaddingLeft() - s.Normal.Title.GetPaddingRight()
	title = ansi.Truncate(title, textwidth, ellipsis)

	var lines []string
	for i, line := range strings.Split(desc, "\n") {
		if i >= t.height-1 {
			break
		}
		lines = append(lines, ansi.Truncate(line, textwidth, ellipsis))
	}
	desc = strings.Join(lines, "\n")

	// Conditions
	var (
		isSelected  = index == m.Index()
		emptyFilter = m.FilterState() == list.Filtering && m.FilterValue() == ""
		isFiltered  = m.FilterState() == list.Filtering || m.FilterState() == list.FilterApplied
	)

	if isFiltered && index < len(m.VisibleItems()) {
		// Get indices of matched characters
		matchedRunes = m.MatchesForItem(index)
	}

	if emptyFilter {
		titleStyle = s.Dimmed.Title
		descStyle = s.Dimmed.Desc
		dateStyle = s.Dimmed.Date
	} else if isSelected && m.FilterState() != list.Filtering {
		if isFiltered {
			// Highlight matches
			unmatched := s.Selected.Title.Inline(true)
			matched := unmatched.Inherit(s.FilterMatch)
			title = lipgloss.StyleRunes(title, matchedRunes, matched, unmatched)
		}
		titleStyle = s.Selected.Title
		descStyle = s.Selected.Desc
		dateStyle = s.Selected.Date
	} else {
		if isFiltered {
			// Highlight matches
			unmatched := s.Normal.Title.Inline(true)
			matched := unmatched.Inherit(s.FilterMatch)
			title = lipgloss.StyleRunes(title, matchedRunes, matched, unmatched)
		}
		titleStyle = s.Normal.Title
		descStyle = s.Normal.Desc
		dateStyle = s.Normal.Date
	}

	// Mark Done
	if done {
		titleStyle = titleStyle.Inherit(titleStyle).Strikethrough(true)
		descStyle = descStyle.Inherit(descStyle).Strikethrough(true)
		dateStyle = dateStyle.Inherit(dateStyle).Strikethrough(true)
	} else {
		titleStyle = titleStyle.Inherit(titleStyle).Strikethrough(false)
		descStyle = descStyle.Inherit(descStyle).Strikethrough(false)
		dateStyle = dateStyle.Inherit(dateStyle).Strikethrough(false)
	}

	n, err := fmt.Fprintf(
		w,
		"%s\n%s\n%s",
		titleStyle.Render(title),
		descStyle.Render(desc),
		dateStyle.Render(date),
	)
	if err != nil {
		log.Error("Failed to render delegate", "err", err, "bytes", n)
	}
}

// ShortHelp implements the help.KeyMap interface for condensed help
// for task-related key bindings.
func (t TaskDelegate) ShortHelp() []key.Binding {
	return []key.Binding{
		t.keymap.ToggleDone,
		t.keymap.EditItem,
		t.keymap.RemoveItem,
	}
}

// FullHelp implements the help.KeyMap interface for the full help view
// of task-related key bindings.
func (t TaskDelegate) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{
			t.keymap.ToggleDone,
			t.keymap.EditItem,
			t.keymap.RemoveItem,
		},
	}
}
