package models

import (
	"database/sql"
	"time"
)

type Post struct {
	ID              uint64         `db:"id" json:"id"`
	Title           string         `db:"title" json:"title"`
	Slug            string         `db:"slug" json:"slug"`
	Content         string         `db:"content" json:"content,omitempty"`
	Excerpt         sql.NullString `db:"excerpt" json:"excerpt,omitempty"`
	MetaDescription sql.NullString `db:"meta_description" json:"meta_description,omitempty"`
	CanonicalURL    sql.NullString `db:"canonical_url" json:"canonical_url,omitempty"`
	PostImage       sql.NullString `db:"post_image" json:"post_image,omitempty"`
	Tags            []string       `db:"-" json:"tags,omitempty"`
	CreatedAt       time.Time      `db:"created_at" json:"created_at"`
	UpdatedAt       time.Time      `db:"updated_at" json:"updated_at"`
}

type PostSummary struct {
	ID        uint64         `db:"id" json:"id"`
	Title     string         `db:"title" json:"title"`
	Slug      string         `db:"slug" json:"slug"`
	Excerpt   sql.NullString `db:"excerpt" json:"excerpt,omitempty"`
	PostImage sql.NullString `db:"post_image" json:"post_image,omitempty"`
	Tags      []string       `db:"-" json:"tags,omitempty"`
	CreatedAt time.Time      `db:"created_at" json:"created_at"`
}
