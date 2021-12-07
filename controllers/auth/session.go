package auth

import (
	"errors"
	"github.com/google/uuid"
	"net/http"
	. "photo_blog/models/user"
)

const sessionCookieName = "pb_session_id"

type sessionId string

var sessions map[sessionId]*User

func init() {
	sessions = make(map[sessionId]*User)
}

var ErrNoAuthUser = errors.New("user is not authorized")

func getAuthUser(req *http.Request) (*User, error) {
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

func CreateSession(user *User) (sessionId, error) {
	//TODO: realize create and save session
	sessId, err := uuid.NewRandom()
	if err != nil {
		return err
	}

}
