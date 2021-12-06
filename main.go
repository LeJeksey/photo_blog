package main

import (
	"html/template"
	"net/http"
	"photo_blog/auth"
)

var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseGlob("templates/*"))
}

func main() {
	http.HandleFunc("/", index)
	http.HandleFunc("/login", login)
	http.HandleFunc("/logout", logout)
	http.Handle("/favicon.ico", http.NotFoundHandler())
	http.ListenAndServe(":8080", nil)
}

//test git user
//func authMiddleware(w http.ResponseWriter, r *http.Request)

func index(w http.ResponseWriter, req *http.Request) {
	authUser, err := auth.GetAuthUser(req)
	if err == auth.ErrNoAuthUser {
		http.Redirect(w, req, "/login", http.StatusSeeOther)
	}

	tpl.ExecuteTemplate(w, "index.gohtml", nil)
}
