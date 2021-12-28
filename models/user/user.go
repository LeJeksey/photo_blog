package user

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"photo_blog/components/db"
)

type User struct {
	Id       primitive.ObjectID `bson:"_id"`
	Login    string             `bson:"login"`
	PassHash []byte             `bson:"passHash"`
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
