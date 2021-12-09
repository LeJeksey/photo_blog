package auth

import (
	"log"
	"net/http"
	"photo_blog/components/auth"
	"photo_blog/models/user"
	"photo_blog/views"
)

func Login(w http.ResponseWriter, r *http.Request) {
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

	authUser, err := user.CheckPassword(login, password)
	if err != nil {
		http.Error(w, "Bad login or password", http.StatusForbidden)
		log.Println("Bad try to login:", err)
		return
	}

	if _, err = auth.CreateSession(w, authUser); err != nil {
		http.Error(w, "Invalid session", http.StatusInternalServerError)
		log.Println("createSession:", err)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
