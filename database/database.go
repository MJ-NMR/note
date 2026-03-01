package database

import (
	"database/sql"
	"os"

	"github.com/MJ-NMR/note/modules"
	_ "modernc.org/sqlite"
)

func GetDBPath() (string, error) {
	path := os.Getenv("notesDB")
	if path == "" {
		path = os.Getenv("HOME") + "/.config/notes"
	}

	_, err := os.Stat(path)

	if os.IsNotExist(err) {
		err = os.MkdirAll(path, 0755)
	}
	return path, err
}

func NewDBConnection() (*DB, error) {
	dbPath, err := GetDBPath()
	if err != nil {
		return nil, err
	}
	db, err := sql.Open("sqlite", dbPath+"/database.db")
	if err != nil {
		return nil, err
	}
	_, err = db.Exec(
		`CREATE TABLE IF NOT EXISTS notes (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			content TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			user string NOT NULL
		);`)

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
	err := row.Scan(&note.Id, &note.Content, &note.CreatedAt, &note.User)
	return &note, err
}

func (d DB) DeleteOneNote(id string) error {
	_, err := d.db.Exec("delete from notes where id=?", id)
	return err
}

func (d DB) AddOneNote(user, constant string) error {
	_, err := d.db.Exec("insert into notes (user, content) values (?,?)", user, constant)
	return  err
}
