package app

import (
	"testing"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
	"github.com/jacobdanielrose/terminaltask/internal/task"
	"github.com/jacobdanielrose/terminaltask/internal/task/editmenu"
)

func TestView_StateListUsersListView(t *testing.T) {
	delegate := task.NewTaskDelegate()
	l := list.New(nil, delegate, 0, 0)
	l = configureListModel(l, ListStyles{})

	l.SetItems([]list.Item{task.Task{TitleStr: "one"}})

	m := model{
		list:   l,
		state:  stateList,
		styles: AppStyles{Frame: lipgloss.NewStyle()},
	}

	out := m.View()

	if out == "" {
		t.Fatalf("View() return empty string, want non-empty string")
	}

	if !contains(out, "one") {
		t.Fatalf("View() output %q, want it to contain task title %q", out, "one")
	}
}

func TestView_StateEditUsesEditmenuView(t *testing.T) {
	// Create a simple task and editmenu model with a predictable title.
	tk := task.Task{TitleStr: "edit-me"}
	em := editmenu.New(tk)

	m := model{
		editmenu: em,
		state:    stateEdit,
		styles:   AppStyles{Frame: lipgloss.NewStyle()}, // no-op frame
	}

	out := m.View()
	if out == "" {
		t.Fatalf("View() returned empty string, want non-empty")
	}
	if !contains(out, "edit-me") {
		t.Fatalf("View() output %q, want it to contain edit menu title %q", out, "edit-me")
	}
}

// contains is a small helper to avoid importing strings just for tests
func contains(s, substr string) bool {
	if substr == "" {
		return true
	}
	n := len(substr)
	for i := 0; i+n <= len(s); i++ {
		if s[i:i+n] == substr {
			return true
		}
	}
	return false
}
