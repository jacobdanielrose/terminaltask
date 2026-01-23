package editmenu

import "github.com/charmbracelet/bubbles/key"

type EditTaskKeyMap struct {
	SaveField      key.Binding
	EscapeEditMode key.Binding
	SaveTask       key.Binding
	Help           key.Binding
	Quit           key.Binding
}

func newEditTaskKeyMap() *EditTaskKeyMap {
	return &EditTaskKeyMap{
		SaveField: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "next field"),
		),
		EscapeEditMode: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "exit edit mode"),
		),
		SaveTask: key.NewBinding(
			key.WithKeys("ctrl+s"),
			key.WithHelp("ctrl+s", "save task"),
		),
		Help: key.NewBinding(
			key.WithKeys("ctrl+o"),
			key.WithHelp("ctrl+o", "help"),
		),
		Quit: key.NewBinding(
			key.WithKeys("ctrl+c"),
			key.WithHelp("ctrl+c", "quit"),
		),
	}
}

func (e EditTaskKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		e.SaveField,
		e.EscapeEditMode,
		e.SaveTask,
		e.Help,
	}
}

func (e EditTaskKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{
			e.SaveField,
			e.EscapeEditMode,
			e.SaveTask,
			e.Help,
			e.Quit,
		},
	}
}
