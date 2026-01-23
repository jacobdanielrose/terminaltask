package main

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
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
	task := m.list.SelectedItem().(task.Task)
	index := m.list.Index()
	task.TitleStr = msg.Title
	task.DescStr = msg.Desc
	task.DueDate = msg.Date
	m.list.SetItem(index, task)
	m.state = stateList
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
	m.editmenu.Title.SetValue(task.Title())
	m.editmenu.Title.SetCursor(len(task.TitleStr))
	m.editmenu.Desc.SetValue(task.Description())
	m.editmenu.Desc.SetCursor(len(task.DescStr))
	m.editmenu.DatePicker.SetTime(task.DueDate)
	return m
}

func (m model) newTask() model {
	m.state = stateEdit
	m.editmenu.Title.SetValue("New Task")
	m.editmenu.Desc.SetValue("Description")
	return m
}
