package database

import "time"

type Note struct {
	Id        int64
	Content   string
	CreatedAt time.Time
	User_id   int
}

type User struct {
	Id       int64
	Name     string
	Password string
}
