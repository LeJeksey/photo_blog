package auth

import (
	"log"
	"net/http"
)

func Logout(w http.ResponseWriter, r *http.Request) {
	if err := destroySession(w, r); err != nil {
		log.Println("Logout error:", err)
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
