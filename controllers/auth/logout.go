package auth

import (
	"log"
	"net/http"
	"photo_blog/components/auth"
)

func Logout() http.HandlerFunc {
	return logout
}
func logout(w http.ResponseWriter, r *http.Request) {
	if err := auth.DestroySession(w, r); err != nil {
		log.Println("Logout error:", err)
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
