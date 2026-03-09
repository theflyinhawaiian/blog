package models

type Reaction struct {
	ID        uint64 `db:"id" json:"id"`
	CommentID uint64 `db:"comment_id" json:"comment_id"`
	Emoji     string `db:"emoji" json:"emoji"`
	Count     uint32 `db:"count" json:"count"`
}
