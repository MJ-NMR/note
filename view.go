package main

import (
	"strings"

	tea "charm.land/bubbletea/v2"
	lipgloss "charm.land/lipgloss/v2"
)

var (
	appNameStyle = lipgloss.NewStyle().Background(lipgloss.Color("99")).Padding(0, 1)

	faint = lipgloss.NewStyle().Foreground(lipgloss.Color("255")).Faint(true)

	listEnumeratorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("99")).MarginRight(1)
)

func (m model) View() tea.View {
	var content string
	switch m.state {
	case titleView:
		content = lipgloss.JoinVertical(lipgloss.Left,
			"Note title:",
			m.textinput.View(),
			faint.Render("enter - save • esc - discard"),
		)
	case bodyView:
		content = lipgloss.JoinVertical(lipgloss.Left,
			"Note:",
			m.textarea.View(),
			faint.Render("ctrl+s - save • esc - discard"),
		)
	case listView:
		rows := make([]string, 0, len(m.notes))
		for i, n := range m.notes {
			prefix := " "
			if i == m.listIndex {
				prefix = ">"
			}
			shortBody := strings.ReplaceAll(n.body, "\n", " ")
			if len(shortBody) > 30 {
				shortBody = shortBody[:30] + "…"
			}
			rows = append(rows, listEnumeratorStyle.Render(prefix)+n.title+" | "+faint.Render(shortBody))
		}
		rows = append(rows, faint.Render("n - new note • q - quit"))
		content = lipgloss.JoinVertical(lipgloss.Left, rows...)
	}

	v := tea.NewView(lipgloss.JoinVertical(lipgloss.Left,
		appNameStyle.Render("NOTES APP"),
		content,
	))

	v.AltScreen = true
	return v
}
