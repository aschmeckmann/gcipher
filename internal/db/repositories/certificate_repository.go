package repositories

import (
	"context"
	"gcipher/internal/db"
	"gcipher/internal/db/models"
	"time"

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

func (repo *CertificateRepository) FindBySerialNumberAndUsername(serialNumber, username string) (*models.Certificate, error) {
	filter := bson.M{"serial_number": serialNumber, "username": username}
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

func (repo *CertificateRepository) GetRevokedCertificates() ([]models.Certificate, error) {
	filter := bson.M{"revoked_at": bson.M{"$exists": true}}
	cursor, err := repo.certCollection.Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var revokedCerts []models.Certificate
	for cursor.Next(context.Background()) {
		var cert models.Certificate
		if err := cursor.Decode(&cert); err != nil {
			return nil, err
		}
		revokedCerts = append(revokedCerts, cert)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return revokedCerts, nil
}

func (repo *CertificateRepository) RevokeCertificate(serialNumber string) error {
	filter := bson.M{"serial_number": serialNumber}
	update := bson.M{"$set": bson.M{"revoked_at": time.Now()}}
	_, err := repo.certCollection.UpdateOne(context.Background(), filter, update)
	return err
}
