package ocsp

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"gcipher/internal/config"
	"gcipher/internal/db/models"
	"gcipher/internal/db/repositories"
	"gcipher/internal/server/api"
	"math/big"
	"net/http"
	"time"
)

// CRLUpdateInterval is the interval at which the CRL should be generated/updated
const CRLUpdateInterval = 24 * time.Hour // Generate/update CRL every 24 hours

func StartCRLUpdater() {
	// Start a goroutine to periodically generate/update the CRL
	go func() {
		for {
			updateCRL()
			time.Sleep(CRLUpdateInterval)
		}
	}()
}

func HandleCRL(w http.ResponseWriter, r *http.Request) {
	crl, err := repositories.GetCRLRepository().FindLatest()
	if err != nil {
		api.EncodeErrorResponse(w, http.StatusInternalServerError, "Failed to retrieve CRL")
		return
	}

	w.Header().Set("Content-Type", "application/pkix-crl")
	w.Header().Set("Content-Disposition", "attachment; filename=crl.crl")

	// Encode CRL bytes to PEM format
	pemBlock := &pem.Block{
		Type:  "X509 CRL",
		Bytes: crl.CRLBytes,
	}
	pem.Encode(w, pemBlock)
}

func updateCRL() {
	cfg, err := config.GetConfig()
	if err != nil {
		fmt.Println("Failed to get config:", err)
		return
	}

	certRepo := repositories.GetCertificateRepository()
	revokedCerts, err := certRepo.GetRevokedCertificates()
	if err != nil {
		fmt.Println("Failed to get revoked certificates:", err)
		return
	}

	var crlBytes []byte

	switch key := cfg.CAKey.(type) {
	case *rsa.PrivateKey:
		// Generate CRL using RSA private key
		crlBytes, err = generateCRL(cfg.CACert, key, revokedCerts)
		if err != nil {
			fmt.Println("Failed to generate CRL:", err)
			return
		}

	case *ecdsa.PrivateKey:
		// Generate CRL using EC private key
		crlBytes, err = generateCRL(cfg.CACert, key, revokedCerts)
		if err != nil {
			fmt.Println("Failed to generate CRL:", err)
			return
		}

	default:
		// Handle unsupported key type
		fmt.Println("Unsupported private key type")
		return
	}

	crl := models.NewCRL(cfg.CACert.Issuer.SerialNumber, crlBytes)
	err = repositories.GetCRLRepository().InsertOrUpdate(*crl)
	if err != nil {
		fmt.Println("Failed to insert/update CRL:", err)
		return
	}
}

func generateCRL(cert *x509.Certificate, key crypto.Signer, revokedCerts []models.Certificate) ([]byte, error) {
	template := x509.RevocationList{
		SignatureAlgorithm:  cert.SignatureAlgorithm,
		RevokedCertificates: []pkix.RevokedCertificate{},
		ThisUpdate:          time.Now(),
		NextUpdate:          time.Now().Add(CRLUpdateInterval),
	}

	for _, cert := range revokedCerts {
		serialNumber, _ := new(big.Int).SetString(cert.SerialNumber, 16)
		template.RevokedCertificates = append(template.RevokedCertificates, pkix.RevokedCertificate{
			SerialNumber:   serialNumber,
			RevocationTime: time.Now(),
		})
	}

	crlBytes, err := x509.CreateRevocationList(rand.Reader, &template, cert, key)
	if err != nil {
		return nil, err
	}

	return crlBytes, nil
}
