package editmenu

import (
	"testing"
	"time"

	datepicker "github.com/ethanefung/bubble-datepicker"
)

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
