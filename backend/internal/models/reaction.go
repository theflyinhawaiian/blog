package models

type Reaction struct {
	Emoji       string `db:"emoji" json:"emoji"`
	Count       int    `db:"count" json:"count"`
	ReactedByMe bool   `db:"reacted_by_me" json:"reacted_by_me"`
}
