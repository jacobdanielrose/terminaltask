package store

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/jacobdanielrose/terminaltask/internal/task"
)

func TestNewFileStore(t *testing.T) {

	tfs := NewFileTaskStore(".")

	if tfs == nil {
		t.Errorf("Expected non-nil FileTaskStore")
	}

}

func TestFileStoreSave(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "tasks.json")

	store := NewFileTaskStore(path)

	// Fixed times for deterministic tests
	t1 := time.Date(2024, 1, 2, 15, 4, 5, 0, time.UTC)
	t2 := time.Date(2024, 1, 3, 10, 0, 0, 0, time.UTC)

	tasks := []task.Task{
		task.NewWithOptions("Task 1", "Description 1", t1, false),
		task.NewWithOptions("Task 2", "Description 2", t2, true),
	}

	if err := store.Save(tasks); err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if _, err := os.Stat(path); err != nil {
		t.Fatalf("expected file %s to exist after Save, got error: %v", path, err)
	}

	loadedStore := NewFileTaskStore(path)
	loaded, err := loadedStore.Load()
	if err != nil {
		t.Fatalf("Load() returned error: %v", err)
	}

	if len(loaded) != len(tasks) {
		t.Fatalf("expected %d tasks after Load, got %d", len(tasks), len(loaded))
	}
	for i, tk := range tasks {
		if loaded[i].TitleStr != tk.TitleStr {
			t.Fatalf("expected task %d title to be %s, got %s", i, tk.TitleStr, loaded[i].TitleStr)
		}
		if loaded[i].DescStr != tk.DescStr {
			t.Fatalf("expected task %d description to be %s, got %s", i, tk.DescStr, loaded[i].DescStr)
		}
		if loaded[i].DueDate != tk.DueDate {
			t.Fatalf("expected task %d due time to be %v, got %v", i, tk.DueDate, loaded[i].DueDate)
		}
		if loaded[i].Done != tk.Done {
			t.Fatalf("expected task %d completed status to be %v, got %v", i, tk.Done, loaded[i].Done)
		}
	}

}
