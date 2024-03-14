package database

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type Database struct {
	Client *mongo.Client
	Ctx    context.Context
}

func MustRun() *Database {
	ctx := context.Background()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		panic(err)
	}

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		panic(err)
	}

	err = client.Database("local").CreateCollection(ctx, "tokens")
	if err != nil {
		fmt.Println("collection already exists")
	}
	return &Database{
		Client: client,
		Ctx:    ctx,
	}
}
