package db

import (
	"context"
	"gcipher/internal/config"
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
	cfg, err := config.GetConfig()
	if err != nil {
		return nil, err
	}

	clientOnce.Do(func() {
		clientOptions := options.Client().ApplyURI(cfg.DatabaseURL)
		client, err = mongo.Connect(context.Background(), clientOptions)
		if err != nil {
			panic(err)
		}
	})
	return client, err
}
