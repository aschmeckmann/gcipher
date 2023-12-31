package models

import "time"

type Certificate struct {
	SerialNumber   string     `bson:"serial_number"`
	CertificatePEM []byte     `bson:"certificate_pem"`
	Username       string     `bson:"username"`
	RevokedAt      *time.Time `bson:"revoked_at,omitempty"`
}

// Create a new certificate instance
func NewCertificate(serialNumber string, certificatePEM []byte, username string) *Certificate {
	return &Certificate{
		SerialNumber:   serialNumber,
		CertificatePEM: certificatePEM,
		Username:       username,
	}
}
