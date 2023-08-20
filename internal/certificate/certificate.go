package certificate

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"gcipher/internal/config"
	"gcipher/internal/db/models"
	"gcipher/internal/db/repositories"
	"gcipher/internal/server/api"
	"gcipher/internal/user"
	"math/big"
	"net/http"
	"time"
)

// HandleCertificateRequest handles incoming certificate signing requests (CSRs) and generates signed certificates.
func HandleCertificateRequest(w http.ResponseWriter, r *http.Request) {
	var request api.Request

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		api.EncodeErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	authUser, err := user.Authenticate(request.Auth.Username, request.Auth.Password)
	if err != nil {
		api.EncodeErrorResponse(w, http.StatusUnauthorized, "Unauthenticated")
		return
	}

	cfg, err := config.GetConfig()
	if err != nil {
		api.EncodeErrorResponse(w, http.StatusBadRequest, "Couldn't read config from context")
		return
	}

	// Determine key usage from type
	var keyUsage x509.KeyUsage
	var extKeyUsage []x509.ExtKeyUsage

	switch request.Data.Type {
	case "client":
		keyUsage = x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature
		extKeyUsage = []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth}
	default:
		keyUsage = x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature
		extKeyUsage = []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth}
	}

	// Decode CSR from base64
	csrBytes, err := base64.StdEncoding.DecodeString(request.Data.CSR)
	if err != nil {
		api.EncodeErrorResponse(w, http.StatusBadRequest, "Invalid CSR format")
		return
	}

	// Parse CSR
	csr, err := x509.ParseCertificateRequest(csrBytes)
	if err != nil {
		api.EncodeErrorResponse(w, http.StatusBadRequest, "Failed to parse CSR")
		return
	}

	serialNumber, err := generateRandomSerialNumber()
	if err != nil {
		api.EncodeErrorResponse(w, http.StatusBadRequest, "Failed to generate serial number")
		return
	}

	if request.Data.Lifetime < 1 {
		request.Data.Lifetime = cfg.CertificateLifetimeDefault
	}

	// Create certificate template
	template := x509.Certificate{
		SerialNumber:          serialNumber,
		Subject:               csr.Subject,
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(0, 0, request.Data.Lifetime),
		KeyUsage:              keyUsage,
		ExtKeyUsage:           extKeyUsage,
		BasicConstraintsValid: true,
	}

	// Generate private key
	priv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		api.EncodeErrorResponse(w, http.StatusInternalServerError, "Failed to generate private key")
		return
	}

	// Generate certificate
	certBytes, err := x509.CreateCertificate(rand.Reader, &template, cfg.CACert, &priv.PublicKey, cfg.CAKey)
	if err != nil {
		api.EncodeErrorResponse(w, http.StatusInternalServerError, "Failed to create certificate")
		return
	}

	// Encode certificate to PEM format
	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certBytes})

	// Save certificate to database
	cert := models.NewCertificate(template.Subject.SerialNumber, certPEM, authUser.Username)

	err = repositories.GetCertificateRepository().Insert(*cert)
	if err != nil {
		api.EncodeErrorResponse(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	// Return certificate to the client
	api.EncodeResponse(w, api.CertificateResponseData{CertificatePEM: string(certPEM)})
}

// HandleCertificateRetrieval retrieves a certificate by serial number.
func HandleCertificateRetrieval(w http.ResponseWriter, r *http.Request) {
	var request api.Request

	_, err := user.Authenticate(request.Auth.Username, request.Auth.Password)
	if err != nil {
		api.EncodeErrorResponse(w, http.StatusUnauthorized, "Unauthenticated")
		return
	}

	queryValues := r.URL.Query()
	certSerialNumber := queryValues.Get("serialNumber")

	if certSerialNumber == "" {
		api.EncodeErrorResponse(w, http.StatusBadRequest, "Missing serialnumber parameter")
		return
	}

	cert, err := repositories.GetCertificateRepository().FindBySerialNumber(certSerialNumber)
	if err != nil {
		api.EncodeErrorResponse(w, http.StatusNotFound, "Certificate not found")
		return
	}

	api.EncodeResponse(w, api.CertificateResponseData{CertificatePEM: string(cert.CertificatePEM)})
}

func generateRandomSerialNumber() (*big.Int, error) {
	maxSerialNumber := new(big.Int).Lsh(big.NewInt(1), 128) // 2^128
	serialNumber, err := rand.Int(rand.Reader, maxSerialNumber)
	if err != nil {
		return nil, err
	}
	return serialNumber, nil
}
