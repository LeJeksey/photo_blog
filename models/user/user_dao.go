package user

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"photo_blog/components/db"
)

type DaoInterface interface {
	GetUserByLogin(context.Context, string) (*User, error)
}

type userDao struct{}

var Dao DaoInterface

func init() {
	Dao = &userDao{}
}

const usersCollectionName = "users"

func (ud *userDao) GetUserByLogin(ctx context.Context, login string) (*User, error) {
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

var getUserCollection = func(ctx context.Context) (*mongo.Collection, error) {
	client, err := db.Client(ctx)
	if err != nil {
		return nil, err
	}

	return client.Database(db.CommonDbName).Collection(usersCollectionName), nil
}
