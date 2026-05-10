package main

import (
	"log"

	"charm.land/bubbles/v2/textarea"
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
)

const (
	listView uint = iota
	titleView
	bodyView
)

type model struct {
	store     *Store
	state     uint
	textarea  textarea.Model
	textinput textinput.Model
	currNote  note
	notes     []note
	listIndex int
	width     int
	height    int
}

func NewModel(store *Store) model {
	notes, err := store.GetNotes()
	if err != nil {
		log.Fatalf("unable to get notes: %v", err)
	}

	return model{
		store:     store,
		state:     listView,
		textarea:  textarea.New(),
		textinput: textinput.New(),
		notes:     notes,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmds []tea.Cmd
		cmd  tea.Cmd
	)

	m.textarea, cmd = m.textarea.Update(msg)
	cmds = append(cmds, cmd)

	m.textinput, cmd = m.textinput.Update(msg)
	cmds = append(cmds, cmd)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		m.textarea.SetWidth(msg.Width)
		m.textarea.SetHeight(max(msg.Height-4, 0)) // leave room for header/footer
		m.textinput.SetWidth(max(msg.Width-4, 0))

	case tea.KeyPressMsg:
		key := msg.String()
		switch m.state {
		case listView:
			switch key {
			case "q":
				return m, tea.Quit
			case "n":
				m.textinput.SetValue("")
				m.textinput.Focus()
				m.currNote = note{}
				m.state = titleView
			case "up", "k":
				if m.listIndex > 0 {
					m.listIndex--
				}
			case "down", "j":
				if m.listIndex < len(m.notes)-1 {
					m.listIndex++
				}
			case "enter":
				m.currNote = m.notes[m.listIndex]
				m.state = bodyView
				m.textarea.SetValue(m.currNote.body)
				m.textarea.Focus()
				m.textarea.CursorEnd()
			}

		// Title Input View key bindings
		case titleView:
			switch key {
			case "enter":
				title := m.textinput.Value()
				if title != "" {
					m.currNote.title = title

					m.state = bodyView
					m.textarea.SetValue("")
					m.textarea.Focus()
					m.textarea.CursorEnd()
				}
			case "esc":
				m.state = listView
			}

		// Body Textarea key bindings
		case bodyView:
			switch key {
			case "ctrl+s":
				m.currNote.body = m.textarea.Value()

				var err error
				if err = m.store.SaveNote(m.currNote); err != nil {
					// TODO: handle error instead of quitting
					return m, tea.Quit
				}

				m.notes, err = m.store.GetNotes()
				if err != nil {
					// TODO: handle error instead of quitting
					return m, tea.Quit
				}

				m.state = listView
			case "esc":
				m.state = listView
			}
		}
	}

	return m, tea.Batch(cmds...)
}
