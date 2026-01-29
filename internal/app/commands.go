package app

import (
	tea "github.com/charmbracelet/bubbletea"
	task "github.com/jacobdanielrose/terminaltask/internal/task"
)

func (m model) saveTasksCmd(tasks []task.Task, msg string) tea.Cmd {
	return func() tea.Msg {
		err := m.service.SaveTasks(tasks)
		if err != nil {
			return TasksSaveErrorMsg{Err: err}
		}
		return TasksSavedMsg{msg: msg}
	}
}

func (m model) loadTasksCmd() tea.Cmd {
	return func() tea.Msg {
		tasks, err := m.service.LoadTasks()
		if err != nil {
			return TasksLoadErrorMsg{Err: err}
		}
		return TasksLoadedMsg{Tasks: tasks}
	}
}
