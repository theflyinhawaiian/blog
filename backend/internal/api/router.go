package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
	"github.com/peterblog/blog/internal/middleware"
)

func NewRouter(database *sqlx.DB) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.CORS)

	r.Route("/api", func(r chi.Router) {
		r.Get("/posts", listPosts(database))
		r.Get("/posts/{slug}", getPost(database))
		r.Get("/posts/{slug}/comments", getComments(database))

		r.With(middleware.RequireAuth).Post("/posts/{slug}/comments", createComment(database))
		r.With(middleware.RequireAuth).Post("/comments/{id}/reactions", addReaction(database))
	})

	r.Route("/auth", func(r chi.Router) {
		r.Get("/{provider}", startAuth(database))
		r.Get("/{provider}/callback", callbackAuth(database))
		r.Get("/me", getMe(database))
		r.Post("/logout", logout(database))
	})

	return r
}
