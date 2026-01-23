package task

import (
	"time"
)

type Task struct {
	TitleStr string
	DescStr  string
	DueDate  time.Time
	Done     bool
}

func (t Task) FilterValue() string { return t.TitleStr }
func (t Task) Title() string       { return t.TitleStr }
func (t Task) Description() string { return t.DescStr }
