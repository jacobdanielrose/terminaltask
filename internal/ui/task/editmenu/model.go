package editmenu

import (
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	datepicker "github.com/ethanefung/bubble-datepicker"
)

type EscapeEditMsg struct{}
type SaveTaskMsg struct {
	Title string
	Desc  string
	Date  time.Time
}

type Model struct {
	Title      textinput.Model
	Desc       textinput.Model
	focusIdx   int // 0=title, 1=desc, 2=duedate
	keymap     *EditTaskKeyMap
	isNew      bool
	styles     Styles
	help       help.Model
	width      int
	height     int
	DatePicker datepicker.Model
}

func New() Model {
	title := textinput.New()
	title.Prompt = "Title: "
	title.Placeholder = "Title"
	title.Focus()

	desc := textinput.New()
	desc.Prompt = "Description: "
	desc.Placeholder = "Description"

	return Model{
		Title:      title,
		Desc:       desc,
		focusIdx:   0,
		keymap:     newEditTaskKeyMap(),
		styles:     NewStyles(),
		help:       help.New(),
		width:      10,
		height:     10,
		DatePicker: datepicker.New(time.Now()),
	}
}

func (m *Model) SetSize(width int, height int) {
	m.width = width
	m.help.Width = width

	m.height = height
}

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	// Is it a key press?
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keymap.SaveField):
			m.focusIdx = (m.focusIdx + 1) % 3
			m.setFocus()
		case key.Matches(msg, m.keymap.SaveTask):
			return m, func() tea.Msg {
				return SaveTaskMsg{
					m.Title.Value(),
					m.Desc.Value(),
					m.DatePicker.Time,
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
	m.Title, cmd = m.Title.Update(msg)
	cmds = append(cmds, cmd)
	m.Desc, cmd = m.Desc.Update(msg)
	cmds = append(cmds, cmd)
	m.DatePicker, cmd = m.DatePicker.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	m.Title.TextStyle = m.styles.Normal
	m.Title.PromptStyle = m.styles.Normal
	m.Desc.TextStyle = m.styles.Normal
	m.Desc.PromptStyle = m.styles.Normal

	switch m.focusIdx {
	case 0:
		m.Title.TextStyle = m.styles.Normal
		m.Title.PromptStyle = m.styles.Focused
	case 1:
		m.Desc.TextStyle = m.styles.Normal
		m.Desc.PromptStyle = m.styles.Focused
		//	case 2:
		//	m.datePicker.Time = m.Duedate
	}

	mainView := m.Title.View() + "\n\n" +
		m.Desc.View() + "\n\n" +
		m.DatePicker.View()

	helpView := m.help.View(m.keymap)
	height := 8 - strings.Count(helpView, "\n")
	if mainView != "" {
		mainView += "\n"
	}

	return mainView + strings.Repeat("\n", height) + helpView
}

func (m *Model) setFocus() {
	m.Title.Blur()
	m.Desc.Blur()
	m.DatePicker.Blur()
	switch m.focusIdx {
	case 0:
		m.Title.Focus()
	case 1:
		m.Desc.Focus()
	case 2:
		m.DatePicker.SelectDate()
		m.DatePicker.SetFocus(datepicker.FocusHeaderMonth)
	}
}
