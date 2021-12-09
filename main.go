package main

import (
	"log"
	"net/http"
	"photo_blog/components/auth"
	authController "photo_blog/controllers/auth"
	"photo_blog/controllers/pictures"
	"photo_blog/views"
)

func main() {
	mux := http.NewServeMux()

	mux.Handle("/", auth.Middleware(index()))
	mux.Handle("/login", authController.Login())
	mux.Handle("/logout", auth.Middleware(authController.Logout()))

	//mux.Handle("/upload", auth.Middleware(http.HandlerFunc()pictures.Upload))

	mux.Handle("/favicon.ico", http.NotFoundHandler())

	log.Fatal(http.ListenAndServe("", mux))
}

func index() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		authUser := auth.GetUserFromSession(req)
		pictures.StoreValuesToCookie(w, req)

		if err := views.Tpl().ExecuteTemplate(w, "index.gohtml", authUser); err != nil {
			log.Println(err)
		}
	}
}
