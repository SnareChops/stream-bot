package db

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	Database *mongo.Database
)

func Start(connString, dbName string) func() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connString))
	if err != nil {
		panic(err)
	}
	Database = client.Database(dbName)
	return func() {
		if err := client.Disconnect(context.Background()); err != nil {
			panic(err)
		}
	}
}
