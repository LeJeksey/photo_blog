package auth

import (
	"log"
	"net/http"
	"photo_blog/components/auth"
	"photo_blog/services/users"
	"photo_blog/views"
)

func Login() http.HandlerFunc {
	return login
}

func login(w http.ResponseWriter, r *http.Request) {
	if auth.IsUserAuthorized(r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	if r.Method == http.MethodPost {
		handleLogin(w, r)
		return
	}

	if err := views.Tpl().ExecuteTemplate(w, "login.gohtml", nil); err != nil {
		log.Println(err)
	}
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Invalid login form", http.StatusInternalServerError)
		log.Println("Invalid login form", err)
		return
	}

	login := r.FormValue("login")
	password := r.FormValue("password")

	authUser, err := users.Service.CheckPassword(r.Context(), login, password)
	if err != nil {
		http.Error(w, "Bad login or password", http.StatusForbidden)
		log.Println("Bad try to login:", err)
		return
	}

	if _, err = auth.SessionManager.Create(w, authUser); err != nil {
		http.Error(w, "Invalid session", http.StatusInternalServerError)
		log.Println("createSession:", err)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
