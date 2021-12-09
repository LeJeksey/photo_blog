package auth

import (
	"net/http"
	"photo_blog/models/user"
)

func GetUserFromSession(req *http.Request) *user.User {
	if req.Context().Value(userContextKey) == nil {
		return nil
	}

	return req.Context().Value(userContextKey).(*user.User)
}

func IsUserAuthorized(req *http.Request) bool {
	_, err := getAuthUser(req)
	return err == nil
}
