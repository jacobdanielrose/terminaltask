package editmenu

import (
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	datepicker "github.com/ethanefung/bubble-datepicker"
)

type clearStatusMsg struct{}

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

func (m *Model) showStatus(msg string) tea.Cmd {
	m.statusMsg = msg
	return tea.Tick(1*time.Second, func(t time.Time) tea.Msg {
		return clearStatusMsg{}
	})
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case clearStatusMsg:
		m.statusMsg = ""
		return m, nil
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keymap.SaveField):
			m.focusIdx = (m.focusIdx + 1) % 3
			m.setFocus()
		case key.Matches(msg, m.keymap.SaveTask):
			if m.DatePicker.Time.Before(time.Now().Truncate(24 * time.Hour)) {
				return m, m.showStatus("You cannot pick a date in the past!")
			}
			m.focusIdx = 0
			m.setFocus()
			return m, func() tea.Msg {
				return SaveTaskMsg{
					m.TaskID,
					m.TaskTitle.Value(),
					m.Desc.Value(),
					m.DatePicker.Time,
					false,
					m.IsNew,
				}
			}
		case key.Matches(msg, m.keymap.EscapeEditMode):
			return m, func() tea.Msg {
				return EscapeEditMsg{}
			}
		case key.Matches(msg, m.keymap.Help):
			m.help.ShowAll = !m.help.ShowAll
		case key.Matches(msg, m.keymap.Quit):
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.TaskTitle, cmd = m.TaskTitle.Update(msg)
	cmds = append(cmds, cmd)
	m.Desc, cmd = m.Desc.Update(msg)
	cmds = append(cmds, cmd)
	m.DatePicker, cmd = m.DatePicker.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m *Model) setFocus() {
	m.TaskTitle.Blur()
	m.Desc.Blur()
	m.DatePicker.Blur()
	switch m.focusIdx {
	case 0:
		m.TaskTitle.Focus()
	case 1:
		m.Desc.Focus()
	case 2:
		m.DatePicker.SelectDate()
		m.DatePicker.SetFocus(datepicker.FocusCalendar)
	}
}

func (m *Model) SetSize(width int, height int) {
	m.width = width
	m.help.Width = width

	m.height = height
}
