// Package app contains the primary Bubble Tea application update loop for
// terminaltask. It defines how the root model initializes, responds to
// incoming messages, and coordinates state transitions between the list
// and edit views.
package app

import (
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
	"github.com/jacobdanielrose/terminaltask/internal/task"
	"github.com/jacobdanielrose/terminaltask/internal/task/editmenu"
)

// Extract magic strings to constants.
const (
	// Generic error text for failed saves.
	statusMsgSaveError   = "Error saving!"
	statusMsgDeleteError = "Error deleting task!"

	// Success status templates.
	statusMsgEditedTask    = "Edited: \"%s\""
	statusMsgDeletedTask   = "Deleted: \"%s\""
	statusMsgCompletedTask = "Completed: \"%s\""
	statusMsgCreatedTask   = "Created new task: \"%s\""
)

// Init implements tea.Model and, in this application, triggers loading
// tasks from the backing service.
func (m model) Init() tea.Cmd {
	return m.loadTasksCmd()
}

// renderSuccessStatus formats a success status message using the
// application-wide success style. It is used when operations like
// saving or editing tasks complete successfully.
func (m model) renderSuccessStatus(msg string) string {
	return m.styles.Status.SuccessStyle.Render(msg)
}

// renderErrorStatus formats an error status message using the
// application-wide error style. It is used when persistence or other
// operations fail and a message should be shown in the status bar.
func (m model) renderErrorStatus(msg string) string {
	return m.styles.Status.ErrorStyle.Render(msg)
}

// Update implements the Bubble Tea Update method for the root
// application model. It delegates high-level messages first, then
// routes to the appropriate state-specific update function.
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Global key handling (applies in all states).
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		if key.Matches(keyMsg, m.keymap.Quit) {
			// Ctrl+C: always quit the app, no matter where we are.
			return m, tea.Quit
		}
	}

	// High-level message routing.
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		return m.resizeWindow(msg)

	case editmenu.EscapeEditMsg:
		m.state = stateList
		return m, nil

	case task.EnterEditMsg:
		return m.enterEditMenu()

	case editmenu.SaveTaskMsg:
		return m.saveTask(msg)

	case task.ToggleDoneMsg:
		return m.toggleDone()

	case TasksSavedMsg:
		return m.taskSaved(msg)

	case TasksSaveErrorMsg:
		return m.taskSaveError(msg)

	case TasksLoadedMsg:
		m.list.SetItems(tasksToItems(msg.Tasks))
		return m, nil

	case task.DeleteMsg:
		return m.deleteTask()
	}

	// Fallback to state-specific handling.
	switch m.state {
	case stateList:
		return m.stateListUpdate(msg)
	case stateEdit:
		return m.stateEditUpdate(msg)
	default:
		return m, nil
	}
}

func (m model) resizeWindow(msg tea.WindowSizeMsg) (tea.Model, tea.Cmd) {
	h, v := m.styles.Frame.GetFrameSize()
	contentW, contentH := msg.Width-h, msg.Height-v
	m.list.SetSize(contentW, contentH)
	m.editmenu = m.editmenu.SetSize(contentW, contentH)
	return m, nil
}

func (m model) enterEditMenu() (tea.Model, tea.Cmd) {
	item := m.list.SelectedItem()
	if item == nil {
		return m, nil
	}
	t := item.(task.Task)
	w, h := m.editmenu.Width(), m.editmenu.Height()
	m.editmenu = editmenu.NewWithSize(w, h, t)
	m.state = stateEdit
	return m, nil
}

func (m model) toggleDone() (tea.Model, tea.Cmd) {
	item := m.list.SelectedItem()
	if item == nil {
		return m, nil
	}

	taskItem := item.(task.Task)
	index := m.list.Index()
	taskItem.Done = !taskItem.Done
	m.list.SetItem(index, taskItem)

	tasks := itemsToTasks(m.list.Items())

	var statusText string
	if taskItem.Done {
		statusText = fmt.Sprintf(statusMsgCompletedTask, taskItem.Title())
	} else {
		statusText = fmt.Sprintf(statusMsgEditedTask, taskItem.Title())
	}

	return m, m.saveTasksCmd(tasks, statusText)
}

func (m model) taskSaveError(msg TasksSaveErrorMsg) (tea.Model, tea.Cmd) {
	cmd := m.list.NewStatusMessage(
		m.renderErrorStatus(statusMsgSaveError),
	)
	log.Error("Error saving tasks", "err", msg.Err, "store", m.service.Name())
	return m, cmd
}

func (m model) taskSaved(msg TasksSavedMsg) (tea.Model, tea.Cmd) {
	cmd := m.list.NewStatusMessage(
		m.renderSuccessStatus(msg.msg),
	)
	return m, cmd
}

func (m model) deleteTask() (tea.Model, tea.Cmd) {
	index := m.list.Index()
	taskItem, ok := m.list.SelectedItem().(task.Task)
	if !ok {
		return m, nil
	}

	if err := m.service.DeleteByID(taskItem.GetID()); err != nil {
		cmd := m.list.NewStatusMessage(
			m.renderErrorStatus(statusMsgDeleteError),
		)
		log.Error("Error deleting task", "err", err, "store", m.service.Name())
		return m, cmd
	}

	m.list.RemoveItem(index)
	cmd := m.list.NewStatusMessage(
		m.renderSuccessStatus(fmt.Sprintf(statusMsgDeletedTask, taskItem.Title())),
	)
	return m, cmd
}

// saveTask handles an editmenu.SaveTaskMsg by either updating an
// existing task in the list or inserting a new one. It then switches
// back to the list state and returns a command that persists the
// updated tasks and shows an appropriate status message.
func (m model) saveTask(msg editmenu.SaveTaskMsg) (model, tea.Cmd) {
	t := task.NewWithOptions(
		msg.Title,
		msg.Desc,
		msg.Date,
		msg.Done,
	)

	index := -1
	var statusText string

	if len(m.list.Items()) != 0 && !msg.IsNew {
		// Existing task, preserve ID.
		t.SetID(msg.TaskID)

		index = m.list.Index()
		m.list.SetItem(index, t)
		statusText = fmt.Sprintf(statusMsgEditedTask, t.Title())
	} else {
		// New task: use the ID generated by NewWithOptions.
		m.list.InsertItem(index, t)
		statusText = fmt.Sprintf(statusMsgCreatedTask, t.Title())
	}

	m.state = stateList
	tasks := itemsToTasks(m.list.Items())

	if err := m.service.SaveTasks(tasks); err != nil {
		cmd := m.list.NewStatusMessage(
			m.renderErrorStatus(statusMsgSaveError),
		)
		log.Error("Error saving tasks", "err", err, "store", m.service.Name())
		return m, cmd
	}

	cmd := m.list.NewStatusMessage(
		m.renderSuccessStatus(statusText),
	)
	return m, cmd
}

// stateListUpdate handles messages that should be processed while the
// application is in the list state. It is responsible for responding
// to list-specific keybindings (such as creating a new task) and for
// delegating messages down to the list sub-model.
func (m model) stateListUpdate(msg tea.Msg) (model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Disable other keys if actively filtering.
		if m.list.FilterState() == list.Filtering {
			break
		}
		// New item: open the edit menu with an empty task.
		if key.Matches(msg, m.keymap.NewItem) {
			newTask := task.New()
			w, h := m.editmenu.Width(), m.editmenu.Height()
			m.editmenu = editmenu.NewWithSize(w, h, newTask)
			m.state = stateEdit
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

// stateEditUpdate handles messages that should be processed while the
// application is in the edit state. It delegates the message to the
// edit menu sub-model and returns the updated root model and command.
func (m model) stateEditUpdate(msg tea.Msg) (model, tea.Cmd) {
	var cmd tea.Cmd
	m.editmenu, cmd = m.editmenu.Update(msg)
	return m, cmd
}
