package app

import (
	"errors"
	"testing"

	"github.com/google/uuid"
	taskservice "github.com/jacobdanielrose/terminaltask/internal/service"
	task "github.com/jacobdanielrose/terminaltask/internal/task"
)

// commandsFakeService is a dedicated test double for commands.go tests.
// It implements taskservice.Service without conflicting with the fakeService
// defined in model_test.go.
type commandsFakeService struct {
	saveTasksFn  func([]task.Task) error
	loadTasksFn  func() ([]task.Task, error)
	toggleFn     func(t task.Task) (task.Task, error)
	deleteByIDFn func(id uuid.UUID) error
	upsertFn     func(t task.Task) error
	nameFn       func() string
}

// Ensure commandsFakeService satisfies the Service interface.
var _ taskservice.Service = (*commandsFakeService)(nil)

func (f *commandsFakeService) LoadTasks() ([]task.Task, error) {
	if f.loadTasksFn != nil {
		return f.loadTasksFn()
	}
	return nil, nil
}

func (f *commandsFakeService) SaveTasks(tasks []task.Task) error {
	if f.saveTasksFn != nil {
		return f.saveTasksFn(tasks)
	}
	return nil
}

func (f *commandsFakeService) ToggleCompleted(t task.Task) (task.Task, error) {
	if f.toggleFn != nil {
		return f.toggleFn(t)
	}
	return t, nil
}

func (f *commandsFakeService) DeleteByID(id uuid.UUID) error {
	if f.deleteByIDFn != nil {
		return f.deleteByIDFn(id)
	}
	return nil
}

func (f *commandsFakeService) UpsertTask(t task.Task) error {
	if f.upsertFn != nil {
		return f.upsertFn(t)
	}
	return nil
}

func (f *commandsFakeService) Name() string {
	if f.nameFn != nil {
		return f.nameFn()
	}
	return "commands-fake"
}

func TestSaveTasksCmd_Success(t *testing.T) {
	tasks := []task.Task{
		task.NewWithOptions("t1", "d1", task.Task{}.DueDate, false),
	}

	var saved []task.Task
	m := Model{
		service: &commandsFakeService{
			saveTasksFn: func(ts []task.Task) error {
				saved = append(saved, ts...)
				return nil
			},
		},
	}

	const msgText = "saved!"
	cmd := m.saveTasksCmd(tasks, msgText)
	if cmd == nil {
		t.Fatalf("expected non-nil cmd from saveTasksCmd")
	}

	out := cmd()
	savedMsg, ok := out.(TasksSavedMsg)
	if !ok {
		t.Fatalf("expected TasksSavedMsg, got %T", out)
	}

	// Assert the message content.
	if savedMsg.msg != msgText {
		t.Errorf("TasksSavedMsg.msg = %q, want %q", savedMsg.msg, msgText)
	}

	// Assert the service was called with the correct tasks.
	if len(saved) != len(tasks) {
		t.Fatalf("saved %d tasks, want %d", len(saved), len(tasks))
	}
	for i := range tasks {
		if saved[i].Title() != tasks[i].Title() {
			t.Errorf("saved[%d].Title() = %q, want %q", i, saved[i].Title(), tasks[i].Title())
		}
	}
}

func TestSaveTasksCmd_Error(t *testing.T) {
	saveErr := errors.New("save failed")

	m := Model{
		service: &commandsFakeService{
			saveTasksFn: func([]task.Task) error {
				return saveErr
			},
		},
	}

	cmd := m.saveTasksCmd(nil, "ignored")
	if cmd == nil {
		t.Fatalf("expected non-nil cmd from saveTasksCmd")
	}

	out := cmd()
	errMsg, ok := out.(TasksSaveErrorMsg)
	if !ok {
		t.Fatalf("expected TasksSaveErrorMsg, got %T", out)
	}
	if errMsg.Err != saveErr {
		t.Errorf("TasksSaveErrorMsg.Err = %v, want %v", errMsg.Err, saveErr)
	}
}

func TestLoadTasksCmd_Success(t *testing.T) {
	expectedTasks := []task.Task{
		task.NewWithOptions("t1", "d1", task.Task{}.DueDate, false),
		task.NewWithOptions("t2", "d2", task.Task{}.DueDate, true),
	}

	m := Model{
		service: &commandsFakeService{
			loadTasksFn: func() ([]task.Task, error) {
				return expectedTasks, nil
			},
		},
	}

	cmd := m.loadTasksCmd()
	if cmd == nil {
		t.Fatalf("expected non-nil cmd from loadTasksCmd")
	}

	out := cmd()
	loadedMsg, ok := out.(TasksLoadedMsg)
	if !ok {
		t.Fatalf("expected TasksLoadedMsg, got %T", out)
	}

	if len(loadedMsg.Tasks) != len(expectedTasks) {
		t.Fatalf("TasksLoadedMsg.Tasks length = %d, want %d", len(loadedMsg.Tasks), len(expectedTasks))
	}
	for i := range expectedTasks {
		if loadedMsg.Tasks[i].Title() != expectedTasks[i].Title() {
			t.Errorf("TasksLoadedMsg.Tasks[%d].Title() = %q, want %q",
				i, loadedMsg.Tasks[i].Title(), expectedTasks[i].Title())
		}
	}
}

func TestLoadTasksCmd_Error(t *testing.T) {
	loadErr := errors.New("load failed")

	m := Model{
		service: &commandsFakeService{
			loadTasksFn: func() ([]task.Task, error) {
				return nil, loadErr
			},
		},
	}

	cmd := m.loadTasksCmd()
	if cmd == nil {
		t.Fatalf("expected non-nil cmd from loadTasksCmd")
	}

	out := cmd()
	errMsg, ok := out.(TasksLoadErrorMsg)
	if !ok {
		t.Fatalf("expected TasksLoadErrorMsg, got %T", out)
	}
	if errMsg.Err != loadErr {
		t.Errorf("TasksLoadErrorMsg.Err = %v, want %v", errMsg.Err, loadErr)
	}
}
