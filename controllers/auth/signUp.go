package auth

import (
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"photo_blog/models/user"
	"photo_blog/views"
)

func SignUp() http.HandlerFunc {
	return signUp
}

func signUp(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		handleSignUp(w, r)
	}

	if err := views.Tpl().ExecuteTemplate(w, "signUp.gohtml", nil); err != nil {
		log.Println(err)
	}
}

func handleSignUp(w http.ResponseWriter, r *http.Request) {
	login := r.FormValue("login")
	password := r.FormValue("password")

	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		log.Println("handleSignUp GenerateFromPassword:", err)
		return
	}

	u := &user.User{Login: login, PassHash: passHash}
	_, err = u.Save(r.Context())
	if err != nil {
		http.Error(w, "Can't sign up", http.StatusInternalServerError)
		log.Println("handleSignUp User.Save:", err)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
