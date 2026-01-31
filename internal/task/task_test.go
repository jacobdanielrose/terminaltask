package task

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

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
