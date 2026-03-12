package models

type Reaction struct {
	Emoji       string `db:"emoji" json:"emoji"`
	Count       int    `db:"count" json:"count"`
	ReactedByMe bool   `db:"-" json:"reacted_by_me"`
}
