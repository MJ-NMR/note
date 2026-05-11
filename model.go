package main

import (
	"log"

	"charm.land/bubbles/v2/list"
	"charm.land/bubbles/v2/textarea"
	"charm.land/bubbles/v2/textinput"

	// "charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
)

const (
	listView uint = iota
	titleView
	bodyView
)

type note struct {
	id    int64
	title string
	body  string
}

func (n note) Title() string { return n.title }
func (n note) Description() string {
	return n.body
}
func (n note) FilterValue() string { return n.title }

func NewList(notes []note) list.Model {
	items := make([]list.Item, len(notes))
	for i := range notes {
		items[i] = notes[i]
	}

	l := list.New(items, list.NewDefaultDelegate(), 0, 0)
	l.SetShowFilter(false)
	l.SetFilteringEnabled(false)
	return l
}

type model struct {
	store *Store
	state uint
	// vp        viewport.Model
	l        list.Model
	ta       textarea.Model
	ti       textinput.Model
	currNote note
	notes    []note
	width    int
	height   int
}

func NewModel(store *Store) model {
	notes, err := store.GetNotes()
	if err != nil {
		log.Fatalf("unable to get notes: %v", err)
	}
	ti := textinput.New()
	ti.Prompt = "Note Title :"

	return model{
		store: store,
		state: listView,
		ta:    textarea.New(),
		l:     NewList(notes),
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

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		m.ta.SetWidth(msg.Width)
		m.ta.SetHeight(msg.Height - 3)
		m.ti.SetWidth(msg.Width)
		m.l.SetHeight(msg.Height - 2) // leave room for header/footer
		m.l.SetWidth(msg.Width)

	case tea.KeyPressMsg:
		key := msg.String()
		switch m.state {
		case listView:
			m.l, cmd = m.l.Update(msg)
			cmds = append(cmds, cmd)

			if m.l.SettingFilter() {
				break
			}

			switch key {
			case "q":
				return m, tea.Quit
			case "n":
				m.ti.SetValue("")
				m.ti.Focus()
				m.currNote = note{}
				m.state = titleView
			case "enter":
				m.currNote = m.notes[m.l.Cursor()]
				m.state = bodyView
				m.ta.SetValue(m.currNote.body)
				m.ta.Focus()
				m.ta.CursorEnd()
			}

		// Title Input View key bindings
		case titleView:
			m.ti, cmd = m.ti.Update(msg)
			cmds = append(cmds, cmd)
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
			m.ta, cmd = m.ta.Update(msg)
			cmds = append(cmds, cmd)
			switch key {
			case "ctrl+s":
				m.currNote.body = m.ta.Value()

				var err error
				if err = m.store.SaveNote(m.currNote); err != nil {
					// TODO: handle error instead of quitting
					return m, tea.Quit
				}

				m.notes, err = m.store.GetNotes()
				m.l = NewList(m.notes)
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
