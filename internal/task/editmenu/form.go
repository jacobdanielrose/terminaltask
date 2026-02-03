package editmenu

import (
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	datepicker "github.com/ethanefung/bubble-datepicker"
)

const (
	focusIdxTitle = iota
	focusIdxDesc
	focusIdxDate
	focusIdxMax

	// Padding to indent the datepicker calendar within the form.
	formCalendarPadding = 10
)

// Form is a Bubble Tea sub-model that encapsulates the editable task
type Form struct {
	Title    textinput.Model
	Desc     textinput.Model
	Date     datepicker.Model
	Done     bool
	focusIdx int
	keymap   *EditTaskKeyMap

	styles Styles
}

// NewForm constructs a Form with initial values for the task fields,
// defaulting the due date to now if it is zero.
func NewForm(
	title, desc string,
	dueDate time.Time,
	done bool,
	keymap *EditTaskKeyMap,
	styles Styles,
) Form {
	// Default due date to "now" if zero so the datepicker has a sensible value.
	if dueDate.IsZero() {
		dueDate = time.Now()
	}

	dp := datepicker.NewWithRange(dueDate, dueDate, time.Time{})

	return Form{
		Title:    newTitleInput(title),
		Desc:     newDescInput(desc),
		Date:     dp,
		Done:     done,
		focusIdx: focusIdxTitle,
		keymap:   keymap,
		styles:   styles,
	}
}

// newTitleInput configures a text input for the title field.
func newTitleInput(initial string) textinput.Model {
	ti := textinput.New()
	ti.Prompt = defaultTitlePrompt
	ti.PromptStyle.Underline(true)
	ti.Placeholder = defaultTitlePlaceholder
	ti.SetValue(initial)
	ti.SetCursor(len(initial))
	ti.Focus()
	return ti
}

// newDescInput configures a text input for the description field.
func newDescInput(initial string) textinput.Model {
	ti := textinput.New()
	ti.Prompt = defaultDescPrompt
	ti.PromptStyle.Underline(true)
	ti.Placeholder = defaultDescPlaceholder
	ti.SetValue(initial)
	ti.SetCursor(len(initial))
	return ti
}

// Update processes incoming messages for the form, cycling focus on
// SaveField and delegating updates to each sub-input.
func (f Form) Update(msg tea.Msg) (Form, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, f.keymap.SaveField):
			f.focusIdx = (f.focusIdx + 1) % focusIdxMax
			f = f.setFocus()
		}
	}

	var cmd tea.Cmd

	switch f.focusIdx {
	case focusIdxTitle:
		f.Title, cmd = f.Title.Update(msg)
		cmds = append(cmds, cmd)

	case focusIdxDesc:
		f.Desc, cmd = f.Desc.Update(msg)
		cmds = append(cmds, cmd)

	case focusIdxDate:
		f.Date, cmd = f.Date.Update(msg)
		cmds = append(cmds, cmd)
	}

	return f, tea.Batch(cmds...)
}

// setFocus applies focus to the active field based on the current
// focus index and blurs all others.
func (f Form) setFocus() Form {
	f.Title.Blur()
	f.Desc.Blur()
	f.Date.Blur()

	switch f.focusIdx {
	case focusIdxTitle:
		f.Title.Focus()
	case focusIdxDesc:
		f.Desc.Focus()
	case focusIdxDate:
		f.Date.SelectDate()
		f.Date.SetFocus(datepicker.FocusCalendar)
	}

	return f
}

// View renders the form
func (f Form) View() string {
	f.Title.TextStyle = f.styles.Normal
	f.Title.PromptStyle = f.styles.Normal
	f.Desc.TextStyle = f.styles.Normal
	f.Desc.PromptStyle = f.styles.Normal

	// Base calendar style: keep a fixed left padding so the calendar
	// is horizontally aligned regardless of focus state.
	calendarStyle := lipgloss.NewStyle().
		AlignHorizontal(lipgloss.Center).
		PaddingLeft(formCalendarPadding)

	switch f.focusIdx {
	case focusIdxTitle:
		f.Title.PromptStyle = f.styles.Focused
	case focusIdxDesc:
		f.Desc.PromptStyle = f.styles.Focused
	case focusIdxDate:
		calendarStyle = calendarStyle.
			Border(lipgloss.NormalBorder(), false, false, false, true).
			BorderForeground(f.styles.Focused.GetForeground())
	}

	calendar := calendarStyle.Render(f.Date.View())

	return lipgloss.JoinVertical(
		lipgloss.Left,
		f.Title.View(),
		f.Desc.View(),
		calendar,
	)
}
