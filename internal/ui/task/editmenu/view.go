package editmenu

import (
	"github.com/charmbracelet/lipgloss"
)

func (m Model) View() string {

	var (
		sections    []string
		availHeight = m.height
	)

	if m.showTitle {
		v := m.titleView()
		sections = append(sections, v)
		availHeight -= lipgloss.Height(v)
	}

	var help string
	if m.showHelp {
		help = m.helpView()
		availHeight -= lipgloss.Height(help)
	}

	editContent := lipgloss.NewStyle().Height(availHeight).Render(m.editView())
	sections = append(sections, editContent)

	if m.showHelp {
		sections = append(sections, help)
	}

	return lipgloss.JoinVertical(lipgloss.Left, sections...)
}

func (m Model) titleView() string {
	var (
		view          string
		titleBarStyle = m.styles.TitleBar
	)

	view += m.styles.Title.Render(m.Title)

	if len(view) > 0 {
		return titleBarStyle.Render(view)
	}

	return view
}

func (m Model) ShowTitle() bool {
	return m.showTitle
}

func (m Model) editView() string {
	m.TaskTitle.TextStyle = m.styles.Normal
	m.TaskTitle.PromptStyle = m.styles.Normal
	m.Desc.TextStyle = m.styles.Normal
	m.Desc.PromptStyle = m.styles.Normal

	lipgloss.JoinVertical(0.5)
	switch m.focusIdx {
	case 0:
		m.TaskTitle.TextStyle = m.styles.Normal
		m.TaskTitle.PromptStyle = m.styles.Focused
	case 1:
		m.Desc.TextStyle = m.styles.Normal
		m.Desc.PromptStyle = m.styles.Focused
		//	case 2:
		//	m.datePicker.Time = m.Duedate
	}

	titleDescCombo := lipgloss.JoinVertical(lipgloss.Left,
		m.TaskTitle.View(),
		m.Desc.View(),
	)

	return lipgloss.JoinVertical(0.5,
		titleDescCombo,
		m.DatePicker.View(),
	)

}

// SetShowHelp shows or hides the help view.
func (m *Model) SetShowHelp(v bool) {
	m.showHelp = v
}

// ShowHelp returns whether or not the help is set to be rendered.
func (m Model) ShowHelp() bool {
	return m.showHelp
}

func (m Model) helpView() string {
	return m.styles.HelpStyle.Render(m.help.View(m.keymap))
}
