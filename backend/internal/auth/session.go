package auth

import (
	"net/http"
	"os"

	"github.com/gorilla/sessions"
)

const sessionName = "blog_session"
const sessionUserKey = "user_id"

var store *sessions.CookieStore

func InitSessionStore() {
	secret := os.Getenv("SESSION_SECRET")
	if secret == "" {
		secret = "dev-secret-change-me"
	}
	store = sessions.NewCookieStore([]byte(secret))
	store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 30, // 30 days
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   false, // set true in production
	}
}

func GetSession(r *http.Request) (*sessions.Session, error) {
	return store.Get(r, sessionName)
}

func SetUserSession(w http.ResponseWriter, r *http.Request, userID uint64) error {
	session, err := store.Get(r, sessionName)
	if err != nil {
		return err
	}
	session.Values[sessionUserKey] = userID
	return session.Save(r, w)
}

func GetSessionUserID(r *http.Request) (uint64, bool) {
	session, err := store.Get(r, sessionName)
	if err != nil {
		return 0, false
	}
	id, ok := session.Values[sessionUserKey]
	if !ok {
		return 0, false
	}
	userID, ok := id.(uint64)
	return userID, ok
}

func ClearSession(w http.ResponseWriter, r *http.Request) error {
	session, err := store.Get(r, sessionName)
	if err != nil {
		return err
	}
	session.Options.MaxAge = -1
	return session.Save(r, w)
}
