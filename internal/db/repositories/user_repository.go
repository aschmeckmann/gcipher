package repositories

import (
	"context"
	"gcipher/internal/db"
	"gcipher/internal/db/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository struct {
	userCollection *mongo.Collection
}

func NewUserRepository() (*UserRepository, error) {
	client, err := db.GetDBClient()
	if err != nil {
		return nil, err
	}

	userCollection := client.Database("gcipher").Collection("users")
	return &UserRepository{userCollection: userCollection}, nil
}

func (repo *UserRepository) Insert(user models.User) error {
	_, err := repo.userCollection.InsertOne(context.Background(), user)
	return err
}

func (repo *UserRepository) FindByUsername(username string) (*models.User, error) {
	filter := bson.M{"username": username}
	var result models.User
	err := repo.userCollection.FindOne(context.Background(), filter).Decode(&result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (repo *UserRepository) Update(user models.User) error {
	filter := bson.M{"username": user.Username}
	update := bson.M{"$set": user}
	_, err := repo.userCollection.UpdateOne(context.Background(), filter, update)
	return err
}

func (repo *UserRepository) Delete(username string) error {
	filter := bson.M{"username": username}
	_, err := repo.userCollection.DeleteOne(context.Background(), filter)
	return err
}
