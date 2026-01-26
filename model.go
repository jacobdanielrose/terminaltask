package main

import (
	"log"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
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
	store    TaskStore
}

func initialModel(store TaskStore) model {

	tasks, err := store.Load()
	if err != nil {
		tasks = []task.Task{}
		log.Fatal(err)
	}

	// tasks := []list.Item{
	// 	task.Task{TitleStr: "Buy groceries", DescStr: "Milk, eggs, bread", DueDate: time.Now().Add(time.Hour * 24), Done: false},
	// 	task.Task{TitleStr: "Complete project report", DescStr: "Finish the final draft of the project report before the meeting.", DueDate: time.Now().Add(time.Hour * 48), Done: false},
	// 	task.Task{TitleStr: "Plan weekend getaway", DescStr: "Research destination options and make a booking.", DueDate: time.Now().AddDate(0, 0, 5), Done: true},
	// 	task.Task{TitleStr: "Book dentist appointment", DescStr: "Schedule a routine check-up and cleaning.", DueDate: time.Now().AddDate(0, 1, 0), Done: false},
	// 	task.Task{TitleStr: "Exercise", DescStr: "Morning workout: jogging and stretching.", DueDate: time.Now().Add(time.Hour * 2), Done: true},
	// 	task.Task{TitleStr: "Call mom", DescStr: "Check in and see how she's doing.", DueDate: time.Now().Add(time.Hour * 24 * 2), Done: false},
	// 	task.Task{TitleStr: "Read a book", DescStr: "Finish reading the current novel.", DueDate: time.Now().AddDate(0, 0, 7), Done: false},
	// 	task.Task{TitleStr: "Pay utility bills", DescStr: "Pay electricity, water, and gas bills online.", DueDate: time.Now().Add(time.Hour * 72), Done: true},
	// }

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
