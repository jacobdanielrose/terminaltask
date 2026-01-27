package task

import (
	"time"

	"github.com/google/uuid"
)

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
