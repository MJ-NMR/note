package main

import (
	"database/sql"
	"time"

	_ "modernc.org/sqlite"
)

type Store struct {
	conn *sql.DB
}

func (s *Store) Init() error {
	var err error
	s.conn, err = sql.Open("sqlite", "./database.db")
	if err != nil {
		return err
	}

	createTableStmt := `CREATE TABLE IF NOT EXISTS notes (
		id integer not null primary key,
		title text not null,
		body text not null
	);`

	if _, err := s.conn.Exec(createTableStmt); err != nil {
		return err
	}

	return nil
}

func (s *Store) GetNotes() ([]note, error) {
	rows, err := s.conn.Query("SELECT * FROM notes")
	if err != nil {
		return nil, err
	}

	notes := []note{}
	defer rows.Close()
	for rows.Next() {
		note := note{}
		rows.Scan(&note.id, &note.title, &note.body)
		notes = append(notes, note)
	}

	return notes, nil
}

func (s *Store) SaveNote(note note) error {
	if note.id == 0 {
		// pseudo-unique id
		note.id = time.Now().UTC().Unix()
	}

	upsertQuery := `INSERT INTO notes (id, title, body)
	VALUES (?, ?, ?)
	ON CONFLICT(id) DO UPDATE
	SET title=excluded.title, body=excluded.body;`

	if _, err := s.conn.Exec(upsertQuery, note.id, note.title, note.body); err != nil {
		return err
	}

	return nil
}
