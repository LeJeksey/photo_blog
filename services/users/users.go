package users

import (
	"context"
	"golang.org/x/crypto/bcrypt"
	"log"
	"photo_blog/models/user"
)

type ServiceInterface interface {
	CheckPassword(context.Context, string, string) (*user.User, error)
}

type usersService struct{}

var Service ServiceInterface

func init() {
	Service = &usersService{}
}

type UnknownLoginError string

func (ul UnknownLoginError) Error() string {
	return "user not found with login:" + string(ul)
}

type BadPasswordError string

func (bp BadPasswordError) Error() string {
	return "bad password for user with login:" + string(bp)
}

func (us *usersService) CheckPassword(ctx context.Context, login string, password string) (*user.User, error) {
	u, err := user.Dao.GetUserByLogin(ctx, login)
	if err != nil {
		log.Println(err)
		return nil, UnknownLoginError(login)
	}

	if err := bcrypt.CompareHashAndPassword(u.PassHash, []byte(password)); err != nil {
		if err != bcrypt.ErrMismatchedHashAndPassword {
			log.Println("unknown error in password matching: ", err)
		}

		return nil, BadPasswordError(login)
	}

	return u, nil
}
