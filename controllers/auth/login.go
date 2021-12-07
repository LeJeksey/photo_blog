package auth

import (
	"log"
	"net/http"
	"photo_blog/models/user"
	"photo_blog/views"
)

func Login(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		handleLogin(w, r)
		return
	}

	views.Tpl().ExecuteTemplate(w, "login.gohtml", nil)
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
	}

	// TODO: create session, write sessId to Cookie, redirect to index
	CreateSession(authUser)
}
