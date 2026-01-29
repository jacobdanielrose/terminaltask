package taskservice

import (
	"github.com/google/uuid"
	"github.com/jacobdanielrose/terminaltask/internal/task"
)

type Service interface {
	// Loads all tasks
	LoadTasks() ([]task.Task, error)

	// Saves given tasks, overwriting existing ones
	SaveTasks(tasks []task.Task) error

	// Highlevel operations
	ToggleCompleted(t task.Task) (task.Task, error)
	DeleteByID(id uuid.UUID) error
	UpsertTask(t task.Task) error

	// For logging
	Name() string
}
