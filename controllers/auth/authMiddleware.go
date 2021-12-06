package auth

import (
	"log"
	"net/http"
	"photo_blog/components/auth"
)

func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authUser, err := auth.GetAuthUser(r)
		if err != nil {
			if err == auth.ErrNoAuthUser {
				http.Redirect(w, r, "/login", http.StatusSeeOther)
				return
			}
			log.Println("Error while checking authUser:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		_ = authUser
		next.ServeHTTP(w, r)
	})
}
