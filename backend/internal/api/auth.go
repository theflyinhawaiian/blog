package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
	"github.com/peterblog/blog/internal/auth"
	dbpkg "github.com/peterblog/blog/internal/db"
	"golang.org/x/oauth2"
)

func startAuth(database *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		providerName := chi.URLParam(r, "provider")

		cfg, ok := auth.GetProvider(providerName)
		if !ok {
			jsonError(w, "unknown provider", http.StatusBadRequest)
			return
		}

		state := auth.OAuthStateParam()
		if err := auth.StoreOAuthState(w, r, state); err != nil {
			jsonError(w, "session error", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, auth.AuthCodeURL(cfg, state), http.StatusFound)
	}
}

func callbackAuth(database *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		providerName := chi.URLParam(r, "provider")

		cfg, ok := auth.GetProvider(providerName)
		if !ok {
			http.Redirect(w, r, auth.RedirectToFrontend("/login?error=unknown_provider"), http.StatusFound)
			return
		}

		code := r.URL.Query().Get("code")
		state := r.URL.Query().Get("state")

		if !auth.ValidateOAuthState(r, state) {
			http.Redirect(w, r, auth.RedirectToFrontend("/login?error=invalid_state"), http.StatusFound)
			return
		}

		var token *oauth2.Token
		var exchangeErr error

		token, exchangeErr = auth.ExchangeCode(r.Context(), cfg, code)

		if exchangeErr != nil {
			http.Redirect(w, r, auth.RedirectToFrontend("/login?error=token_exchange"), http.StatusFound)
			return
		}

		userInfo, err := auth.FetchUserInfo(r.Context(), providerName, cfg, token)
		if err != nil {
			http.Redirect(w, r, auth.RedirectToFrontend("/login?error=user_info"), http.StatusFound)
			return
		}

		user, err := dbpkg.UpsertUser(database, providerName, userInfo.ProviderUserID, userInfo.DisplayName, userInfo.Email)
		if err != nil {
			http.Redirect(w, r, auth.RedirectToFrontend("/login?error=db"), http.StatusFound)
			return
		}

		if err := auth.SetUserSession(w, r, user.ID); err != nil {
			http.Redirect(w, r, auth.RedirectToFrontend("/login?error=session"), http.StatusFound)
			return
		}

		http.Redirect(w, r, auth.RedirectToFrontend("/"), http.StatusFound)
	}
}

func getMe(database *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := auth.GetSessionUserID(r)
		if !ok {
			jsonError(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		user, err := dbpkg.GetUserByID(database, userID)
		if err != nil {
			jsonError(w, "user not found", http.StatusNotFound)
			return
		}

		jsonResponse(w, user)
	}
}

func logout(database *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := auth.ClearSession(w, r); err != nil {
			jsonError(w, "session error", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}
