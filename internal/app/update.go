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

// Extract magic strings to constants
const (
	// Generic error text for failed saves
	statusMsgSaveError = "Error saving!"

	// Success status templace
	statusMsgEditedTask    = "Edited: \"%s\""
	statusMsgDeletedTask   = "Deleted: \"%s\""
	statusMsgCompletedTask = "Completed: \"%s\""
	statusMsgCreatedTask   = "Created new task: \"%s\""
)

func (m model) Init() tea.Cmd {
	return m.loadTasksCmd()
}

// renderSuccessStatus formats a success status message with global styles.
func (m model) renderSuccessStatus(msg string) string {
	return m.styles.Status.SuccessStyle.Render(msg)
}

// renderErrorStatus formats an error status message with global styles.
func (m model) renderErrorStatus(msg string) string {
	return m.styles.Status.ErrorStyle.Render(msg)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Global key handling (applies in all states)
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		if key.Matches(keyMsg, m.keymap.Quit) {
			// ctrl+c: always quit the app, no matter where we are
			return m, tea.Quit
		}
	}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h, v := m.styles.Frame.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
		m.editmenu = m.editmenu.SetSize(msg.Width-h, msg.Height-v)
		return m, nil

	case editmenu.EscapeEditMsg:
		m.state = stateList
		return m, nil

	case task.EnterEditMsg:
		item := m.list.SelectedItem()
		if item == nil {
			return m, nil
		}
		t := item.(task.Task)
		m.editmenu = editmenu.New(t)
		m.state = stateEdit
		return m, nil

	case editmenu.SaveTaskMsg:
		return m.saveTask(msg)

	case task.ToggleDoneMsg:
		item := m.list.SelectedItem()
		if item == nil {
			return m, nil
		}
		task := item.(task.Task)
		index := m.list.Index()
		task.Done = !task.Done
		m.list.SetItem(index, task)
		tasks := itemsToTasks(m.list.Items())

		var statusText string
		if task.Done {
			statusText = fmt.Sprintf(statusMsgCompletedTask, task.Title())
		} else {
			statusText = fmt.Sprintf(statusMsgEditedTask, task.Title())
		}
		return m, m.saveTasksCmd(tasks, statusText)

	case TasksSavedMsg:
		cmd := m.list.NewStatusMessage(
			m.renderSuccessStatus(msg.msg),
		)
		return m, cmd

	case TasksSaveErrorMsg:
		cmd := m.list.NewStatusMessage(
			m.renderErrorStatus(statusMsgSaveError),
		)
		log.Error("Error saving tasks", "err", msg.Err, "store", m.service.Name())
		return m, cmd

	case TasksLoadedMsg:
		m.list.SetItems(tasksToItems(msg.Tasks))
		return m, nil

	case task.DeleteMsg:
		index := m.list.Index()
		tasks := itemsToTasks(m.list.Items())
		deletedTask := tasks[index]
		m.list.RemoveItem(index)
		statusText := fmt.Sprintf(statusMsgDeletedTask, deletedTask.Title())
		return m, m.saveTasksCmd(tasks, statusText)
	}

	switch m.state {
	case stateList:
		return m.stateListUpdate(msg)
	case stateEdit:
		return m.stateEditUpdate(msg)
	}
	return m, nil
}

func (m model) saveTask(msg editmenu.SaveTaskMsg) (model, tea.Cmd) {
	t := task.Task{
		TitleStr: msg.Title,
		DescStr:  msg.Desc,
		DueDate:  msg.Date,
		Done:     msg.Done,
	}
	t.SetID(msg.TaskID)

	index := -1
	var statusText string

	if len(m.list.Items()) != 0 && !msg.IsNew {
		index = m.list.Index()
		m.list.SetItem(index, t)
		statusText = fmt.Sprintf(statusMsgEditedTask, t.Title())
	} else {
		m.list.InsertItem(index, t)
		statusText = fmt.Sprintf(statusMsgCreatedTask, t.Title())
	}

	m.state = stateList
	tasks := itemsToTasks(m.list.Items())

	return m, m.saveTasksCmd(tasks, statusText)
}

func (m model) stateListUpdate(msg tea.Msg) (model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Disables other keys if actively filtering
		if m.list.FilterState() == list.Filtering {
			break
		}
		if key.Matches(msg, m.keymap.NewItem) {
			newTask := task.New()
			m.editmenu = editmenu.New(newTask)
			m.state = stateEdit
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m model) stateEditUpdate(msg tea.Msg) (model, tea.Cmd) {
	var cmd tea.Cmd
	m.editmenu, cmd = m.editmenu.Update(msg)
	return m, cmd
}
