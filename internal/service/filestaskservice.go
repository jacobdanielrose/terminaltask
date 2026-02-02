package taskservice

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/jacobdanielrose/terminaltask/internal/store"
	"github.com/jacobdanielrose/terminaltask/internal/task"
)

type FileTaskService struct {
	store store.TaskStore
}

func NewFileTaskService(s store.TaskStore) Service {
	return &FileTaskService{store: s}
}

func (s *FileTaskService) Name() string {
	return s.store.Name()
}

func (s *FileTaskService) LoadTasks() ([]task.Task, error) {
	return s.store.Load()
}

func (s *FileTaskService) SaveTasks(tasks []task.Task) error {
	return s.store.Save(tasks)
}

func (s *FileTaskService) ToggleCompleted(t task.Task) (task.Task, error) {
	t.Done = !t.Done

	tasks, err := s.store.Load()
	if err != nil {
		return t, fmt.Errorf("load tasks: %w", err)
	}

	for i := range tasks {
		if tasks[i].GetID() == t.GetID() {
			tasks[i].Done = t.Done
			break
		}
	}

	if err := s.store.Save(tasks); err != nil {
		return t, fmt.Errorf("save tasks: %w", err)
	}

	return t, nil
}

func (s *FileTaskService) DeleteByID(id uuid.UUID) error {
	tasks, err := s.store.Load()
	if err != nil {
		return fmt.Errorf("load tasks: %w", err)
	}

	out := tasks[:0]
	for _, t := range tasks {
		if t.GetID() != id {
			out = append(out, t)
		}
	}

	if err := s.store.Save(out); err != nil {
		return fmt.Errorf("save tasks: %w", err)
	}

	return nil
}

func (s *FileTaskService) UpsertTask(t task.Task) error {
	tasks, err := s.store.Load()
	if err != nil {
		return fmt.Errorf("load tasks: %w", err)
	}

	found := false
	for i := range tasks {
		if tasks[i].GetID() == t.GetID() {
			tasks[i] = t
			found = true
			break
		}
	}

	if !found {
		tasks = append(tasks, t)
	}

	if err := s.store.Save(tasks); err != nil {
		return fmt.Errorf("save tasks: %w", err)
	}

	return nil
}
