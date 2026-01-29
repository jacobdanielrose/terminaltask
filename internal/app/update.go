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
	errSavingTasks = "Error saving!"
	msgEditedTask  = "Edited: \"%s\""
	msgDeletedTask = "Deleted: \"%s\""
)

func (m model) Init() tea.Cmd {
	return m.loadTasksCmd()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
		return m, m.saveTasksCmd(tasks, fmt.Sprintf(msgEditedTask, task.Title()))
	case TasksSavedMsg:
		m.list.NewStatusMessage(
			m.styles.
				Status.
				SuccessStyle.
				Render(msg.msg),
		)
		return m, nil
	case TasksSaveErrorMsg:
		m.list.NewStatusMessage(errSavingTasks)
		log.Error("Error saving tasks", "err", msg.Err, "store", m.service.Name())
		return m, nil
	case TasksLoadedMsg:
		m.list.SetItems(tasksToItems(msg.Tasks))
		return m, nil

	case task.DeleteMsg:
		index := m.list.Index()
		tasks := itemsToTasks(m.list.Items())
		deletedTask := tasks[index]
		m.list.RemoveItem(index)
		return m, m.saveTasksCmd(tasks, fmt.Sprintf(msgDeletedTask, deletedTask.Title()))
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
	if len(m.list.Items()) != 0 && !msg.IsNew {
		index = m.list.Index()
		m.list.SetItem(index, t)
	} else {
		m.list.InsertItem(index, t)
	}

	m.state = stateList
	tasks := itemsToTasks(m.list.Items())

	return m, m.saveTasksCmd(tasks, fmt.Sprintf(msgEditedTask, t.Title()))
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
