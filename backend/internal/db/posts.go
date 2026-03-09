package db

import (
	"github.com/jmoiron/sqlx"
	"github.com/peterblog/blog/internal/models"
)

func ListPosts(db *sqlx.DB) ([]models.PostSummary, error) {
	var posts []models.PostSummary
	err := db.Select(&posts,
		`SELECT id, title, slug, excerpt, post_image, tags, created_at
		 FROM posts ORDER BY created_at DESC`)
	return posts, err
}

func GetPostBySlug(db *sqlx.DB, slug string) (*models.Post, error) {
	var post models.Post
	err := db.Get(&post,
		`SELECT id, title, slug, content, excerpt, meta_description, canonical_url, post_image, tags, created_at, updated_at
		 FROM posts WHERE slug = ?`, slug)
	if err != nil {
		return nil, err
	}
	return &post, nil
}
