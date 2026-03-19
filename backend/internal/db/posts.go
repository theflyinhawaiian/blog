package db

import (
	"github.com/jmoiron/sqlx"
	"github.com/peterblog/blog/internal/models"
)

type postTag struct {
	PostID uint64 `db:"post_id"`
	Name   string `db:"name"`
}

func fetchTagsForPosts(db *sqlx.DB, ids []uint64) (map[uint64][]string, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	query, args, err := sqlx.In(
		`SELECT pt.post_id, t.name
		 FROM post_tags pt JOIN tags t ON t.id = pt.tag_id
		 WHERE pt.post_id IN (?) ORDER BY t.name`, ids)
	if err != nil {
		return nil, err
	}
	var rows []postTag
	if err := db.Select(&rows, db.Rebind(query), args...); err != nil {
		return nil, err
	}
	tagMap := make(map[uint64][]string)
	for _, r := range rows {
		tagMap[r.PostID] = append(tagMap[r.PostID], r.Name)
	}
	return tagMap, nil
}

func ListPosts(db *sqlx.DB) ([]models.PostSummary, error) {
	var posts []models.PostSummary
	err := db.Select(&posts,
		`SELECT id, title, slug, excerpt, post_image, created_at
		 FROM posts ORDER BY created_at DESC`)
	if err != nil || len(posts) == 0 {
		return posts, err
	}

	ids := make([]uint64, len(posts))
	for i, p := range posts {
		ids[i] = p.ID
	}
	tagMap, err := fetchTagsForPosts(db, ids)
	if err != nil {
		return nil, err
	}
	for i := range posts {
		posts[i].Tags = tagMap[posts[i].ID]
	}
	return posts, nil
}

func GetPostBySlug(db *sqlx.DB, slug string) (*models.Post, error) {
	var post models.Post
	err := db.Get(&post,
		`SELECT id, title, slug, content, excerpt, meta_description, canonical_url, post_image, created_at, updated_at
		 FROM posts WHERE slug = ?`, slug)
	if err != nil {
		return nil, err
	}

	tagMap, err := fetchTagsForPosts(db, []uint64{post.ID})
	if err != nil {
		return nil, err
	}
	post.Tags = tagMap[post.ID]
	return &post, nil
}
