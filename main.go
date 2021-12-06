package main

import (
	"html/template"
	"io"
	"log"
	"net/http"
	"photo_blog/controllers/auth"
)

var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseGlob("views/templates/*"))
}

func main() {
	mux := http.NewServeMux()

	mux.Handle("/", auth.Middleware(http.HandlerFunc(index)))
	mux.HandleFunc("/login", index2)

	//mux.HandleFunc("/logout", logout)

	mux.Handle("/favicon.ico", http.NotFoundHandler())

	log.Fatal(http.ListenAndServe("", mux))
}

func index(w http.ResponseWriter, req *http.Request) {
	tpl.ExecuteTemplate(w, "index.gohtml", nil)
}

func index2(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, "check")
	tpl.ExecuteTemplate(w, "index.gohtml", nil)
}
