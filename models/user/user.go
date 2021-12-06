package user

type User struct {
	Id       int
	Login    string
	PassHash string
	Name     string
}

var users = map[int]User{
	1: {Id: 1, Name: "admin", Login: "admin", PassHash: "admin"},
	2: {Id: 2, Name: "test", Login: "test", PassHash: "test"},
}
