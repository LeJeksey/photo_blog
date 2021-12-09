package user

import (
	"golang.org/x/crypto/bcrypt"
	"log"
)

type User struct {
	Login    string
	PassHash []byte
}

var users = map[string]*User{
	"admin": {Login: "admin", PassHash: []byte("$2a$04$fU.P6Id4ntlMnQWwqny27OLKxwpDUEVhtVPtTfXj5SSGePOCeZnkK")},
	"test":  {Login: "test", PassHash: []byte("$2a$04$L6JsPcj86.PIFvD5alKVse2ZxhZz5VevOx3m3Q.tYrPF.Go0eeMgq")},
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

	if err := bcrypt.CompareHashAndPassword(user.PassHash, []byte(password)); err != nil {

		if err != bcrypt.ErrMismatchedHashAndPassword {
			log.Println("unknown error in password matching: ", err)
		}

		return nil, BadPasswordError(login)
	}

	return user, nil
}
