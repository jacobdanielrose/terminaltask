package task

import (
	"bytes"
	"testing"
	"time"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

// Helpers

// newTestList creates a minimal list.Model using the TaskDelegate.
func newTestList(items []list.Item, d list.ItemDelegate) list.Model {
	m := list.New(items, d, 0, 0)
	// Width will be set explicitly in tests that care about it.
	return m
}

// Tests

func TestNewTaskDelegateAccessors(t *testing.T) {
	delegate := NewTaskDelegate()

	if delegate.Height() != defaultDelegateHeight {
		t.Errorf("NewTaskDelegate.Height() = %d, want %d", delegate.Height(), defaultDelegateHeight)
	}

	if delegate.Spacing() != defaultDelegateSpacing {
		t.Errorf("NewTaskDelegate.Spacing() = %d, want %d", delegate.Spacing(), defaultDelegateSpacing)
	}

	// Smoke check: styles should be usable (render something non-empty).
	out := delegate.Styles.Normal.Title.Render("x")
	if out == "" {
		t.Errorf("expected Normal.Title style to render non-empty string")
	}
}

// Update loop tests

func TestTaskDelegateUpdate_ToggleDone(t *testing.T) {
	d := NewTaskDelegate()
	items := []list.Item{Task{TitleStr: "title"}}
	m := newTestList(items, d)

	// Current key binding for ToggleDone is " " (space).
	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{' '}}

	cmd := d.Update(msg, &m)
	if cmd == nil {
		t.Fatalf("expected non-nil cmd for ToggleDone key")
	}

	out := cmd()
	if _, ok := out.(ToggleDoneMsg); !ok {
		t.Fatalf("expected ToggleDoneMsg, got %T", out)
	}
}

func TestTaskDelegateUpdate_EnterEdit(t *testing.T) {
	d := NewTaskDelegate()
	items := []list.Item{Task{TitleStr: "title"}}
	m := newTestList(items, d)

	// Current key binding for EditItem is "e".
	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}}

	cmd := d.Update(msg, &m)
	if cmd == nil {
		t.Fatalf("expected non-nil cmd for Edit key")
	}

	out := cmd()
	if _, ok := out.(EnterEditMsg); !ok {
		t.Fatalf("expected EnterEditMsg, got %T", out)
	}
}

func TestTaskDelegateUpdate_Delete_WithItems(t *testing.T) {
	d := NewTaskDelegate()
	items := []list.Item{Task{TitleStr: "title"}}
	m := newTestList(items, d)

	// Current key binding for RemoveItem is "r".
	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'r'}}

	cmd := d.Update(msg, &m)
	if cmd == nil {
		t.Fatalf("expected non-nil cmd for Remove key when items exist")
	}

	// RemoveItem should remain enabled when items exist.
	if !d.keymap.RemoveItem.Enabled() {
		t.Fatalf("expected RemoveItem to remain enabled when items exist")
	}

	out := cmd()
	if _, ok := out.(DeleteMsg); !ok {
		t.Fatalf("expected DeleteMsg, got %T", out)
	}
}

func TestTaskDelegateUpdate_Delete_NoItemsDisablesRemove(t *testing.T) {
	d := NewTaskDelegate()
	items := []list.Item{} // empty list
	m := newTestList(items, d)

	// Current key binding for RemoveItem is "r".
	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'r'}}

	cmd := d.Update(msg, &m)
	if cmd == nil {
		t.Fatalf("expected non-nil cmd for Remove key even when no items")
	}

	// RemoveItem should be disabled when there are no items.
	if d.keymap.RemoveItem.Enabled() {
		t.Fatalf("expected RemoveItem to be disabled when there are no items")
	}

	out := cmd()
	if _, ok := out.(DeleteMsg); !ok {
		t.Fatalf("expected DeleteMsg, got %T", out)
	}
}

func TestTaskDelegateUpdate_IgnoresOtherKeys(t *testing.T) {
	d := NewTaskDelegate()
	items := []list.Item{Task{TitleStr: "title"}}
	m := newTestList(items, d)

	// Some unrelated key we don't bind.
	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'x'}}

	cmd := d.Update(msg, &m)
	if cmd != nil {
		t.Fatalf("expected nil cmd for unhandled key, got non-nil")
	}
}

func TestTaskDelegateUpdate_IgnoresNonKeyMsgs(t *testing.T) {
	d := NewTaskDelegate()
	items := []list.Item{Task{TitleStr: "title"}}
	m := newTestList(items, d)

	// A non-key message should be ignored.
	cmd := d.Update(tea.WindowSizeMsg{Width: 80, Height: 24}, &m)
	if cmd != nil {
		t.Fatalf("expected nil cmd for non-key msg, got non-nil")
	}
}

// Render tests

func TestTaskDelegateRender_NonTaskItem_NoOutput(t *testing.T) {
	d := NewTaskDelegate()
	m := newTestList(nil, d)
	m.SetWidth(40)

	var buf bytes.Buffer
	// Use a dummy list.Item that is not a Task.
	type otherItem struct{ list.Item }
	var it otherItem

	d.Render(&buf, m, 0, it)

	if buf.Len() != 0 {
		t.Fatalf("expected no output when rendering non-Task item, got %q", buf.String())
	}
}

func TestTaskDelegateRender_ZeroWidth_NoOutput(t *testing.T) {
	d := NewTaskDelegate()
	items := []list.Item{Task{TitleStr: "title", DescStr: "desc", DueDate: time.Now()}}
	m := newTestList(items, d)
	m.SetWidth(0)

	var buf bytes.Buffer
	d.Render(&buf, m, 0, items[0])

	if buf.Len() != 0 {
		t.Fatalf("expected no output when model width &#60;= 0, got %q", buf.String())
	}
}

// Help bindings from the delegate

func TestTaskDelegateShortHelpBindings(t *testing.T) {
	d := NewTaskDelegate()

	short := d.ShortHelp()
	if len(short) != 3 {
		t.Fatalf("ShortHelp length = %d, want 3", len(short))
	}

	// We rely on key order matching the current keymap:
	// [ToggleDone (" "), EditItem ("e"), RemoveItem ("r")].
	if keys := short[0].Keys(); len(keys) == 0 || keys[0] != " " {
		t.Errorf("ShortHelp[0].Keys() = %v, want first key %q", keys, " ")
	}
	if keys := short[1].Keys(); len(keys) == 0 || keys[0] != "e" {
		t.Errorf("ShortHelp[1].Keys() = %v, want first key %q", keys, "e")
	}
	if keys := short[2].Keys(); len(keys) == 0 || keys[0] != "r" {
		t.Errorf("ShortHelp[2].Keys() = %v, want first key %q", keys, "r")
	}
}

func TestTaskDelegateFullHelpBindings(t *testing.T) {
	d := NewTaskDelegate()

	full := d.FullHelp()
	if len(full) != 1 {
		t.Fatalf("FullHelp length = %d, want 1 row", len(full))
	}
	if len(full[0]) != 3 {
		t.Fatalf("FullHelp[0] length = %d, want 3", len(full[0]))
	}

	row := full[0]

	if keys := row[0].Keys(); len(keys) == 0 || keys[0] != " " {
		t.Errorf("FullHelp[0][0].Keys() = %v, want first key %q", keys, " ")
	}
	if keys := row[1].Keys(); len(keys) == 0 || keys[0] != "e" {
		t.Errorf("FullHelp[0][1].Keys() = %v, want first key %q", keys, "e")
	}
	if keys := row[2].Keys(); len(keys) == 0 || keys[0] != "r" {
		t.Errorf("FullHelp[0][2].Keys() = %v, want first key %q", keys, "r")
	}
}
