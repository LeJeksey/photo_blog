package main

import (
	"log"
	"net/http"
	"photo_blog/controllers/auth"
	"photo_blog/views"
)

func main() {
	mux := http.NewServeMux()

	mux.Handle("/", auth.Middleware(http.HandlerFunc(index)))
	mux.HandleFunc("/login", auth.Login)
	mux.Handle("/logout", auth.Middleware(http.HandlerFunc(auth.Logout)))

	mux.Handle("/favicon.ico", http.NotFoundHandler())

	log.Fatal(http.ListenAndServe("", mux))
}

func index(w http.ResponseWriter, req *http.Request) {
	authUser := auth.GetUserFromSession(req)

	if err := views.Tpl().ExecuteTemplate(w, "index.gohtml", authUser); err != nil {
		log.Println(err)
	}
}
