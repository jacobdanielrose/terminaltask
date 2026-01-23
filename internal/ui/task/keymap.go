package task

import "github.com/charmbracelet/bubbles/key"

type TaskKeyMap struct {
	EditItem   key.Binding
	ToggleDone key.Binding
	RemoveItem key.Binding
}

func newTaskKeyMap() *TaskKeyMap {
	return &TaskKeyMap{
		EditItem: key.NewBinding(
			key.WithKeys("e"),
			key.WithHelp("e", "edit item"),
		),
		ToggleDone: key.NewBinding(
			key.WithKeys(" "),
			key.WithHelp("space", "toggle done"),
		),
		RemoveItem: key.NewBinding(
			key.WithKeys("r"),
			key.WithHelp("r", "remove item"),
		),
	}
}

func (t TaskKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		t.ToggleDone,
		t.EditItem,
		t.RemoveItem,
	}
}

func (t TaskKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{
			t.ToggleDone,
			t.EditItem,
			t.RemoveItem,
		},
	}
}
