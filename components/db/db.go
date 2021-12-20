package db

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var clientOptions *options.ClientOptions

const CommonDbName = "photoBlogDB"

func init() {
	clientOptions = options.Client().ApplyURI("mongodb://localhost:27017")
}

func Client(ctx context.Context) (*mongo.Client, error) {
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}

	if err = client.Ping(ctx, nil); err != nil {
		return nil, err
	}

	return client, nil
}
