package database

import (
	"database/sql"
	"os"

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

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			password TEXT NOT NULL
		);

		CREATE TABLE IF NOT EXISTS notes (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			content TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id)
		);
		`)

	if err != nil {
		return nil, err
	}
	return &DB{db}, nil
}

type DB struct {
	db *sql.DB
}

func (d DB) GetAllNots(userId string) ([]Note, error) {
	notes := []Note{}
	rows, err := d.db.Query("Select * from notes where user_id=?;", userId)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var note Note
		err = rows.Scan(&note.Id, &note.User_id, &note.Content, &note.CreatedAt)
		if err != nil {
			return nil, err
		}
		notes = append(notes, note)
	}

	return notes, nil
}

func (d DB) GetOneNote(id, userId string) (*Note, error) {
	note := Note{}
	row := d.db.QueryRow("select * from notes where id=? and user_id=?;", id, userId)
	err := row.Scan(&note.Id, &note.Content, &note.CreatedAt, &note.User_id)
	return &note, err
}

func (d DB) DeleteOneNote(id, user_id string) error {
	_, err := d.db.Exec("delete from notes where id=? and user_id=?;", id, user_id)
	return err
}

func (d DB) AddOneNote(userId, constant string) error {
	_, err := d.db.Exec("insert into notes (user_id, content) values (?,?);", userId, constant)
	return err
}


func (d DB) AddUser(name, password string) (int64, error) {
	res, err := d.db.Exec("insert into users (name, password) values (?,?);", name, password)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	return id, err
}

func (d DB) GetUser(name, password string) (int64, error) {
	row := d.db.QueryRow("select id from users where name=? and password=?;", name, password)

	var id int64
	err := row.Scan(&id)
	return id, err
}
