package database

import (
	"database/sql"

	"github.com/MJ-NMR/note/modules"
	_ "modernc.org/sqlite"
)

func NewDB() (*DB, error) {
	db, err := sql.Open("sqlite", "database/database.db")
	if err != nil {
		return nil, err
	}
	return &DB{db}, nil
}

type DB struct {
	db *sql.DB
}

func (d DB) GetAllNots() ([]modules.Note, error) {
	notes := []modules.Note{}
	rows, err := d.db.Query("Select * from notes;")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var note modules.Note
		err = rows.Scan(&note.Id, &note.Content, &note.CreatedAt, &note.User)
		if err != nil {
			return nil, err
		}
		notes = append(notes, note)
	}

	return notes, nil
}

func (d DB) GetOneNote(id string) (*modules.Note, error) {
	note := modules.Note{}
	row := d.db.QueryRow("select * from notes where id=?;", id)
	row.Scan(&note.Id, &note.Content, &note.CreatedAt, &note.User)
	return &note, nil
}
