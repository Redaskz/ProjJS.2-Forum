package models

import "time"

type Like struct {
	ID         string
	UserID     string
	TargetID   string
	TargetType string
	Value      int
	CreatedAt  time.Time
}
