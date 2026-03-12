package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
	"github.com/microcosm-cc/bluemonday"
	"github.com/peterblog/blog/internal/auth"
	"github.com/peterblog/blog/internal/db"
	"github.com/peterblog/blog/internal/middleware"
)

var sanitizer = bluemonday.UGCPolicy()

func getComments(database *sqlx.DB) http.HandlerFunc {
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

		userID, _ := auth.GetSessionUserID(r)

		comments, err := db.GetCommentsByPostID(database, post.ID, userID)
		if err != nil {
			log.Printf("getComments: %v", err)
			jsonError(w, "failed to fetch comments", http.StatusInternalServerError)
			return
		}

		jsonResponse(w, comments)
	}
}

func createComment(database *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slug := chi.URLParam(r, "slug")
		userID, ok := middleware.GetUserID(r)
		if !ok {
			jsonError(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		var body struct {
			Content string `json:"content"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			jsonError(w, "invalid request body", http.StatusBadRequest)
			return
		}

		content := sanitizer.Sanitize(body.Content)
		if content == "" {
			jsonError(w, "comment content is required", http.StatusBadRequest)
			return
		}

		post, err := db.GetPostBySlug(database, slug)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				jsonError(w, "post not found", http.StatusNotFound)
				return
			}
			jsonError(w, "failed to fetch post", http.StatusInternalServerError)
			return
		}

		comment, err := db.CreateComment(database, post.ID, userID, content)
		if err != nil {
			jsonError(w, "failed to create comment", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		jsonResponse(w, comment)
	}
}

func addReaction(database *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")
		commentID, err := strconv.ParseUint(idStr, 10, 64)
		if err != nil {
			jsonError(w, "invalid comment id", http.StatusBadRequest)
			return
		}

		var body struct {
			Emoji string `json:"emoji"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			jsonError(w, "invalid request body", http.StatusBadRequest)
			return
		}

		if body.Emoji == "" {
			jsonError(w, "emoji is required", http.StatusBadRequest)
			return
		}
		if !isValidEmoji(body.Emoji) {
			jsonError(w, "invalid emoji", http.StatusBadRequest)
			return
		}

		userID, ok := middleware.GetUserID(r)
		if !ok {
			jsonError(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		reaction, err := db.ToggleReaction(database, commentID, userID, body.Emoji)
		if err != nil {
			jsonError(w, "failed to add reaction", http.StatusInternalServerError)
			return
		}

		jsonResponse(w, reaction)
	}
}

// isValidEmoji returns true if s looks like an emoji rather than text.
// Rejects anything containing ASCII characters (blocks letters, punctuation, angle braces, etc.)
// and enforces a max byte length. Intentionally permissive — blocks words/phrases, not every
// non-emoji Unicode character.
func isValidEmoji(s string) bool {
	if len(s) > 32 {
		return false
	}
	for _, r := range s {
		if r < 128 {
			return false
		}
	}
	return true
}
