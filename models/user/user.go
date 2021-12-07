package user

import (
	"golang.org/x/crypto/bcrypt"
	"log"
)

type User struct {
	Id       int
	Login    string
	PassHash []byte
	Name     string
}

var users = map[string]*User{
	"admin": {Login: "admin", PassHash: []byte("admin")},
	"test":  {Login: "test", PassHash: []byte("test")},
}

type UnknownLoginError string

func (ul UnknownLoginError) Error() string {
	return "user not found with login:" + string(ul)
}

type BadPasswordError string

func (bp BadPasswordError) Error() string {
	return "bad password for user with login:" + string(bp)
}

func CheckPassword(login string, password string) (*User, error) {
	user, ok := users[login]
	if !ok {
		return nil, UnknownLoginError(login)
	}

	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		return nil, err
	}

	if err = bcrypt.CompareHashAndPassword(user.PassHash, passHash); err != nil {
		// TODO: delete it
		log.Println(string(passHash))

		if err != bcrypt.ErrMismatchedHashAndPassword {
			log.Println("unknown error in password matching: ", err)
		}

		return nil, BadPasswordError(login)
	}

	return user, nil
}
