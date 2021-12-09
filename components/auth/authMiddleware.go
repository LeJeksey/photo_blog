package auth

import (
	"context"
	"log"
	"net/http"
)

type contextKey string

const userContextKey contextKey = "userContext"

func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authUser, err := getAuthUser(r)
		if err != nil {
			if err == ErrNoAuthUser {
				http.Redirect(w, r, "/login", http.StatusSeeOther)
				return
			}
			log.Println("Error while checking authUser:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		ctxWithUser := context.WithValue(r.Context(), userContextKey, authUser)

		next.ServeHTTP(w, r.WithContext(ctxWithUser))
	})
}
