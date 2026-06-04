package models

import "time"

type Post struct {
	ID         string
	UserID     string
	Username   string
	Title      string
	Content    string
	ImagePath  string
	Categories []Category
	Likes      int
	Dislikes   int
	UserVote   int
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type Category struct {
	ID   string
	Name string
}
