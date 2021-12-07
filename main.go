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

	//mux.HandleFunc("/logout", logout)

	mux.Handle("/favicon.ico", http.NotFoundHandler())

	log.Fatal(http.ListenAndServe("", mux))
}

func index(w http.ResponseWriter, req *http.Request) {
	views.Tpl().ExecuteTemplate(w, "index.gohtml", nil)
}
