// Package app contains the primary Bubble Tea application model for
// terminaltask. It wires together the list view, edit menu, storage
// service, and configuration to provide the main interactive program.
package app

import (
	"time"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jacobdanielrose/terminaltask/internal/config"
	taskservice "github.com/jacobdanielrose/terminaltask/internal/service"
	task "github.com/jacobdanielrose/terminaltask/internal/task"
	"github.com/jacobdanielrose/terminaltask/internal/task/editmenu"
)

type state int

const (
	stateList state = iota
	stateEdit
)

const (
	listModelTitle = "Terminal Task"
)

// statusMessageStyles groups the global styles used for rendering
// success and error status messages throughout the application.
type statusMessageStyles struct {
	SuccessStyle lipgloss.Style
	ErrorStyle   lipgloss.Style
}

// ListStyles contains styles for the task list view, including the
// list title styling applied at the top of the list component.
type ListStyles struct {
	Title lipgloss.Style
}

// AppStyles is the top-level style graph for the application. It
// centralizes all styling concerns so that the model, update, and view
// logic can share a consistent visual configuration.
type AppStyles struct {
	// Frame is the global application frame (padding/margins) applied
	// around both the list view and the editmenu view.
	Frame lipgloss.Style

	// Status contains global success/error status styles shared by
	// all views (list, edit menu, etc.).
	Status statusMessageStyles

	// List contains styles specific to the list component.
	List ListStyles

	// EditMenu contains styles for the edit menu container.
	EditMenu editmenu.Styles

	// Form contains styles for the inner edit form.
	Form editmenu.Styles
}

// newAppStyles constructs the top-level styles for the app.
func newAppStyles() AppStyles {
	listTitle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFFDF5")).
		Background(lipgloss.Color("#25A065")).
		Padding(0, 1)

	status := statusMessageStyles{
		SuccessStyle: lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "#04B575", Dark: "#04B575"}),
		ErrorStyle: lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "#FF0000", Dark: "#FF0000"}),
	}

	editMenuStyles := editmenu.DefaultStyles()
	formStyles := editmenu.DefaultStyles()

	// Use global status styles for the edit menu status message.
	editMenuStyles.StatusMessage = status.SuccessStyle

	return AppStyles{
		Frame:  lipgloss.NewStyle().Padding(1, 2),
		Status: status,
		List: ListStyles{
			Title: listTitle,
		},
		EditMenu: editMenuStyles,
		Form:     formStyles,
	}
}

// model is the root Bubble Tea model for terminaltask. It coordinates
// the task list, edit menu, current application state, styles, and
// backing task service.
type model struct {
	// list is the main task list view.
	list list.Model

	// editmenu is the sub-model responsible for editing/creating tasks.
	editmenu editmenu.Model

	// state tracks whether the user is interacting with the list or
	// the edit menu.
	state state

	// keymap holds global key bindings for list-level interactions.
	keymap *listKeyMap

	// styles contains all top-level styling information for the app.
	styles AppStyles

	// service abstracts persistence and higher-level task operations.
	service taskservice.Service
}

// NewModel constructs a new application model wired with the provided
// configuration and task service. It initializes the list and edit
// menu sub-models and returns a Bubble Tea model ready for use in a
// tea.Program.
func NewModel(cfg config.Config, service taskservice.Service) tea.Model {
	appStyles := newAppStyles()

	delegate := task.NewTaskDelegate()
	listModel := list.New(nil, delegate, 0, 0)
	editmenuModel := editmenu.NewWithStyles(task.Task{}, appStyles.EditMenu, appStyles.Form)

	listModel = configureListModel(listModel, appStyles.List)

	return model{
		list:     listModel,
		editmenu: editmenuModel,
		state:    stateList,
		keymap:   NewListKeyMap(),
		styles:   appStyles,
		service:  service,
	}
}

// configureListModel applies application-specific configuration and
// styling to the given list model
func configureListModel(listModel list.Model, styles ListStyles) list.Model {
	listModel.Title = listModelTitle
	listModel.Styles.Title = styles.Title
	listModel.SetShowStatusBar(true)
	listModel.SetStatusBarItemName("task", "tasks")
	listModel.StatusMessageLifetime = 1 * time.Second
	return listModel
}

// taskToItem converts a Task into a list.Item.
func taskToItem(t task.Task) list.Item {
	return t
}

// tasksToItems converts a slice of Task values into a slice of
// list.Item so they can be displayed by the list component. It always
// returns a non-nil slice.
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

// itemToTask converts a list.Item back into a Task.
func itemToTask(i list.Item) task.Task {
	return i.(task.Task)
}

// itemsToTasks converts a slice of list.Item values back into a slice
// of Task. It always returns a non-nil slice.
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
