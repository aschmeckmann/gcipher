package repositories

import (
	"context"
	"gcipher/internal/db"
	"gcipher/internal/db/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type CRLRepository struct {
	crlCollection *mongo.Collection
}

func NewCRLRepository() (*CRLRepository, error) {
	client, err := db.GetDBClient()
	if err != nil {
		return nil, err
	}

	crlCollection := client.Database("gcipher").Collection("crls")
	return &CRLRepository{crlCollection: crlCollection}, nil
}

func (repo *CRLRepository) Insert(crl models.CRL) error {
	_, err := repo.crlCollection.InsertOne(context.Background(), crl)
	return err
}

func (repo *CRLRepository) InsertOrUpdate(crl models.CRL) error {
	filter := bson.M{"issuer": crl.Issuer}
	update := bson.M{"$set": crl, "$currentDate": bson.M{"updated_at": true}}
	opts := options.Update().SetUpsert(true)

	_, err := repo.crlCollection.UpdateOne(context.Background(), filter, update, opts)
	return err
}

func (repo *CRLRepository) FindLatest() (*models.CRL, error) {
	options := options.FindOne().SetSort(bson.M{"updated_at": -1})
	var result models.CRL
	err := repo.crlCollection.FindOne(context.Background(), bson.M{}, options).Decode(&result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
