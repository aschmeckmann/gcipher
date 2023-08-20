package db

import (
	"context"
	"sync"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	client     *mongo.Client
	clientOnce sync.Once
)

func GetDBClient() (*mongo.Client, error) {
	var err error
	clientOnce.Do(func() {
		clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
		client, err = mongo.Connect(context.Background(), clientOptions)
		if err != nil {
			panic(err)
		}
	})
	return client, err
}
