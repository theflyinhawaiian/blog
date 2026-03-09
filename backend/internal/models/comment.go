package models

import "time"

type Comment struct {
	ID          uint64     `db:"id" json:"id"`
	PostID      uint64     `db:"post_id" json:"post_id"`
	UserID      uint64     `db:"user_id" json:"user_id"`
	DisplayName string     `db:"display_name" json:"display_name"`
	Content     string     `db:"content" json:"content"`
	CreatedAt   time.Time  `db:"created_at" json:"created_at"`
	Reactions   []Reaction `db:"-" json:"reactions"`
}
