package app

func (m model) View() string {
	switch m.state {
	case stateList:
		return m.styles.appStyle.Render(m.list.View())
	case stateEdit:
		return m.styles.appStyle.Render(m.editmenu.View())
	default:
		return "Unknown State"
	}
}
