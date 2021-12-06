package auth

import (
	"errors"
	"net/http"
	. "photo_blog/user"
)

const sessionCookieName = "pb_session_id"

var sessions map[string]*User

func init() {
	sessions = make(map[string]*User)
}

var ErrNoAuthUser = errors.New("user is not authorized")

func GetAuthUser(req *http.Request) (*User, error) {
	sCookie, err := req.Cookie(sessionCookieName)
	if err != nil {
		return nil, ErrNoAuthUser
	}

	if u, ok := sessions[sCookie.Value]; !ok {
		return nil, ErrNoAuthUser
	} else {
		return u, nil
	}
}
