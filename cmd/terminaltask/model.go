package main

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	"github.com/jacobdanielrose/terminaltask/internal/storage"
	task "github.com/jacobdanielrose/terminaltask/internal/ui/task"
	"github.com/jacobdanielrose/terminaltask/internal/ui/task/editmenu"
)

type state int

const (
	stateList state = iota
	stateEdit
)

type Styles struct {
	appStyle   lipgloss.Style
	titleStyle lipgloss.Style
}

func newStyles() Styles {
	return Styles{
		appStyle: lipgloss.NewStyle().Padding(1, 2),
		titleStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFDF5")).
			Background(lipgloss.Color("#25A065")).
			Padding(0, 1),
	}
}

type model struct {
	list     list.Model
	editmenu editmenu.Model
	state    state
	keymap   *listKeyMap
	styles   Styles
	store    storage.TaskStore
}

func initialModel(store storage.TaskStore) model {

	tasks, err := store.Load()
	if err != nil {
		tasks = []task.Task{}
		log.Error("Unable to open saved tasks!", "err", err)
	}

	styles := newStyles()

	delegate := task.NewTaskDelegate()

	listModel := list.New(nil, delegate, 0, 0)
	listModel.Title = "Terminal Task"
	listModel.Styles.Title = styles.titleStyle
	listModel.SetShowStatusBar(true)
	listModel.SetStatusBarItemName("task", "tasks")

	listModel.SetItems(tasksToItems(tasks))
	return model{
		list:     listModel,
		editmenu: editmenu.New(0, 0),
		state:    stateList,
		keymap:   newListKeyMap(),
		styles:   styles,
		store:    store,
	}
}

// Convert Task to Item
func taskToItem(t task.Task) list.Item {
	return task.Task{TitleStr: t.TitleStr, DescStr: t.DescStr, DueDate: t.DueDate, Done: t.Done}
}

// Convert []task.Task to []list.Item
func tasksToItems(tasks []task.Task) []list.Item {
	if tasks == nil {
		return []list.Item{}
	}
	items := make([]list.Item, len(tasks))
	for i, t := range tasks {
		items[i] = taskToItem(t)
	}
	return items
}

// Convert Item to Task
func itemToTask(i list.Item) task.Task {
	t := i.(task.Task)
	return task.Task{TitleStr: t.TitleStr, DescStr: t.DescStr, DueDate: t.DueDate, Done: t.Done}
}

// Convert []list.Item to []task.Task
func itemsToTasks(items []list.Item) []task.Task {
	if items == nil {
		return []task.Task{}
	}
	tasks := make([]task.Task, len(items))
	for i, item := range items {
		tasks[i] = itemToTask(item)
	}
	return tasks
}
