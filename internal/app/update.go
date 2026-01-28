package app

import (
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jacobdanielrose/terminaltask/internal/task"
	"github.com/jacobdanielrose/terminaltask/internal/task/editmenu"
)

var statusMessageStyle = lipgloss.NewStyle().
	Foreground(lipgloss.AdaptiveColor{Light: "#04B575", Dark: "#04B575"}).
	Render

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h, v := m.styles.appStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
		m.editmenu.SetSize(msg.Width-h, msg.Height-v)
	case editmenu.EscapeEditMsg:
		m.state = stateList
		return m, nil
	case task.EnterEditMsg:
		task := m.list.SelectedItem().(task.Task)
		m.editmenu = editmenu.New(task)
		m.state = stateEdit
		return m, nil
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

func (m Model) saveTask(msg editmenu.SaveTaskMsg) Model {
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

	m.list.NewStatusMessage(
		statusMessageStyle(
			fmt.Sprintf("Edited: \"%s\"", task.Title()),
		),
	)
	return m
}

func (m Model) stateListUpdate(msg tea.Msg) (Model, tea.Cmd) {
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

func (m Model) stateEditUpdate(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	m.editmenu, cmd = m.editmenu.Update(msg)
	return m, cmd
}
