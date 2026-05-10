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

func (m *model) updateListContent() {
	rows := make([]string, 0, len(m.notes))
	for i, n := range m.notes {
		prefix := " "
		if i == m.listIndex {
			prefix = ">"
		}

		rest := m.width - len(n.title)
		var shortBody string
		if rest > 0 {
			shortBody = strings.ReplaceAll(n.body, "\n", " ")
			if len(shortBody) > rest {
				shortBody = shortBody[:rest-3] + "…"
			}
		}
		rows = append(rows, listEnumeratorStyle.Render(prefix)+n.title+" | "+faint.Render(shortBody))
	}
	m.vp.SetContent(lipgloss.JoinVertical(lipgloss.Left, rows...))
}

func (m model) View() tea.View {
	var body string
	var help string
	switch m.state {
	case titleView:
		help = "enter: save , esc: discard"
		body = m.ti.View()
	case bodyView:
		help = "ctrl+s: save , esc: discard"
		body = lipgloss.JoinVertical(lipgloss.Left,
			m.currNote.title,
			m.ta.View(),
		)
	case listView:
		help = "n - new note • q - quit"
		body = m.vp.View()
	}

	m.vp.SetContent(body)
	v := tea.NewView(lipgloss.JoinVertical(lipgloss.Left,
		appNameStyle.Render("NOTES APP"),
		m.vp.View(),
		help,
	))
	v.AltScreen = true
	return v
}
