package models

import "time"

type Post struct {
	ID       int64     `json:"id,omitempty" db:"id"`
	Parent   int64     `json:"parent,omitempty" db:"parent"`
	Author   string    `json:"author" db:"author"`
	Message  string    `json:"message,omitempty" db:"message"`
	IsEdited bool      `json:"isEdited" db:"is_edited"`
	Forum    string    `json:"forum,omitempty" db:"forum"`
	Thread   int64     `json:"thread" db:"thread"`
	Created  time.Time `json:"created,omitempty" db:"created"`
}

type PostDetails struct {
	Post   *Post   `json:"post,omitempty"`
	User   *User   `json:"author,omitempty"`
	Forum  *Forum  `json:"forum,omitempty"`
	Thread *Thread `json:"thread,omitempty"`
}
