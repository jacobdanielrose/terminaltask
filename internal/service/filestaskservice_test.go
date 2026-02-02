package taskservice

import (
	"errors"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/jacobdanielrose/terminaltask/internal/task"
)

// mockStore is an in-memory implementation of store.TaskStore used for testing
// FileTaskService behavior without touching the filesystem.
type mockStore struct {
	name      string
	tasks     []task.Task
	loadErr   error
	saveErr   error
	saveCalls int
	lastSaved []task.Task
}

func newMockStore(name string, tasks []task.Task) *mockStore {
	// Make a copy so tests don't share slices accidentally.
	cp := make([]task.Task, len(tasks))
	copy(cp, tasks)

	return &mockStore{
		name:  name,
		tasks: cp,
	}
}

// TaskStore interface methods (from internal/store, repeated conceptually here):
// Load() ([]task.Task, error)
// Save([]task.Task) error
// Name() string

func (m *mockStore) Load() ([]task.Task, error) {
	if m.loadErr != nil {
		return nil, m.loadErr
	}
	cp := make([]task.Task, len(m.tasks))
	copy(cp, m.tasks)
	return cp, nil
}

func (m *mockStore) Save(tasks []task.Task) error {
	m.saveCalls++
	if m.saveErr != nil {
		return m.saveErr
	}
	cp := make([]task.Task, len(tasks))
	copy(cp, tasks)
	m.tasks = cp
	m.lastSaved = cp
	return nil
}

func (m *mockStore) Name() string {
	if m.name == "" {
		return "mockStore"
	}
	return m.name
}

// -----------------------------------------------------------------------------
// Helpers
// -----------------------------------------------------------------------------

func newTaskWithID(id uuid.UUID, title string, done bool) task.Task {
	t := task.Task{
		TitleStr: title,
		Done:     done,
	}
	t.SetID(id)
	return t
}

// -----------------------------------------------------------------------------
// Basic delegation tests
// -----------------------------------------------------------------------------

func TestFileTaskService_NameDelegatesToStore(t *testing.T) {
	ms := newMockStore("Test Store", nil)
	svc := NewFileTaskService(ms)

	if got, want := svc.Name(), "Test Store"; got != want {
		t.Fatalf("Name() = %q, want %q", got, want)
	}
}

func TestFileTaskService_LoadTasksDelegatesToStore(t *testing.T) {
	id := uuid.New()
	expected := []task.Task{newTaskWithID(id, "one", false)}
	ms := newMockStore("mock", expected)
	svc := NewFileTaskService(ms)

	got, err := svc.LoadTasks()
	if err != nil {
		t.Fatalf("LoadTasks() error = %v, want nil", err)
	}

	if len(got) != len(expected) {
		t.Fatalf("len(LoadTasks()) = %d, want %d", len(got), len(expected))
	}
	if got[0].TitleStr != expected[0].TitleStr {
		t.Errorf("first task TitleStr = %q, want %q", got[0].TitleStr, expected[0].TitleStr)
	}
}

func TestFileTaskService_SaveTasksDelegatesToStore(t *testing.T) {
	id := uuid.New()
	input := []task.Task{newTaskWithID(id, "save me", false)}
	ms := newMockStore("mock", nil)
	svc := NewFileTaskService(ms)

	if err := svc.SaveTasks(input); err != nil {
		t.Fatalf("SaveTasks() error = %v, want nil", err)
	}

	if ms.saveCalls != 1 {
		t.Fatalf("Save() was called %d times, want 1", ms.saveCalls)
	}
	if len(ms.tasks) != len(input) {
		t.Fatalf("store.tasks len = %d, want %d", len(ms.tasks), len(input))
	}
	if ms.tasks[0].TitleStr != input[0].TitleStr {
		t.Errorf("saved task TitleStr = %q, want %q", ms.tasks[0].TitleStr, input[0].TitleStr)
	}
}

// -----------------------------------------------------------------------------
// ToggleCompleted
// -----------------------------------------------------------------------------

func TestFileTaskService_ToggleCompleted_TogglesAndPersists(t *testing.T) {
	id := uuid.New()
	orig := newTaskWithID(id, "toggle me", false)
	ms := newMockStore("mock", []task.Task{orig})
	svc := NewFileTaskService(ms)

	updated, err := svc.ToggleCompleted(orig)
	if err != nil {
		t.Fatalf("ToggleCompleted() error = %v, want nil", err)
	}

	// Returned value toggled
	if !updated.Done {
		t.Fatalf("updated.Done = %v, want true", updated.Done)
	}

	// Store state updated and saved once
	if ms.saveCalls != 1 {
		t.Fatalf("Save() was called %d times, want 1", ms.saveCalls)
	}
	if len(ms.tasks) != 1 {
		t.Fatalf("store.tasks len = %d, want 1", len(ms.tasks))
	}
	if ms.tasks[0].GetID() != id {
		t.Fatalf("store.tasks[0].ID = %v, want %v", ms.tasks[0].GetID(), id)
	}
	if !ms.tasks[0].Done {
		t.Fatalf("store.tasks[0].Done = %v, want true", ms.tasks[0].Done)
	}
}

func TestFileTaskService_ToggleCompleted_PropagatesLoadError(t *testing.T) {
	id := uuid.New()
	orig := newTaskWithID(id, "broken load", false)
	ms := newMockStore("mock", nil)
	ms.loadErr = errors.New("load failed")
	svc := NewFileTaskService(ms)

	_, err := svc.ToggleCompleted(orig)
	if err == nil {
		t.Fatalf("ToggleCompleted() error = nil, want non-nil")
	}
	if wantSub := "load tasks"; !containsSubstr(err.Error(), wantSub) {
		t.Errorf("ToggleCompleted() error = %q, want substring %q", err.Error(), wantSub)
	}
}

func TestFileTaskService_ToggleCompleted_PropagatesSaveError(t *testing.T) {
	id := uuid.New()
	orig := newTaskWithID(id, "broken save", false)
	ms := newMockStore("mock", []task.Task{orig})
	ms.saveErr = errors.New("save failed")
	svc := NewFileTaskService(ms)

	_, err := svc.ToggleCompleted(orig)
	if err == nil {
		t.Fatalf("ToggleCompleted() error = nil, want non-nil")
	}
	if wantSub := "save tasks"; !containsSubstr(err.Error(), wantSub) {
		t.Errorf("ToggleCompleted() error = %q, want substring %q", err.Error(), wantSub)
	}
}

// -----------------------------------------------------------------------------
// DeleteByID
// -----------------------------------------------------------------------------

func TestFileTaskService_DeleteByID_RemovesMatchingTask(t *testing.T) {
	id1 := uuid.New()
	id2 := uuid.New()
	t1 := newTaskWithID(id1, "keep", false)
	t2 := newTaskWithID(id2, "delete me", false)
	ms := newMockStore("mock", []task.Task{t1, t2})
	svc := NewFileTaskService(ms)

	if err := svc.DeleteByID(id2); err != nil {
		t.Fatalf("DeleteByID() error = %v, want nil", err)
	}

	if ms.saveCalls != 1 {
		t.Fatalf("Save() was called %d times, want 1", ms.saveCalls)
	}
	if len(ms.tasks) != 1 {
		t.Fatalf("store.tasks len = %d, want 1", len(ms.tasks))
	}
	if ms.tasks[0].GetID() != id1 {
		t.Fatalf("remaining task ID = %v, want %v", ms.tasks[0].GetID(), id1)
	}
}

func TestFileTaskService_DeleteByID_LoadError(t *testing.T) {
	ms := newMockStore("mock", nil)
	ms.loadErr = errors.New("load failed")
	svc := NewFileTaskService(ms)

	err := svc.DeleteByID(uuid.New())
	if err == nil {
		t.Fatalf("DeleteByID() error = nil, want non-nil")
	}
	if wantSub := "load tasks"; !containsSubstr(err.Error(), wantSub) {
		t.Errorf("DeleteByID() error = %q, want substring %q", err.Error(), wantSub)
	}
}

func TestFileTaskService_DeleteByID_SaveError(t *testing.T) {
	id := uuid.New()
	t1 := newTaskWithID(id, "only", false)
	ms := newMockStore("mock", []task.Task{t1})
	ms.saveErr = errors.New("save failed")
	svc := NewFileTaskService(ms)

	err := svc.DeleteByID(id)
	if err == nil {
		t.Fatalf("DeleteByID() error = nil, want non-nil")
	}
	if wantSub := "save tasks"; !containsSubstr(err.Error(), wantSub) {
		t.Errorf("DeleteByID() error = %q, want substring %q", err.Error(), wantSub)
	}
}

// -----------------------------------------------------------------------------
// UpsertTask
// -----------------------------------------------------------------------------

func TestFileTaskService_UpsertTask_UpdatesExisting(t *testing.T) {
	id := uuid.New()
	orig := newTaskWithID(id, "original", false)
	ms := newMockStore("mock", []task.Task{orig})
	svc := NewFileTaskService(ms)

	updated := newTaskWithID(id, "updated title", true)

	if err := svc.UpsertTask(updated); err != nil {
		t.Fatalf("UpsertTask() error = %v, want nil", err)
	}

	if ms.saveCalls != 1 {
		t.Fatalf("Save() was called %d times, want 1", ms.saveCalls)
	}
	if len(ms.tasks) != 1 {
		t.Fatalf("store.tasks len = %d, want 1", len(ms.tasks))
	}
	if ms.tasks[0].TitleStr != "updated title" {
		t.Errorf("stored task TitleStr = %q, want %q", ms.tasks[0].TitleStr, "updated title")
	}
	if !ms.tasks[0].Done {
		t.Errorf("stored task Done = %v, want true", ms.tasks[0].Done)
	}
}

func TestFileTaskService_UpsertTask_AppendsWhenNotFound(t *testing.T) {
	existingID := uuid.New()
	existing := newTaskWithID(existingID, "existing", false)
	ms := newMockStore("mock", []task.Task{existing})
	svc := NewFileTaskService(ms)

	newID := uuid.New()
	newTask := newTaskWithID(newID, "new task", true)

	if err := svc.UpsertTask(newTask); err != nil {
		t.Fatalf("UpsertTask() error = %v, want nil", err)
	}

	if ms.saveCalls != 1 {
		t.Fatalf("Save() was called %d times, want 1", ms.saveCalls)
	}
	if len(ms.tasks) != 2 {
		t.Fatalf("store.tasks len = %d, want 2", len(ms.tasks))
	}

	// Order should be existing, then new.
	if ms.tasks[0].GetID() != existingID {
		t.Errorf("tasks[0].ID = %v, want %v", ms.tasks[0].GetID(), existingID)
	}
	if ms.tasks[1].GetID() != newID {
		t.Errorf("tasks[1].ID = %v, want %v", ms.tasks[1].GetID(), newID)
	}
}

func TestFileTaskService_UpsertTask_LoadError(t *testing.T) {
	ms := newMockStore("mock", nil)
	ms.loadErr = errors.New("load failed")
	svc := NewFileTaskService(ms)

	err := svc.UpsertTask(task.Task{})
	if err == nil {
		t.Fatalf("UpsertTask() error = nil, want non-nil")
	}
	if wantSub := "load tasks"; !containsSubstr(err.Error(), wantSub) {
		t.Errorf("UpsertTask() error = %q, want substring %q", err.Error(), wantSub)
	}
}

func TestFileTaskService_UpsertTask_SaveError(t *testing.T) {
	id := uuid.New()
	orig := newTaskWithID(id, "orig", false)
	ms := newMockStore("mock", []task.Task{orig})
	ms.saveErr = errors.New("save failed")
	svc := NewFileTaskService(ms)

	err := svc.UpsertTask(orig)
	if err == nil {
		t.Fatalf("UpsertTask() error = nil, want non-nil")
	}
	if wantSub := "save tasks"; !containsSubstr(err.Error(), wantSub) {
		t.Errorf("UpsertTask() error = %q, want substring %q", err.Error(), wantSub)
	}
}

// -----------------------------------------------------------------------------
// Small helpers
// -----------------------------------------------------------------------------

func containsSubstr(s, substr string) bool {
	return substr == "" || (len(s) >= len(substr) && (indexOf(s, substr) >= 0))
}

// indexOf is a minimal substring search to avoid importing strings just for tests.
func indexOf(s, substr string) int {
	n := len(substr)
	if n == 0 {
		return 0
	}
	for i := 0; i+n <= len(s); i++ {
		if s[i:i+n] == substr {
			return i
		}
	}
	return -1
}

// Sanity check to ensure mockStore satisfies the TaskStore interface at compile time.
// This mirrors the store.TaskStore shape without importing it here.
type taskStoreLike interface {
	Load() ([]task.Task, error)
	Save([]task.Task) error
	Name() string
}

func _() {
	var _ taskStoreLike = (*mockStore)(nil)
}

// Optional: ensure errors contain expected wrapping prefixes (demonstrative).
func ExampleFileTaskService_ToggleCompleted_errorWrapping() {
	ms := newMockStore("mock", nil)
	ms.loadErr = errors.New("boom")
	svc := NewFileTaskService(ms)

	_, err := svc.ToggleCompleted(task.Task{})
	fmt.Println(err != nil)
	// Output: true
}
