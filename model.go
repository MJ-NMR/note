package main

import (
	"log"

	"charm.land/bubbles/v2/textarea"
	"charm.land/bubbles/v2/textinput"
	"charm.land/bubbles/v2/viewport"
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
	vp        viewport.Model
	ta        textarea.Model
	ti        textinput.Model
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
	ti := textinput.New()
	ti.Prompt = "Note title:"

	return model{
		store: store,
		state: listView,
		vp:    viewport.New(),
		ta:    textarea.New(),
		ti:    ti,
		notes: notes,
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

	m.vp, cmd = m.vp.Update(msg)
	cmds = append(cmds, cmd)

	m.ta, cmd = m.ta.Update(msg)
	cmds = append(cmds, cmd)

	m.ti, cmd = m.ti.Update(msg)
	cmds = append(cmds, cmd)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		m.vp.SetWidth(msg.Width)
		m.vp.SetHeight(msg.Height - 2) // leave room for header/footer
		m.ta.SetWidth(msg.Width)
		m.ta.SetHeight(msg.Height)
		m.ti.SetWidth(msg.Width)
		m.updateListContent()

	case tea.KeyPressMsg:
		key := msg.String()
		switch m.state {
		case listView:
			switch key {
			case "q":
				return m, tea.Quit
			case "n":
				m.ti.SetValue("")
				m.ti.Focus()
				m.currNote = note{}
				m.state = titleView
			case "up", "k":
				if m.listIndex > 0 {
					m.listIndex--
					m.updateListContent()
				}
			case "down", "j":
				if m.listIndex < len(m.notes)-1 {
					m.listIndex++
					m.updateListContent()
				}
			case "enter":
				m.currNote = m.notes[m.listIndex]
				m.state = bodyView
				m.ta.SetValue(m.currNote.body)
				m.ta.Focus()
				m.ta.CursorEnd()
			}

		// Title Input View key bindings
		case titleView:
			switch key {
			case "enter":
				title := m.ti.Value()
				if title != "" {
					m.currNote.title = title

					m.state = bodyView
					m.ta.SetValue("")
					m.ta.Focus()
					m.ta.CursorEnd()
				}
			case "esc":
				m.state = listView
			}

		// Body Textarea key bindings
		case bodyView:
			switch key {
			case "ctrl+s":
				m.currNote.body = m.ta.Value()

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
