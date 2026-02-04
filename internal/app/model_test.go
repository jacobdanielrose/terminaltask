package app

import (
	"testing"
	"time"

	"github.com/charmbracelet/bubbles/list"
	"github.com/google/uuid"
	"github.com/jacobdanielrose/terminaltask/internal/config"
	"github.com/jacobdanielrose/terminaltask/internal/task"
)

// fakeService is a minimal implementation of taskservice.Service used
// to exercise NewModel wiring without hitting the real service/store.
type fakeService struct {
	name string
}

func (f *fakeService) LoadTasks() ([]task.Task, error)                { return nil, nil }
func (f *fakeService) SaveTasks([]task.Task) error                    { return nil }
func (f *fakeService) ToggleCompleted(t task.Task) (task.Task, error) { return t, nil }
func (f *fakeService) DeleteByID(uuid.UUID) error                     { return nil }
func (f *fakeService) UpsertTask(task.Task) error                     { return nil }
func (f *fakeService) Name() string                                   { return f.name }

func TestTasksToItemsAndBack(t *testing.T) {
	t1 := task.Task{TitleStr: "one", DescStr: "first"}
	t2 := task.Task{TitleStr: "two", DescStr: "second"}
	tasks := []task.Task{t1, t2}

	items := tasksToItems(tasks)
	if items == nil {
		t.Fatalf("tasksToItems returned nil slice, want non-nil")
	}
	if len(items) != len(tasks) {
		t.Fatalf("len(items) = %d, want %d", len(items), len(tasks))
	}

	roundTripped := itemsToTasks(items)
	if roundTripped == nil {
		t.Fatalf("itemsToTasks returned nil slice, want non-nil")
	}
	if len(roundTripped) != len(tasks) {
		t.Fatalf("len(roundTripped) = %d, want %d", len(roundTripped), len(tasks))
	}

	for i := range tasks {
		if roundTripped[i].TitleStr != tasks[i].TitleStr {
			t.Errorf("task %d TitleStr = %q, want %q", i, roundTripped[i].TitleStr, tasks[i].TitleStr)
		}
		if roundTripped[i].DescStr != tasks[i].DescStr {
			t.Errorf("task %d DescStr = %q, want %q", i, roundTripped[i].DescStr, tasks[i].DescStr)
		}
	}
}

func TestTasksToItemsNilAndItemsToTasksNil(t *testing.T) {
	items := tasksToItems(nil)
	if items == nil {
		t.Fatalf("tasksToItems(nil) returned nil slice, want non-nil empty slice")
	}
	if len(items) != 0 {
		t.Fatalf("len(tasksToItems(nil)) = %d, want 0", len(items))
	}

	tasks := itemsToTasks(nil)
	if tasks == nil {
		t.Fatalf("itemsToTasks(nil) returned nil slice, want non-nil empty slice")
	}
	if len(tasks) != 0 {
		t.Fatalf("len(itemsToTasks(nil)) = %d, want 0", len(tasks))
	}
}

func TestConfigureListModel(t *testing.T) {
	l := list.New(nil, list.NewDefaultDelegate(), 0, 0)
	styles := ListStyles{}
	got := configureListModel(l, styles)

	if got.Title != listModelTitle {
		t.Fatalf("Title = %q, want %q", got.Title, listModelTitle)
	}
	if !got.ShowStatusBar() {
		t.Fatalf("ShowStatusBar() = false, want true")
	}
	if got.StatusMessageLifetime <= 0 {
		t.Fatalf("StatusMessageLifetime = %v, want > 0", got.StatusMessageLifetime)
	}
}

func TestNewModelInitialState(t *testing.T) {
	cfg := config.Config{}
	svc := &fakeService{name: "fake"}
	mAny := NewModel(cfg, svc)

	// NewModel returns tea.Model; assert and inspect concrete model.
	m, ok := mAny.(Model)
	if !ok {
		t.Fatalf("NewModel returned %T, want app.model", mAny)
	}

	if m.state != stateList {
		t.Fatalf("initial state = %v, want %v (stateList)", m.state, stateList)
	}

	if m.list.Title != listModelTitle {
		t.Fatalf("list.Title = %q, want %q", m.list.Title, listModelTitle)
	}

	if m.service != svc {
		t.Fatalf("service on model not equal to provided service")
	}

	// Ensure list status lifetime is configured (not the zero value).
	if m.list.StatusMessageLifetime <= 0*time.Second {
		t.Fatalf("StatusMessageLifetime = %v, want > 0", m.list.StatusMessageLifetime)
	}
}
