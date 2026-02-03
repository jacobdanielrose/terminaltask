package editmenu

import (
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	datepicker "github.com/ethanefung/bubble-datepicker"
)

//
// Construction / datepicker initialization
//

func TestNewForm_DatepickerRange_ZeroDueDate(t *testing.T) {
	var zeroDueDate time.Time
	title := "test title"
	desc := "test desc"
	done := false

	form := NewForm(title, desc, zeroDueDate, done, newEditTaskKeyMap(), Styles{})

	if isUninitializedDatepicker(form.Date) {
		t.Fatalf("expected Date to be initialized, got zero-value datepicker.Model")
	}

	current := form.Date.Time
	min, max := form.Date.StartDate, form.Date.EndDate

	// current should be approximately "now"
	now := time.Now()
	if current.Before(now.Add(-2*time.Second)) || current.After(now.Add(2*time.Second)) {
		t.Errorf("expected current date around now, got %v (now=%v)", current, now)
	}

	// min should not be zero and should be <= current
	if min.IsZero() {
		t.Errorf("expected min bound to be non-zero, got zero")
	}
	if min.After(current) {
		t.Errorf("expected min <= current, got min=%v current=%v", min, current)
	}

	// max should be zero and/or should be >= current
	if !max.IsZero() && max.Before(current) {
		t.Errorf("expected max to be zero (unbounded) or >= current, got max=%v current=%v", max, current)
	}
}

func TestNewForm_DatepickerRange_NonZeroDueDate(t *testing.T) {
	// Arrange: explicit due date
	dueDate := time.Date(2030, 1, 10, 0, 0, 0, 0, time.UTC)
	title := "future task"
	desc := "desc"
	done := false

	form := NewForm(title, desc, dueDate, done, newEditTaskKeyMap(), Styles{})

	if isUninitializedDatepicker(form.Date) {
		t.Fatalf("expected Date to be initialized, got zero-value datepicker.Model")
	}

	current := form.Date.Time
	min, max := form.Date.StartDate, form.Date.EndDate

	// current should be equal to the provided due date
	if !current.Equal(dueDate) {
		t.Errorf("expected current date %v, got %v", dueDate, current)
	}

	// min should be <= dueDate; in many designs min == dueDate
	if min.After(dueDate) {
		t.Errorf("expected min <= dueDate, got min=%v dueDate=%v", min, dueDate)
	}

	if !max.IsZero() && max.Before(dueDate) {
		t.Errorf("expected max to be zero (unbounded) or >= dueDate, got max=%v dueDate=%v", max, dueDate)
	}
}

//
// Focus handling
//

func TestNewForm_SetFocus(t *testing.T) {
	f := NewForm("title", "desc", time.Time{}, false, newEditTaskKeyMap(), Styles{})
	f = f.setFocus()

	if !f.Title.Focused() {
		t.Errorf("expected Title to be focused, got %v", f.Title.Focused())
	}

	f.focusIdx = focusIdxDesc
	f = f.setFocus()

	if !f.Desc.Focused() {
		t.Errorf("expected Desc to be focused, got %v", f.Desc.Focused())
	}

	f.focusIdx = focusIdxDate
	f = f.setFocus()

	if f.Date.Focused == datepicker.FocusNone {
		t.Errorf("expected Date to be focused")
	}
}

//
// Update / focus behavior
//

func TestNewForm_Update_FocusCycle(t *testing.T) {
	f := NewForm("title", "desc", time.Time{}, false, newEditTaskKeyMap(), Styles{})
	msg := tea.KeyMsg{Type: tea.KeyEnter}

	// make sure the index cycles properly
	if f.focusIdx != focusIdxTitle {
		t.Errorf("expected focusIdx to be focusIdxTitle, got %v", f.focusIdx)
	}

	f, _ = f.Update(msg)
	if f.focusIdx != focusIdxDesc {
		t.Errorf("expected focusIdx to be focusIdxDesc, got %v", f.focusIdx)
	}

	f, _ = f.Update(msg)
	if f.focusIdx != focusIdxDate {
		t.Errorf("expected focusIdx to be focusIdxDate, got %v", f.focusIdx)
	}

	// make sure it cycles back to the title
	f, _ = f.Update(msg)
	if f.focusIdx != focusIdxTitle {
		t.Errorf("expected focusIdx to be focusIdxTitle, got %v", f.focusIdx)
	}
}

//
// Text input configuration
//

func TestNewForm_TextInputs(t *testing.T) {
	f := NewForm("title", "description", time.Time{}, false, newEditTaskKeyMap(), Styles{})

	if f.Title.Prompt != defaultTitlePrompt {
		t.Errorf("expected Title.Prompt to be %q, got %q", defaultTitlePrompt, f.Title.Prompt)
	}
	if f.Title.Placeholder != defaultTitlePlaceholder {
		t.Errorf("expected Title.Placeholder to be %q, got %q", defaultTitlePlaceholder, f.Title.Placeholder)
	}

	if f.Desc.Prompt != defaultDescPrompt {
		t.Errorf("expected Desc.Prompt to be %q, got %q", defaultDescPrompt, f.Desc.Prompt)
	}
	if f.Desc.Placeholder != defaultDescPlaceholder {
		t.Errorf("expected Desc.Placeholder to be %q, got %q", defaultDescPlaceholder, f.Desc.Placeholder)
	}

	if f.Title.Value() != "title" {
		t.Errorf("expected Title.Value to be %q, got %q", "title", f.Title.Value())
	}

	if f.Desc.Value() != "description" {
		t.Errorf("expected Desc.Value to be %q, got %q", "description", f.Desc.Value())
	}

	if f.Title.Width != defaultTextInputWidth {
		t.Errorf("expected Title.Width to be %d, got %d", defaultTextInputWidth, f.Title.Width)
	}
	if f.Desc.Width != defaultTextInputWidth {
		t.Errorf("expected Desc.Width to be %d, got %d", defaultTextInputWidth, f.Desc.Width)
	}
}

//
// Test helpers
//

func isUninitializedDatepicker(m datepicker.Model) bool {
	if m.Focused != datepicker.Focus(0) {
		return false
	}
	if !m.Time.IsZero() {
		return false
	}
	if !m.StartDate.IsZero() {
		return false
	}
	return true
}
