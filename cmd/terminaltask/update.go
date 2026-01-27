package main

import (
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/google/uuid"
	"github.com/jacobdanielrose/terminaltask/internal/ui/task"
	"github.com/jacobdanielrose/terminaltask/internal/ui/task/editmenu"
)

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h, v := m.styles.appStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
		m.editmenu.SetSize(msg.Width-h, msg.Height-v)
	case editmenu.EscapeEditMsg:
		m.state = stateList
		return m, nil
	case task.EnterEditMsg:
		return m.editTask(), nil
	case editmenu.SaveTaskMsg:
		return m.saveTask(msg), nil
	case task.ToggleDoneMsg:
		task := m.list.SelectedItem().(task.Task)
		index := m.list.Index()
		m.list.SetItem(index, task)
		tasks := itemsToTasks(m.list.Items())
		_ = m.store.Save(tasks)
		return m, nil
	case task.DeleteMsg:
		// deletion in list already happens in the delegate
		// just need to save to backend here.
		tasks := itemsToTasks(m.list.Items())
		_ = m.store.Save(tasks)
		return m, nil
	}

	switch m.state {
	case stateList:
		return m.stateListUpdate(msg)
	case stateEdit:
		return m.stateEditUpdate(msg)
	}
	return m, nil
}

func (m model) saveTask(msg editmenu.SaveTaskMsg) model {
	task := task.Task{
		TitleStr: msg.Title,
		DescStr:  msg.Desc,
		DueDate:  msg.Date,
		Done:     msg.Done,
	}
	task.SetID(msg.TaskID)

	index := -1
	if len(m.list.Items()) != 0 && !msg.IsNew {
		index = m.list.Index()
		m.list.SetItem(index, task)
	} else {
		m.list.InsertItem(index, task)
	}

	m.state = stateList
	tasks := itemsToTasks(m.list.Items())
	_ = m.store.Save(tasks)
	return m
}

func (m model) stateListUpdate(msg tea.Msg) (model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Disables other keys if actively filtering
		if m.list.FilterState() == list.Filtering {
			break
		}
		if key.Matches(msg, m.keymap.NewItem) {
			m = m.newTask()
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

func (m model) editTask() model {
	m.state = stateEdit
	task := m.list.SelectedItem().(task.Task)
	m.editmenu.IsNew = false
	m.editmenu.TaskID = task.GetID()
	m.editmenu.TaskTitle.SetValue(task.Title())
	m.editmenu.TaskTitle.SetCursor(len(task.TitleStr))
	m.editmenu.Desc.SetValue(task.Description())
	m.editmenu.Desc.SetCursor(len(task.DescStr))
	m.editmenu.DatePicker.SetTime(task.DueDate)
	return m
}

func (m model) newTask() model {
	m.state = stateEdit
	taskID := uuid.New()
	m.editmenu.IsNew = true
	m.editmenu.TaskID = taskID
	m.editmenu.TaskTitle.SetValue("")
	m.editmenu.Desc.SetValue("")
	m.editmenu.TaskTitle.Placeholder = "New Task"
	m.editmenu.Desc.Placeholder = "Description"
	m.editmenu.TaskTitle.SetCursor(0)
	m.editmenu.Desc.SetCursor(0)
	m.editmenu.DatePicker.SetTime(time.Now())
	return m
}
