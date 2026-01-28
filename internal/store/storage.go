package storage

import "github.com/jacobdanielrose/terminaltask/internal/task"

type TaskStore interface {
	Load() ([]task.Task, error)
	Save([]task.Task) error
}
