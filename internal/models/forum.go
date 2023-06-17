package models

type Forum struct {
	Title   string `json:"title,omitempty" db:"title"`
	User    string `json:"user,omitempty" db:"user_nickname"`
	Slug    string `json:"slug,omitempty" db:"slug"`
	Posts   int64  `json:"posts,omitempty" db:"posts"`
	Threads int64  `json:"threads,omitempty" db:"threads"`
}

type ForumUsers struct {
	Forum string `db:"forum"`
	User  string `db:"user_nickname"`
}
