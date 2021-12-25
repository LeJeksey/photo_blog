package auth

import (
	"errors"
	"github.com/google/uuid"
	"net/http"
	. "photo_blog/models/user"
)

const sessionCookieName = "pb_session_id"
const sessionMaxAge = 3600 * 4

type SessionId string

// todo: here i can use redis or analog, but i didn't want =)
var sessions map[SessionId]*User

func init() {
	sessions = make(map[SessionId]*User)
}

var ErrNoAuthUser = errors.New("user is not authorized")

func getAuthUser(req *http.Request) (*User, error) {
	sessId, err := getSessionIdFromCookie(req)
	if err != nil {
		return nil, ErrNoAuthUser
	}

	if u, ok := sessions[sessId]; !ok {
		return nil, ErrNoAuthUser
	} else {
		return u, nil
	}
}

var ErrNoSessCookie = errors.New("session cookie is empty")

func getSessionIdFromCookie(req *http.Request) (SessionId, error) {
	sCookie, err := req.Cookie(sessionCookieName)
	if err != nil {
		return "", ErrNoSessCookie
	}

	sessId := SessionId(sCookie.Value)
	return sessId, nil
}

func CreateSession(w http.ResponseWriter, user *User) (SessionId, error) {
	sUuid, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}

	// TODO: need to save some time for destroying session on the server side. But i don't want to do it now
	sessId := SessionId(sUuid.String())
	sessions[sessId] = user

	sCookie := &http.Cookie{
		Name:   sessionCookieName,
		Value:  string(sessId),
		MaxAge: sessionMaxAge,
	}
	http.SetCookie(w, sCookie)

	return sessId, nil
}

func DestroySession(w http.ResponseWriter, req *http.Request) error {
	sessId, err := getSessionIdFromCookie(req)
	if err != nil {
		return err
	}

	if _, ok := sessions[sessId]; ok {
		delete(sessions, sessId)
	}

	sCookie := &http.Cookie{
		Name:   sessionCookieName,
		Value:  "",
		MaxAge: -1,
	}
	http.SetCookie(w, sCookie)

	return nil
}
