package storage

import "github.com/jacobdanielrose/terminaltask/internal/ui/task"

type TaskStore interface {
	Load() ([]task.Task, error)
	Save([]task.Task) error
}
