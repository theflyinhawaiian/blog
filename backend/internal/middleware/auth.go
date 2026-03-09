package middleware

import (
	"context"
	"net/http"

	"github.com/peterblog/blog/internal/auth"
)

type contextKey string

const UserIDKey contextKey = "user_id"

func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID, ok := auth.GetSessionUserID(r)
		if !ok {
			http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), UserIDKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetUserID(r *http.Request) (uint64, bool) {
	id, ok := r.Context().Value(UserIDKey).(uint64)
	return id, ok
}
