package db

import (
	"github.com/jmoiron/sqlx"
	"github.com/peterblog/blog/internal/models"
)

func GetCommentsByPostID(db *sqlx.DB, postID uint64) ([]models.Comment, error) {
	var comments []models.Comment
	err := db.Select(&comments,
		`SELECT c.id, c.post_id, c.user_id, u.display_name, c.content, c.created_at
		 FROM comments c
		 JOIN users u ON u.id = c.user_id
		 WHERE c.post_id = ?
		 ORDER BY c.created_at ASC`, postID)
	if err != nil {
		return nil, err
	}

	for i := range comments {
		reactions, err := GetReactionsByCommentID(db, comments[i].ID)
		if err != nil {
			return nil, err
		}
		comments[i].Reactions = reactions
	}

	return comments, nil
}

func CreateComment(db *sqlx.DB, postID, userID uint64, content string) (*models.Comment, error) {
	res, err := db.Exec(
		`INSERT INTO comments (post_id, user_id, content) VALUES (?, ?, ?)`,
		postID, userID, content)
	if err != nil {
		return nil, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	var comment models.Comment
	err = db.Get(&comment,
		`SELECT c.id, c.post_id, c.user_id, u.display_name, c.content, c.created_at
		 FROM comments c
		 JOIN users u ON u.id = c.user_id
		 WHERE c.id = ?`, id)
	if err != nil {
		return nil, err
	}
	comment.Reactions = []models.Reaction{}
	return &comment, nil
}

func GetReactionsByCommentID(db *sqlx.DB, commentID uint64) ([]models.Reaction, error) {
	var reactions []models.Reaction
	err := db.Select(&reactions,
		`SELECT id, comment_id, emoji, count FROM comment_reactions WHERE comment_id = ?`,
		commentID)
	if err != nil {
		return nil, err
	}
	if reactions == nil {
		reactions = []models.Reaction{}
	}
	return reactions, nil
}

func UpsertReaction(db *sqlx.DB, commentID uint64, emoji string) (*models.Reaction, error) {
	_, err := db.Exec(
		`INSERT INTO comment_reactions (comment_id, emoji, count)
		 VALUES (?, ?, 1)
		 ON DUPLICATE KEY UPDATE count = count + 1`,
		commentID, emoji)
	if err != nil {
		return nil, err
	}

	var reaction models.Reaction
	err = db.Get(&reaction,
		`SELECT id, comment_id, emoji, count FROM comment_reactions WHERE comment_id = ? AND emoji = ?`,
		commentID, emoji)
	return &reaction, err
}
