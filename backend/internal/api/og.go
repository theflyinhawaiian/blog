package api

import (
	"database/sql"
	"errors"
	"html/template"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
	"github.com/peterblog/blog/internal/auth"
	dbpkg "github.com/peterblog/blog/internal/db"
)

var ogTemplate = template.Must(template.New("og").Parse(`<!DOCTYPE html>
<html>
<head>
  <meta charset="UTF-8">
  <title>{{.Title}}</title>
  <meta name="description" content="{{.Description}}">
  <meta property="og:title" content="{{.Title}}">
  <meta property="og:description" content="{{.Description}}">
  <meta property="og:type" content="article">
  <meta property="og:url" content="{{.URL}}">
  {{- if .Image}}
  <meta property="og:image" content="{{.Image}}">
  {{- end}}
  <meta name="twitter:card" content="{{if .Image}}summary_large_image{{else}}summary{{end}}">
  <meta name="twitter:title" content="{{.Title}}">
  <meta name="twitter:description" content="{{.Description}}">
  {{- if .Image}}
  <meta name="twitter:image" content="{{.Image}}">
  {{- end}}
  <meta http-equiv="refresh" content="0;url={{.URL}}">
</head>
<body>
  <a href="{{.URL}}">{{.Title}}</a>
</body>
</html>
`))

type ogData struct {
	Title       string
	Description string
	URL         string
	Image       string
}

func getOGTags(database *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slug := chi.URLParam(r, "slug")

		post, err := dbpkg.GetPostBySlug(database, slug)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				http.NotFound(w, r)
				return
			}
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		description := ""
		if post.MetaDescription.Valid {
			description = post.MetaDescription.String
		} else if post.Excerpt.Valid {
			description = post.Excerpt.String
		}

		image := ""
		if post.PostImage.Valid {
			image = post.PostImage.String
		}

		url := auth.RedirectToFrontend("/posts/" + slug)
		if post.CanonicalURL.Valid {
			url = post.CanonicalURL.String
		}

		data := ogData{
			Title:       "Pete's blog - " + post.Title,
			Description: description,
			URL:         url,
			Image:       image,
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_ = ogTemplate.Execute(w, data)
	}
}
