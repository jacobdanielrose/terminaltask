package task

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

// Constructors

func TestNewTaskDefaults(t *testing.T) {
	tk := New()

	if tk.TitleStr != "" {
		t.Errorf("expected default TitleStr to be empty, got %q", tk.TitleStr)
	}
	if tk.DescStr != "" {
		t.Errorf("expected default DescStr to be empty, got %q", tk.DescStr)
	}
	if !tk.DueDate.IsZero() {
		t.Errorf("expected default DueDate to be zero, got %v", tk.DueDate)
	}
	if tk.Done {
		t.Errorf("expected default Done to be false, got true")
	}

	// ID should be non-zero
	if tk.GetID() == uuid.Nil {
		t.Errorf("expected New() to assign a non-zero ID")
	}
}

func TestNewWithOptions(t *testing.T) {
	now := time.Now()

	tk := NewWithOptions("title", "desc", now, true)

	if tk.ID == uuid.Nil {
		t.Errorf("expected NewWithOptions to assign a non-zero ID")
	}
	if tk.TitleStr != "title" {
		t.Errorf("TitleStr = %q, want %q", tk.TitleStr, "title")
	}
	if tk.DescStr != "desc" {
		t.Errorf("DescStr = %q, want %q", tk.DescStr, "desc")
	}
	if !tk.DueDate.Equal(now) {
		t.Errorf("DueDate = %v, want %v", tk.DueDate, now)
	}
	if tk.Done != true {
		t.Errorf("Done = %v, want %v", tk.Done, true)
	}
}

// Accessors

func TestTaskAccessors(t *testing.T) {
	id := uuid.New()
	now := time.Now()

	tk := Task{
		ID:       id,
		TitleStr: "title",
		DescStr:  "desc",
		DueDate:  now,
		Done:     true,
	}

	if got := tk.FilterValue(); got != "title" {
		t.Errorf("FilterValue() = %q, want %q", got, "title")
	}
	if got := tk.Title(); got != "title" {
		t.Errorf("Title() = %q, want %q", got, "title")
	}
	if got := tk.Description(); got != "desc" {
		t.Errorf("Description() = %q, want %q", got, "desc")
	}
	if got := tk.GetID(); got != id {
		t.Errorf("GetID() = %v, want %v", got, id)
	}

	newID := uuid.New()
	tk.SetID(newID)
	if got := tk.ID; got != newID {
		t.Errorf("SetID() did not update ID; got %v, want %v", got, newID)
	}
}

// IsEmpty

func TestTaskIsEmpty(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name string
		task Task
		want bool
	}{
		{
			name: "all zero fields",
			task: Task{},
			want: true,
		},
		{
			name: "title set",
			task: Task{TitleStr: "title"},
			want: false,
		},
		{
			name: "description set",
			task: Task{DescStr: "desc"},
			want: false,
		},
		{
			name: "due date set",
			task: Task{DueDate: now},
			want: false,
		},
		{
			name: "done true",
			task: Task{Done: true},
			want: false,
		},
		{
			name: "mixed non-empty fields",
			task: Task{
				TitleStr: "title",
				DescStr:  "desc",
				DueDate:  now,
				Done:     true,
			},
			want: false,
		},
	}

	for _, tt := range tests {
		tt := tt // capture range variable
		t.Run(tt.name, func(t *testing.T) {
			got := tt.task.IsEmpty()
			if got != tt.want {
				t.Errorf("IsEmpty() = %v, want %v; task = %#v", got, tt.want, tt.task)
			}
		})
	}
}

// Keymap

func TestNewTaskKeyMap(t *testing.T) {
	km := newTaskKeyMap()

	if km == nil {
		t.Fatalf("newTaskKeyMap() returned nil")
	}

}

// Styles

func TestNewTaskStylesConstructsStyles(t *testing.T) {
	// Smoke test: must not panic and must return a Styles value.
	_ = newTaskStyles()
}
