package modules

import "time"

type Note struct {
	Id        int
	Content     string
	CreatedAt time.Time
	User      string
}
