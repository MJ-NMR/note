package main

import (
	tea "charm.land/bubbletea/v2"
	lipgloss "charm.land/lipgloss/v2"
)

var (
	appNameStyle = lipgloss.NewStyle().Background(lipgloss.Color("99")).Padding(0, 1)

	faintStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("255")).Faint(true)

	// listEnumeratorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("99")).MarginRight(1)
)

func (m model) View() tea.View {
	header := lipgloss.PlaceHorizontal(m.width, lipgloss.Center, appNameStyle.Render("NOTES APP"))
	var body string
	var footer string

	switch m.state {
	case titleView:
		footer = "enter: save , esc: discard"
		body = m.ti.View()
	case bodyView:
		footer = "ctrl+s: save , esc: discard"
		body = lipgloss.JoinVertical(lipgloss.Left,
			"  "+appNameStyle.Render(m.currNote.title),
			m.ta.View(),
		)
	case listView:
		body = m.l.View()
	}

	footer = faintStyle.Render(footer)
	v := tea.NewView(lipgloss.JoinVertical(lipgloss.Left,
		header,
		body,
		footer,
	))
	v.AltScreen = true
	return v
}
