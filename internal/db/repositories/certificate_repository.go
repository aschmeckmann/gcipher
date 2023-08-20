package repositories

import (
	"context"
	"gcipher/internal/db"
	"gcipher/internal/db/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type CertificateRepository struct {
	certCollection *mongo.Collection
}

func NewCertificateRepository() (*CertificateRepository, error) {
	client, err := db.GetDBClient()
	if err != nil {
		return nil, err
	}

	certCollection := client.Database("gcipher").Collection("certificates")
	return &CertificateRepository{certCollection: certCollection}, nil
}

func (repo *CertificateRepository) Insert(cert models.Certificate) error {
	_, err := repo.certCollection.InsertOne(context.Background(), cert)
	return err
}

func (repo *CertificateRepository) FindBySerialNumber(serialNumber string) (*models.Certificate, error) {
	filter := bson.M{"serial_number": serialNumber}
	var result models.Certificate
	err := repo.certCollection.FindOne(context.Background(), filter).Decode(&result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (repo *CertificateRepository) Update(cert models.Certificate) error {
	filter := bson.M{"serial_number": cert.SerialNumber}
	update := bson.M{"$set": cert}
	_, err := repo.certCollection.UpdateOne(context.Background(), filter, update)
	return err
}

func (repo *CertificateRepository) Delete(serialNumber string) error {
	filter := bson.M{"serial_number": serialNumber}
	_, err := repo.certCollection.DeleteOne(context.Background(), filter)
	return err
}
