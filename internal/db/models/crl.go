package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CRL struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	Issuer   string             `bson:"issuer"`
	CRLBytes []byte             `bson:"crl_bytes"`
}

func NewCRL(issuer string, crlBytes []byte) *CRL {
	return &CRL{
		Issuer:   issuer,
		CRLBytes: crlBytes,
	}
}
