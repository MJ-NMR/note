package main

import (
	"log"
	"os"
	"path/filepath"

	tea "charm.land/bubbletea/v2"
)

func main() {
	store := new(Store)
	if err := store.Init(dataDir()); err != nil {
		log.Printf("unable to init store: %v", err)
		log.Print(dataDir())
		return
	}

	m := NewModel(store)

	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		log.Fatalf("unable to run tui: %v", err)
	}
}

func dataDir() string {
	base := os.Getenv("XDG_DATA_HOME")
	if base == "" {
		base = "."
	}

	dir := filepath.Join(base, "notes")
	if err := os.MkdirAll(dir, 0755); err != nil {
		panic("could not create data dir: " + err.Error())
	}
	return filepath.Join(dir, "database.db")
}
