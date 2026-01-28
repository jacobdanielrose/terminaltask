package app

import (
	"github.com/charmbracelet/bubbles/key"
)

type listKeyMap struct {
	NewItem key.Binding
}

func NewListKeyMap() *listKeyMap {
	return &listKeyMap{
		NewItem: key.NewBinding(
			key.WithKeys("n"),
			key.WithHelp("n", "new item"),
		),
	}
}
