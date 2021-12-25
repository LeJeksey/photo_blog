package user

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
	"log"
	"photo_blog/components/db"
)

type User struct {
	Id       primitive.ObjectID `bson:"_id"`
	Login    string             `bson:"login"`
	PassHash []byte             `bson:"passHash"`
}

const usersCollectionName = "users"

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

func CheckPassword(ctx context.Context, login string, password string) (*User, error) {
	user, err := getUserByLogin(ctx, login)
	if err != nil {
		log.Println(err)
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

func getUserByLogin(ctx context.Context, login string) (*User, error) {
	uc, err := getUserCollection(ctx)
	if err != nil {
		return nil, err
	}

	res := uc.FindOne(ctx, bson.M{"login": login})
	if res.Err() != nil {
		return nil, res.Err()
	}

	var currUser User
	if err := res.Decode(&currUser); err == nil {
		return &currUser, nil
	} else {
		return nil, err
	}
}

func (u *User) Save(ctx context.Context) (primitive.ObjectID, error) {
	if u.Id != primitive.NilObjectID {
		err := u.update(ctx)
		if err != nil {
			return primitive.NilObjectID, err
		}
	} else {
		_, err := u.insert(ctx)
		if err != nil {
			return primitive.NilObjectID, err
		}
	}

	return u.Id, nil
}

var ErrDatabaseConnection = errors.New("error while connecting to database")
var ErrInsert = errors.New("error while inserting to database")

func (u *User) insert(ctx context.Context) (primitive.ObjectID, error) {
	uc, err := getUserCollection(ctx)
	if err != nil {
		log.Println("User.insert: ", err)
		return primitive.NilObjectID, ErrDatabaseConnection
	}

	u.Id = primitive.NewObjectID()
	_, err = uc.InsertOne(ctx, u)
	if err != nil {
		log.Println("User.insert: ", err)
		return primitive.NilObjectID, ErrInsert
	}

	return u.Id, nil
}

var ErrUpdate = errors.New("error while updating to database")

func (u *User) update(ctx context.Context) error {
	uc, err := getUserCollection(ctx)
	if err != nil {
		log.Println("User.update: ", err)
		return ErrDatabaseConnection
	}

	filter := bson.M{"_id": u.Id}
	updateData, err := db.GetStructAsBsonM(u)
	if err != nil {
		log.Println("User.update: ", err)
		return ErrUpdate
	}

	if _, err = uc.UpdateOne(ctx, filter, updateData); err != nil {
		log.Println("User.update: ", err)
		return ErrUpdate
	}

	return nil
}

func getUserCollection(ctx context.Context) (*mongo.Collection, error) {
	client, err := db.Client(ctx)
	if err != nil {
		return nil, err
	}

	return client.Database(db.CommonDbName).Collection(usersCollectionName), nil
}
