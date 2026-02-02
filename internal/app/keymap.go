package app

import (
	"github.com/charmbracelet/bubbles/key"
)

// listKeyMap defines key bindings for interacting with the task list.
type listKeyMap struct {
	NewItem key.Binding
	Quit    key.Binding
}

// NewListKeyMap constructs the default key bindings for the list view.
func NewListKeyMap() *listKeyMap {
	return &listKeyMap{
		NewItem: key.NewBinding(
			key.WithKeys("n"),
			key.WithHelp("n", "new item"),
		),
		Quit: key.NewBinding(
			key.WithKeys("ctrl+c"),
			key.WithHelp("ctrl+c", "quit"),
		),
	}
}
