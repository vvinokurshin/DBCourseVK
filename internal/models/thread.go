package models

import "github.com/go-openapi/strfmt"

type Thread struct {
	ID      int64           `json:"id,omitempty" db:"id"`
	Title   string          `json:"title,omitempty" db:"title"`
	Author  string          `json:"author,omitempty" db:"author"`
	Forum   string          `json:"forum,omitempty" db:"forum"`
	Message string          `json:"message,omitempty" db:"message"`
	Votes   int64           `json:"votes,omitempty" db:"votes"`
	Slug    string          `json:"slug,omitempty" db:"slug"`
	Created strfmt.DateTime `json:"created,omitempty" db:"created"`
}

type Vote struct {
	Thread   int64  `json:"thread" db:"thread"`
	Nickname string `json:"nickname" db:"nickname"`
	Voice    int64  `json:"voice" db:"voice"`
}
