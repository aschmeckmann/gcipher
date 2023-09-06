package certificate

import (
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
		api.EncodeErrorResponse(w, http.StatusBadRequest, "Couldn't read config")
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

	if request.Data.Lifetime < 1 {
		request.Data.Lifetime = cfg.CertificateLifetimeDefault
	}

	serialNumber, succeed := new(big.Int).SetString(csr.Subject.SerialNumber, 16)
	if !succeed {
		api.EncodeErrorResponse(w, http.StatusBadRequest, "CSR doesn't contain valid serial number")
		return
	}

	// Create certificate template
	template := x509.Certificate{
		Issuer:                cfg.CACert.Subject,
		SignatureAlgorithm:    csr.SignatureAlgorithm,
		PublicKeyAlgorithm:    csr.PublicKeyAlgorithm,
		Version:               csr.Version,
		SerialNumber:          serialNumber,
		Subject:               csr.Subject,
		IPAddresses:           csr.IPAddresses,
		EmailAddresses:        csr.EmailAddresses,
		DNSNames:              csr.DNSNames,
		URIs:                  csr.URIs,
		Extensions:            csr.Extensions,
		ExtraExtensions:       csr.ExtraExtensions,
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(0, 0, request.Data.Lifetime),
		KeyUsage:              keyUsage,
		ExtKeyUsage:           extKeyUsage,
		BasicConstraintsValid: true,
	}

	// Generate certificate
	certBytes, err := x509.CreateCertificate(rand.Reader, &template, cfg.CACert, csr.PublicKey, cfg.CAKey)
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

	if request.Data.SerialNumber == "" {
		api.EncodeErrorResponse(w, http.StatusBadRequest, "Missing serialnumber parameter")
		return
	}

	cert, err := repositories.GetCertificateRepository().FindBySerialNumberAndUsername(request.Data.SerialNumber, authUser.Username)
	if err != nil {
		api.EncodeErrorResponse(w, http.StatusNotFound, "Certificate not found")
		return
	}

	// Return certificate to the client
	api.EncodeResponse(w, api.CertificateResponseData{CertificatePEM: string(cert.CertificatePEM)})
}

// HandleRevokeCertificate revokes a certificate by its serial number.
func HandleRevokeCertificate(w http.ResponseWriter, r *http.Request) {
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

	if request.Data.SerialNumber == "" {
		api.EncodeErrorResponse(w, http.StatusBadRequest, "Missing serialnumber parameter")
		return
	}

	cert, err := repositories.GetCertificateRepository().FindBySerialNumberAndUsername(request.Data.SerialNumber, authUser.Username)
	if err != nil {
		api.EncodeErrorResponse(w, http.StatusNotFound, "Certificate not found")
		return
	}

	// Check if the certificate is already revoked
	if cert.RevokedAt != nil {
		api.EncodeErrorResponse(w, http.StatusBadRequest, "Certificate already revoked")
		return
	}

	// Update certificate with revocation time
	now := time.Now()
	cert.RevokedAt = &now
	err = repositories.GetCertificateRepository().Update(*cert)
	if err != nil {
		api.EncodeErrorResponse(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	// Return success response
	api.EncodeResponse(w, api.Response{Success: true})
}

// HandleCertificateList retrieves a list of certificates based on the specified state filter.
func HandleCertificateList(w http.ResponseWriter, r *http.Request) {
	var request api.Request
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		api.EncodeErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	_, err = user.Authenticate(request.Auth.Username, request.Auth.Password)
	if err != nil {
		api.EncodeErrorResponse(w, http.StatusUnauthorized, "Unauthenticated")
		return
	}

	/* TODO: add admin check?
	if !authUser.IsAdmin() {
		api.EncodeErrorResponse(w, http.StatusForbidden, "Access denied")
		return
	}
	*/

	certificates, err := repositories.GetCertificateRepository().FindByState(request.Data.State)
	if err != nil {
		api.EncodeErrorResponse(w, http.StatusInternalServerError, "Failed to retrieve certificates")
		return
	}

	certList := make([]api.CertificateResponseData, 0, len(certificates))
	for _, cert := range certificates {
		certList = append(certList, api.CertificateResponseData{CertificatePEM: string(cert.CertificatePEM)})
	}

	api.EncodeResponse(w, certList)
}
