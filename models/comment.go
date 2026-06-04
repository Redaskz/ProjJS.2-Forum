package models

import "time"

type Comment struct {
	ID        string
	PostID    string
	UserID    string
	Username  string
	Content   string
	Likes     int
	Dislikes  int
	UserVote  int
	CreatedAt time.Time
	UpdatedAt time.Time
}
