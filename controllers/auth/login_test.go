package auth

import (
	"errors"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"photo_blog/components/auth"
	"photo_blog/mocks"
	"photo_blog/models/user"
	"photo_blog/services/users"
	"strings"
	"testing"
)

func TestHandleLoginInvalidLoginForm(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "http://localhost/login", nil)
	req.Body = nil
	w := httptest.NewRecorder()

	handleLogin(w, req)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.EqualValues(t, http.StatusInternalServerError, resp.StatusCode)
	assert.EqualValues(t, "Invalid login form\n", string(body))
}

func TestHandleLoginBadLoginOrPassword(t *testing.T) {
	login := "myBadLogin"
	password := "orMyBadPassword"

	req := httptest.NewRequest(http.MethodPost, "http://localhost/login",
		strings.NewReader(fmt.Sprintf("login=%s&password=%s", login, password)))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	w := httptest.NewRecorder()

	ctrl := gomock.NewController(t)
	uServiceMock := mocks.NewMockUsersServiceInterface(ctrl)

	uServiceMock.
		EXPECT().
		CheckPassword(gomock.Eq(req.Context()), gomock.Eq(login), gomock.Eq(password)).
		Return(nil, users.BadPasswordError(login))
	users.Service = uServiceMock

	handleLogin(w, req)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.EqualValues(t, http.StatusForbidden, resp.StatusCode)
	assert.EqualValues(t, "Bad login or password\n", string(body))
}

func TestHandleLoginInvalidSession(t *testing.T) {
	login := "myLogin"
	password := "myPassword"

	req := httptest.NewRequest(http.MethodPost, "http://localhost/login",
		strings.NewReader(fmt.Sprintf("login=%s&password=%s", login, password)))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	w := httptest.NewRecorder()

	ctrl := gomock.NewController(t)

	authUser := &user.User{Login: login, PassHash: []byte(password)} // it's not really how password expected
	uServiceMock := mocks.NewMockUsersServiceInterface(ctrl)
	uServiceMock.
		EXPECT().
		CheckPassword(gomock.Eq(req.Context()), gomock.Eq(login), gomock.Eq(password)).
		Return(authUser, nil)
	users.Service = uServiceMock

	authSessionManagerMock := mocks.NewMockSessionManagerInterface(ctrl)
	authSessionManagerMock.
		EXPECT().
		Create(gomock.Eq(w), gomock.Eq(authUser)).
		Return(auth.SessionId(""), errors.New("can't create session"))
	auth.SessionManager = authSessionManagerMock

	handleLogin(w, req)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.EqualValues(t, http.StatusInternalServerError, resp.StatusCode)
	assert.EqualValues(t, "Invalid session\n", string(body))
}
func TestHandleLoginOk(t *testing.T) {
	login := "myLogin"
	password := "myPassword"

	req := httptest.NewRequest(http.MethodPost, "http://localhost/login",
		strings.NewReader(fmt.Sprintf("login=%s&password=%s", login, password)))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	w := httptest.NewRecorder()

	ctrl := gomock.NewController(t)

	authUser := &user.User{Login: login, PassHash: []byte(password)} // it's not really how password expected
	uServiceMock := mocks.NewMockUsersServiceInterface(ctrl)
	uServiceMock.
		EXPECT().
		CheckPassword(gomock.Eq(req.Context()), gomock.Eq(login), gomock.Eq(password)).
		Return(authUser, nil)
	users.Service = uServiceMock

	authSessionManagerMock := mocks.NewMockSessionManagerInterface(ctrl)
	authSessionManagerMock.
		EXPECT().
		Create(gomock.Eq(w), gomock.Eq(authUser)).
		Return(auth.SessionId("0xff123456789"), nil)
	auth.SessionManager = authSessionManagerMock

	handleLogin(w, req)

	resp := w.Result()

	assert.EqualValues(t, http.StatusSeeOther, resp.StatusCode)
}
