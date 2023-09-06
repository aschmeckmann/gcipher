package migratectl

import (
	"crypto/x509"
	"encoding/pem"
	"gcipher/internal/db/models"
	"gcipher/internal/db/repositories"
	"os"
	"path/filepath"
)

// MigrateCerts migrates certificates from a given directory into the database.
// The "username" parameter specifies the owner of the certificates.
func MigrateCerts(dir string, username string) error {
	files, err := os.ReadDir(dir)
	if err != nil {
		return err
	}

	for _, file := range files {
		// Skip directories and hidden files
		if file.IsDir() || file.Name()[0] == '.' {
			continue
		}

		certFile := filepath.Join(dir, file.Name())
		certPEMBytes, err := os.ReadFile(certFile)
		if err != nil {
			return err
		}

		// Decode the PEM file
		block, _ := pem.Decode(certPEMBytes)
		if block == nil || block.Type != "CERTIFICATE" {
			continue
		}

		// Parse the certificate to extract the serial number
		cert, err := x509.ParseCertificate(block.Bytes)
		if err != nil {
			return err
		}

		serialNumber := cert.SerialNumber.String()

		// Create a new Certificate model
		certModel := models.NewCertificate(serialNumber, certPEMBytes, username)

		// Insert the certificate into the database
		err = repositories.GetCertificateRepository().Insert(*certModel)
		if err != nil {
			return err
		}
	}

	return nil
}
