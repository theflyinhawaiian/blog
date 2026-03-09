package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
	"github.com/peterblog/blog/internal/db"
)

func listPosts(database *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		posts, err := db.ListPosts(database)
		if err != nil {
			jsonError(w, "failed to fetch posts", http.StatusInternalServerError)
			return
		}
		jsonResponse(w, posts)
	}
}

func getPost(database *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slug := chi.URLParam(r, "slug")
		post, err := db.GetPostBySlug(database, slug)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				jsonError(w, "post not found", http.StatusNotFound)
				return
			}
			jsonError(w, "failed to fetch post", http.StatusInternalServerError)
			return
		}
		jsonResponse(w, post)
	}
}

func jsonResponse(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(v)
}

func jsonError(w http.ResponseWriter, msg string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": msg})
}
