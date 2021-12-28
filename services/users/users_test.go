package users

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"photo_blog/models/user"
	"testing"
)

type userDaoMock struct{}

var getUserByLogin func(ctx context.Context, login string) (*user.User, error)

func (m *userDaoMock) GetUserByLogin(ctx context.Context, login string) (*user.User, error) {
	return getUserByLogin(ctx, login)
}

func init() {
	user.Dao = &userDaoMock{}
}

func TestCheckPasswordUnknownLoginError(t *testing.T) {
	getUserByLogin = func(ctx context.Context, login string) (*user.User, error) {
		return nil, errors.New("can't get user by login")
	}

	ctx := context.Background()
	login := "test_login"
	password := "pswd"

	u, err := Service.CheckPassword(ctx, login, password)

	assert.Nil(t, u)
	assert.NotNil(t, err)
	assert.EqualValues(t, UnknownLoginError(login), err)
}

func TestCheckPasswordBadPasswordError(t *testing.T) {
	ctx := context.Background()
	login := "test_login"
	password := "pswd"

	passHash, _ := bcrypt.GenerateFromPassword([]byte(password+"1"), bcrypt.MinCost)
	expectedUser := &user.User{Login: login, PassHash: passHash}

	getUserByLogin = func(ctx context.Context, login string) (*user.User, error) {
		return expectedUser, nil
	}

	u, err := Service.CheckPassword(ctx, login, password)

	assert.Nil(t, u)
	assert.NotNil(t, err)
	assert.EqualValues(t, BadPasswordError(login), err)
}

func TestCheckPasswordOk(t *testing.T) {
	ctx := context.Background()
	login := "test_login"
	password := "pswd"

	passHash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	expectedUser := &user.User{Login: login, PassHash: passHash}

	getUserByLogin = func(ctx context.Context, login string) (*user.User, error) {
		return expectedUser, nil
	}

	u, err := Service.CheckPassword(ctx, login, password)

	assert.Nil(t, err)
	assert.NotNil(t, u)

	assert.EqualValues(t, expectedUser, u)
}
