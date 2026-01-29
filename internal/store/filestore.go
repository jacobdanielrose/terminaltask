package store

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"

	"github.com/jacobdanielrose/terminaltask/internal/task"
)

const (
	DefaultName = "File Store"
)

type FileTaskStore struct {
	path string
	name string
}

func NewFileTaskStore(path string) TaskStore {
	return &FileTaskStore{path: path, name: DefaultName}
}

func (fts *FileTaskStore) Name() string {
	return fts.name
}

func (fts *FileTaskStore) Load() ([]task.Task, error) {
	b, err := os.ReadFile(fts.path)
	if errors.Is(err, os.ErrNotExist) {
		return []task.Task{}, nil
	}
	if err != nil {
		return nil, err
	}

	var tasks []task.Task
	if err := json.Unmarshal(b, &tasks); err != nil {
		return nil, err
	}
	return tasks, nil
}

func (fts *FileTaskStore) Save(tasks []task.Task) error {
	if err := os.MkdirAll(filepath.Dir(fts.path), 0o755); err != nil {
		return err
	}

	b, err := json.MarshalIndent(tasks, "", " ")
	if err != nil {
		return err
	}

	tmp := fts.path + ".tmp"
	if err := os.WriteFile(tmp, b, 0o644); err != nil {
		return err
	}
	return os.Rename(tmp, fts.path)
}
