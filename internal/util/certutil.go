package util

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
)

// ParseCACertificate parses the CA certificate from a file and returns the x509.Certificate object
func ParseCACertificate(certPath string) (*x509.Certificate, error) {
	certBytes, err := os.ReadFile(certPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read CA certificate file: %v", err)
	}

	block, _ := pem.Decode(certBytes)
	if block == nil {
		return nil, fmt.Errorf("failed to decode PEM block containing CA certificate")
	}

	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse CA certificate: %v", err)
	}

	return cert, nil
}

// ParseCAKey parses the CA private key from a file and returns the private key object
func ParseCAKey(keyPath string) (interface{}, error) {
	keyBytes, err := os.ReadFile(keyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read CA key file: %v", err)
	}

	block, _ := pem.Decode(keyBytes)
	if block == nil {
		return nil, fmt.Errorf("failed to decode PEM block containing CA key")
	}

	var privateKey interface{}
	switch block.Type {
	case "RSA PRIVATE KEY":
		rsaPrivateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
		if err != nil {
			return nil, fmt.Errorf("failed to parse RSA private key: %v", err)
		}
		privateKey = rsaPrivateKey

	case "EC PRIVATE KEY":
		ecPrivateKey, err := x509.ParseECPrivateKey(block.Bytes)
		if err != nil {
			return nil, fmt.Errorf("failed to parse EC private key: %v", err)
		}
		privateKey = ecPrivateKey

	default:
		return nil, fmt.Errorf("unsupported private key type: %s", block.Type)
	}

	return privateKey, nil
}
