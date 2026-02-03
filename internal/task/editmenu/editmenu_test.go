package editmenu

import (
	"strings"
	"testing"
	"time"

	"github.com/charmbracelet/bubbles/help"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/google/uuid"
	"github.com/jacobdanielrose/terminaltask/internal/task"
)

//
// Constructors / New* tests
//

func TestNewWithSizeAndStyles_EmptyTask(t *testing.T) {
	emptyTask := task.Task{}
	menuStyles := DefaultStyles()
	formStyles := DefaultStyles()

	m := NewWithSizeAndStyles(0, 0, emptyTask, menuStyles, formStyles)

	// Window title defaults when there is no task title.
	if m.Title != defaultWindowTitle {
		t.Errorf("Title = %q, want %q", m.Title, defaultWindowTitle)
	}

	// Empty task is treated as "new".
	if !m.IsNew {
		t.Errorf("IsNew = %v, want true for empty task", m.IsNew)
	}

	// When due date is zero, we expect it to be set to something non-zero in the form.
	if m.form.Date.Time.IsZero() {
		t.Errorf("expected form.Date.Time to be initialized for empty task")
	}
}

func TestNewWithSizeAndStyles_ExistingTask(t *testing.T) {
	due := time.Date(2030, 1, 2, 0, 0, 0, 0, time.UTC)
	tk := task.Task{
		ID:       uuid.New(),
		TitleStr: "my task",
		DescStr:  "desc",
		DueDate:  due,
		Done:     true,
	}
	menuStyles := DefaultStyles()
	formStyles := DefaultStyles()

	m := NewWithSizeAndStyles(80, 24, tk, menuStyles, formStyles)

	// Window title should use the task title when present.
	if m.Title != "my task" {
		t.Errorf("Title = %q, want %q", m.Title, "my task")
	}

	// Non-empty task should not be considered "new".
	if m.IsNew {
		t.Errorf("IsNew = %v, want false for non-empty task", m.IsNew)
	}

	// Form fields should reflect the task values.
	if m.form.Title.Value() != "my task" {
		t.Errorf("form.Title.Value() = %q, want %q", m.form.Title.Value(), "my task")
	}
	if m.form.Desc.Value() != "desc" {
		t.Errorf("form.Desc.Value() = %q, want %q", m.form.Desc.Value(), "desc")
	}
	if !m.form.Date.Time.Equal(due) {
		t.Errorf("form.Date.Time = %v, want %v", m.form.Date.Time, due)
	}
	if !m.form.Done {
		t.Errorf("form.Done = %v, want true", m.form.Done)
	}
}

//
// Update behavior tests
//

func TestModelUpdate_SaveTask_ValidationErrors(t *testing.T) {
	baseTask := task.Task{}
	m := New(baseTask)

	// Approximate "save task" key (ctrl+s) for tests. Exact KeyMsg shape
	// is secondary; we care that key.Matches triggers the SaveTask path.
	saveMsg := tea.KeyMsg{Type: tea.KeyCtrlS}

	// Case 1: date in the past
	m.form.Date.Time = time.Now().Add(-24 * time.Hour)
	m2, cmd := m.Update(saveMsg)
	if cmd == nil {
		t.Fatalf("expected cmd for past date validation")
	}
	if !strings.Contains(m2.statusMsg, statusMsgDatePastError) {
		t.Errorf("statusMsg = %q, want to contain %q", m2.statusMsg, statusMsgDatePastError)
	}

	// Case 2: title empty
	m.form.Date.Time = time.Now().Add(24 * time.Hour).Truncate(24 * time.Hour)
	m.form.Title.SetValue("")
	m.form.Desc.SetValue("desc")
	m2, _ = m.Update(saveMsg)
	if !strings.Contains(m2.statusMsg, statusMsgTitleEmptyError) {
		t.Errorf("statusMsg = %q, want to contain %q", m2.statusMsg, statusMsgTitleEmptyError)
	}

	// Case 3: desc empty
	m.form.Title.SetValue("title")
	m.form.Desc.SetValue("")
	m2, _ = m.Update(saveMsg)
	if !strings.Contains(m2.statusMsg, statusMsgDescEmptyError) {
		t.Errorf("statusMsg = %q, want to contain %q", m2.statusMsg, statusMsgDescEmptyError)
	}
}

func TestModelUpdate_SaveTask_Success(t *testing.T) {
	baseTask := task.Task{ID: uuid.New()}
	m := New(baseTask)

	// Satisfy all validation conditions.
	m.form.Date.Time = time.Now().Add(24 * time.Hour).Truncate(24 * time.Hour)
	m.form.Title.SetValue("title")
	m.form.Desc.SetValue("desc")
	m.form.Done = true

	// Approximation of pressing "ctrl+s".
	saveMsg := tea.KeyMsg{Type: tea.KeyCtrlS}

	m2, cmd := m.Update(saveMsg)
	if cmd == nil {
		t.Fatalf("expected non-nil cmd on successful save")
	}

	out := cmd()
	save, ok := out.(SaveTaskMsg)
	if !ok {
		t.Fatalf("expected SaveTaskMsg, got %T", out)
	}

	if save.Title != "title" {
		t.Errorf("Title = %q, want %q", save.Title, "title")
	}
	if save.Desc != "desc" {
		t.Errorf("Desc = %q, want %q", save.Desc, "desc")
	}
	if !save.Date.Equal(m2.form.Date.Time) {
		t.Errorf("Date in message = %v, want %v", save.Date, m2.form.Date.Time)
	}
	if !save.Done {
		t.Errorf("Done = %v, want true", save.Done)
	}
	if save.IsNew != m2.IsNew {
		t.Errorf("IsNew in message = %v, want %v", save.IsNew, m2.IsNew)
	}
}

func TestModelUpdate_EscapeEditMode(t *testing.T) {
	m := New(task.Task{})

	// Approximate an escape key; key.Matches against the binding should
	// drive the EscapeEditMode branch.
	escMsg := tea.KeyMsg{Type: tea.KeyEscape}

	m2, cmd := m.Update(escMsg)
	_ = m2 // model can remain largely unchanged

	if cmd == nil {
		t.Fatalf("expected non-nil cmd for escape")
	}
	if _, ok := cmd().(EscapeEditMsg); !ok {
		t.Fatalf("expected EscapeEditMsg from escape key")
	}
}

func TestModelUpdate_ToggleHelp(t *testing.T) {
	m := New(task.Task{})
	m.help = help.New()

	// Approximate the help key (ctrl+o).
	helpMsg := tea.KeyMsg{Type: tea.KeyCtrlO}

	initial := m.help.ShowAll

	m2, _ := m.Update(helpMsg)
	if m2.help.ShowAll == initial {
		t.Fatalf("expected help.ShowAll to toggle")
	}
}

func TestModelUpdate_ClearStatusMsg(t *testing.T) {
	m := New(task.Task{})
	m.statusMsg = "something"

	m2, cmd := m.Update(clearStatusMsg{})
	if cmd != nil {
		t.Fatalf("expected nil cmd for clearStatusMsg, got non-nil")
	}
	if m2.statusMsg != "" {
		t.Errorf("statusMsg = %q, want empty after clearStatusMsg", m2.statusMsg)
	}
}

//
// Setters / getters / view tests
//

func TestModelSetSizeAndFlags(t *testing.T) {
	m := New(task.Task{})

	m = m.SetSize(80, 24)
	if m.Width() != 80 || m.Height() != 24 {
		t.Errorf("Width,Height = (%d,%d), want (80,24)", m.Width(), m.Height())
	}
	if m.help.Width != 80 {
		t.Errorf("help.Width = %d, want 80", m.help.Width)
	}

	m = m.SetShowTitle(false)
	if m.ShowTitle() {
		t.Errorf("ShowTitle() = true, want false")
	}

	m = m.SetShowHelp(false)
	if m.ShowHelp() {
		t.Errorf("ShowHelp() = true, want false")
	}
}

func TestModelView_IncludesStatusMessage(t *testing.T) {
	m := New(task.Task{})
	m = m.SetSize(80, 10)
	m.statusMsg = "some error"

	view := m.View()
	if !strings.Contains(view, "some error") {
		t.Errorf("View() = %q, want to contain status message %q", view, "some error")
	}
}
