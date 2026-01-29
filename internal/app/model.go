package app

import (
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

type statusMessageStyles struct {
	SuccessStyle lipgloss.Style
	ErrorStyle   lipgloss.Style
}

// ListStyles contains styles for the task list view.
type ListStyles struct {
	Title lipgloss.Style
}

// AppStyles is the top-level style graph for the application.
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

type model struct {
	list     list.Model
	editmenu editmenu.Model
	state    state
	keymap   *listKeyMap
	styles   AppStyles
	service  taskservice.Service
}

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

func configureListModel(listModel list.Model, styles ListStyles) list.Model {
	listModel.Title = listModelTitle
	listModel.Styles.Title = styles.Title
	listModel.SetShowStatusBar(true)
	listModel.SetStatusBarItemName("task", "tasks")
	return listModel
}

// Convert Task to Item
func taskToItem(t task.Task) list.Item {
	return t
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
	return i.(task.Task)
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
