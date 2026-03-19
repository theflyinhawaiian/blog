package db

import (
	"github.com/jmoiron/sqlx"
	"github.com/peterblog/blog/internal/models"
)

func GetCommentsByPostID(db *sqlx.DB, postID, userID uint64) ([]models.Comment, error) {
	var comments []models.Comment
	err := db.Select(&comments,
		`SELECT id, post_id, user_id, display_name, content, created_at, updated_at
		 FROM comments
		 WHERE post_id = ?
		 ORDER BY created_at ASC`, postID)
	if err != nil {
		return nil, err
	}

	for i := range comments {
		reactions, err := GetReactionsByCommentID(db, comments[i].ID, userID)
		if err != nil {
			return nil, err
		}
		comments[i].Reactions = reactions
	}

	return comments, nil
}

func CreateComment(db *sqlx.DB, postID, userID uint64, content string) (*models.Comment, error) {
	res, err := db.Exec(
		`INSERT INTO comments (post_id, user_id, display_name, content)
		 SELECT ?, id, display_name, ? FROM users WHERE id = ?`,
		postID, content, userID)
	if err != nil {
		return nil, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	var comment models.Comment
	err = db.Get(&comment,
		`SELECT id, post_id, user_id, display_name, content, created_at, updated_at
		 FROM comments WHERE id = ?`, id)
	if err != nil {
		return nil, err
	}
	comment.Reactions = []models.Reaction{}
	return &comment, nil
}

func GetReactionsByCommentID(db *sqlx.DB, commentID, userID uint64) ([]models.Reaction, error) {
	type row struct {
		Emoji string `db:"emoji"`
		Count int    `db:"count"`
	}
	var rows []row
	if err := db.Select(&rows,
		`SELECT emoji, COUNT(*) AS count FROM comment_reactions WHERE comment_id = ? GROUP BY emoji`,
		commentID); err != nil {
		return nil, err
	}
	if len(rows) == 0 {
		return []models.Reaction{}, nil
	}

	reacted := make(map[string]bool)
	if userID > 0 {
		var userEmojis []string
		if err := db.Select(&userEmojis,
			`SELECT emoji FROM comment_reactions WHERE comment_id = ? AND user_id = ?`,
			commentID, userID); err != nil {
			return nil, err
		}
		for _, e := range userEmojis {
			reacted[e] = true
		}
	}

	reactions := make([]models.Reaction, len(rows))
	for i, r := range rows {
		reactions[i] = models.Reaction{Emoji: r.Emoji, Count: r.Count, ReactedByMe: reacted[r.Emoji]}
	}
	return reactions, nil
}

func UpdateComment(db *sqlx.DB, commentID, userID uint64, content string) (*models.Comment, error) {
	res, err := db.Exec(
		`UPDATE comments SET content = ?, updated_at = NOW() WHERE id = ? AND user_id = ?`,
		content, commentID, userID)
	if err != nil {
		return nil, err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return nil, err
	}
	if affected == 0 {
		return nil, nil
	}

	var comment models.Comment
	err = db.Get(&comment,
		`SELECT id, post_id, user_id, display_name, content, created_at, updated_at
		 FROM comments WHERE id = ?`, commentID)
	if err != nil {
		return nil, err
	}
	comment.Reactions = []models.Reaction{}
	return &comment, nil
}

func DeleteComment(db *sqlx.DB, commentID, userID uint64) (bool, error) {
	res, err := db.Exec(
		`DELETE FROM comments WHERE id = ? AND user_id = ?`,
		commentID, userID)
	if err != nil {
		return false, err
	}
	affected, err := res.RowsAffected()
	return affected > 0, err
}

func ToggleReaction(db *sqlx.DB, commentID, userID uint64, emoji string) (*models.Reaction, error) {
	res, err := db.Exec(
		`DELETE FROM comment_reactions WHERE comment_id = ? AND user_id = ? AND emoji = ?`,
		commentID, userID, emoji)
	if err != nil {
		return nil, err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return nil, err
	}

	reactedByMe := false
	if affected == 0 {
		if _, err = db.Exec(
			`INSERT INTO comment_reactions (comment_id, user_id, emoji) VALUES (?, ?, ?)`,
			commentID, userID, emoji); err != nil {
			return nil, err
		}
		reactedByMe = true
	}

	var count int
	err = db.Get(&count,
		`SELECT COUNT(*) FROM comment_reactions WHERE comment_id = ? AND emoji = ?`,
		commentID, emoji)
	if err != nil {
		return nil, err
	}

	return &models.Reaction{Emoji: emoji, Count: count, ReactedByMe: reactedByMe}, nil
}
