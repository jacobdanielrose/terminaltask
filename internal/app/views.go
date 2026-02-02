package app

// View renders the root application view based on the current state.
func (m model) View() string {
	switch m.state {
	case stateList:
		return m.styles.Frame.Render(m.list.View())
	case stateEdit:
		return m.styles.Frame.Render(m.editmenu.View())
	default:
		return "Unknown State"
	}
}
