package store

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/jacobdanielrose/terminaltask/internal/task"
)

// newTempStore is a small helper to create a FileTaskStore backed by a
// temporary directory. It returns the concrete *FileTaskStore for tests
// and the absolute path to the underlying file.
func newTempStore(t *testing.T, filename string) (*FileTaskStore, string) {
	t.Helper()

	dir := t.TempDir()
	path := filepath.Join(dir, filename)

	store := NewFileTaskStore(path)
	fs, ok := store.(*FileTaskStore)
	if !ok {
		t.Fatalf("NewFileTaskStore did not return *FileTaskStore, got %T", store)
	}

	return fs, path
}

// -----------------------------------------------------------------------------
// NewFileTaskStore
// -----------------------------------------------------------------------------

func TestNewFileTaskStore(t *testing.T) {
	store := NewFileTaskStore("test.json")
	if store == nil {
		t.Fatalf("NewFileTaskStore returned nil")
	}

	fs, ok := store.(*FileTaskStore)
	if !ok {
		t.Fatalf("NewFileTaskStore did not return *FileTaskStore, got %T", store)
	}

	if got, want := fs.Name(), DefaultName; got != want {
		t.Fatalf("Name() = %q, want %q", got, want)
	}
}

// -----------------------------------------------------------------------------
// Save
// -----------------------------------------------------------------------------

func TestFileTaskStore_Save_CreatesFileAndPersistsTasks(t *testing.T) {
	store, path := newTempStore(t, "tasks.json")

	// Fixed times for deterministic tests.
	t1 := time.Date(2024, 1, 2, 15, 4, 5, 0, time.UTC)
	t2 := time.Date(2024, 1, 3, 10, 0, 0, 0, time.UTC)

	tasks := []task.Task{
		task.NewWithOptions("Task 1", "Description 1", t1, false),
		task.NewWithOptions("Task 2", "Description 2", t2, true),
	}

	if err := store.Save(tasks); err != nil {
		t.Fatalf("Save() error = %v, want nil", err)
	}

	// File should exist.
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("expected file %s to exist after Save, got error: %v", path, err)
	}

	// Round-trip check via Load.
	loadedStore := NewFileTaskStore(path)
	loaded, err := loadedStore.Load()
	if err != nil {
		t.Fatalf("Load() error = %v, want nil", err)
	}

	if len(loaded) != len(tasks) {
		t.Fatalf("len(loaded) = %d, want %d", len(loaded), len(tasks))
	}

	for i, tk := range tasks {
		if loaded[i].TitleStr != tk.TitleStr {
			t.Errorf("task %d TitleStr = %q, want %q", i, loaded[i].TitleStr, tk.TitleStr)
		}
		if loaded[i].DescStr != tk.DescStr {
			t.Errorf("task %d DescStr = %q, want %q", i, loaded[i].DescStr, tk.DescStr)
		}
		if !loaded[i].DueDate.Equal(tk.DueDate) {
			t.Errorf("task %d DueDate = %v, want %v", i, loaded[i].DueDate, tk.DueDate)
		}
		if loaded[i].Done != tk.Done {
			t.Errorf("task %d Done = %v, want %v", i, loaded[i].Done, tk.Done)
		}
	}
}

func TestFileTaskStore_Save_EmptySlice(t *testing.T) {
	store, _ := newTempStore(t, "empty.json")

	if err := store.Save([]task.Task{}); err != nil {
		t.Fatalf("Save([]) error = %v, want nil", err)
	}

	loaded, err := store.Load()
	if err != nil {
		t.Fatalf("Load() error = %v, want nil", err)
	}

	if loaded == nil {
		t.Fatalf("Load() returned nil slice, want non-nil empty slice")
	}
	if len(loaded) != 0 {
		t.Fatalf("len(loaded) = %d, want 0", len(loaded))
	}
}

func TestFileTaskStore_Save_MkdirAllError(t *testing.T) {
	dir := t.TempDir()

	// Create a file that will serve as a fake "directory".
	notADir := filepath.Join(dir, "notadir")
	if err := os.WriteFile(notADir, []byte("x"), 0o644); err != nil {
		t.Fatalf("WriteFile(notadir) error = %v", err)
	}

	// Parent of this path is a file, so MkdirAll should fail.
	path := filepath.Join(notADir, "tasks.json")
	store := NewFileTaskStore(path)

	if err := store.Save([]task.Task{}); err == nil {
		t.Fatalf("Save() error = nil, want non-nil when MkdirAll fails")
	}
}

// -----------------------------------------------------------------------------
// Load
// -----------------------------------------------------------------------------

func TestFileTaskStore_Load_KnownJSON(t *testing.T) {
	store, path := newTempStore(t, "test.json")

	const sampleJSON = `[
 {
  "TitleStr": "Submit that TPS report.",
  "DescStr": "But do it like.... really slowly. ",
  "DueDate": "2026-01-31T00:00:00Z",
  "Done": true
 },
 {
  "TitleStr": "Buy groceries",
  "DescStr": "Carrots, peas, and ice cream.",
  "DueDate": "2026-02-06T00:00:00Z",
  "Done": false
 }
]`

	if err := os.WriteFile(path, []byte(sampleJSON), 0o644); err != nil {
		t.Fatalf("WriteFile error = %v", err)
	}

	got, err := store.Load()
	if err != nil {
		t.Fatalf("Load() error = %v, want nil", err)
	}

	if len(got) != 2 {
		t.Fatalf("len(got) = %d, want 2", len(got))
	}

	tests := []struct {
		got  task.Task
		want task.Task
	}{
		{
			got: got[0],
			want: task.Task{
				TitleStr: "Submit that TPS report.",
				DescStr:  "But do it like.... really slowly. ",
				DueDate:  time.Date(2026, 1, 31, 0, 0, 0, 0, time.UTC),
				Done:     true,
			},
		},
		{
			got: got[1],
			want: task.Task{
				TitleStr: "Buy groceries",
				DescStr:  "Carrots, peas, and ice cream.",
				DueDate:  time.Date(2026, 2, 6, 0, 0, 0, 0, time.UTC),
				Done:     false,
			},
		},
	}

	for i, tt := range tests {
		if tt.got.TitleStr != tt.want.TitleStr {
			t.Errorf("task %d TitleStr = %q, want %q", i, tt.got.TitleStr, tt.want.TitleStr)
		}
		if tt.got.DescStr != tt.want.DescStr {
			t.Errorf("task %d DescStr = %q, want %q", i, tt.got.DescStr, tt.want.DescStr)
		}
		if !tt.got.DueDate.Equal(tt.want.DueDate) {
			t.Errorf("task %d DueDate = %v, want %v", i, tt.got.DueDate, tt.want.DueDate)
		}
		if tt.got.Done != tt.want.Done {
			t.Errorf("task %d Done = %v, want %v", i, tt.got.Done, tt.want.Done)
		}
	}
}

func TestFileTaskStore_Load_NoFile(t *testing.T) {
	store, _ := newTempStore(t, "does_not_exist.json")

	tasks, err := store.Load()
	if err != nil {
		t.Fatalf("Load() error = %v, want nil", err)
	}
	if tasks == nil {
		t.Fatalf("Load() tasks is nil, want empty slice")
	}
	if len(tasks) != 0 {
		t.Fatalf("len(tasks) = %d, want 0", len(tasks))
	}
}

func TestFileTaskStore_Load_InvalidJSON(t *testing.T) {
	store, path := newTempStore(t, "bad.json")

	if err := os.WriteFile(path, []byte(`not valid json`), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v, want nil", err)
	}

	tasks, err := store.Load()
	if err == nil {
		t.Fatalf("Load() error = nil, want non-nil for invalid JSON")
	}
	if tasks != nil {
		t.Fatalf("Load() tasks is not nil, want nil")
	}
}

func TestFileTaskStore_Load_SaveRoundTrip(t *testing.T) {
	store, path := newTempStore(t, "tasks.json")

	t1 := time.Date(2024, 1, 2, 15, 4, 5, 0, time.UTC)
	t2 := time.Date(2024, 1, 3, 10, 0, 0, 0, time.UTC)

	original := []task.Task{
		{
			TitleStr: "task 1",
			DescStr:  "desc 1",
			DueDate:  t1,
			Done:     false,
		},
		{
			TitleStr: "task 2",
			DescStr:  "desc 2",
			DueDate:  t2,
			Done:     true,
		},
	}

	if err := store.Save(original); err != nil {
		t.Fatalf("Save() error = %v, want nil", err)
	}

	loadedStore := NewFileTaskStore(path)
	loaded, err := loadedStore.Load()
	if err != nil {
		t.Fatalf("Load() error = %v, want nil", err)
	}

	if len(loaded) != len(original) {
		t.Fatalf("len(loaded) = %d, want %d", len(loaded), len(original))
	}

	for i := range original {
		if loaded[i].TitleStr != original[i].TitleStr {
			t.Errorf("task %d TitleStr = %q, want %q", i, loaded[i].TitleStr, original[i].TitleStr)
		}
		if loaded[i].DescStr != original[i].DescStr {
			t.Errorf("task %d DescStr = %q, want %q", i, loaded[i].DescStr, original[i].DescStr)
		}
		if !loaded[i].DueDate.Equal(original[i].DueDate) {
			t.Errorf("task %d DueDate = %v, want %v", i, loaded[i].DueDate, original[i].DueDate)
		}
		if loaded[i].Done != original[i].Done {
			t.Errorf("task %d Done = %v, want %v", i, loaded[i].Done, original[i].Done)
		}
	}
}
