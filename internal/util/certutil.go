package util

import (
	"bytes"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/youmark/pkcs8"
)

// ParseCertificate parses the CA certificate from a file and returns the x509.Certificate object
func ParseCertificate(certPath string) (*x509.Certificate, error) {
	certBytes, err := os.ReadFile(certPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read CA certificate file: %v", err)
	}

	return ParseCertificateFromBytes(certBytes)
}

// ParseKey parses the CA private key from a file and returns the private key object
func ParseKey(keyPath string, passphrase []byte) (interface{}, error) {
	keyBytes, err := os.ReadFile(keyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read CA key file: %v", err)
	}

	return ParseKeyFromBytes(keyBytes, passphrase)
}

// LoadKeyFromS3 loads a key from an S3 bucket and returns the key bytes
func LoadKeyFromS3(bucket, key, region string) ([]byte, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
	})
	if err != nil {
		return nil, err
	}

	svc := s3.New(sess)

	input := &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}

	result, err := svc.GetObject(input)
	if err != nil {
		return nil, err
	}
	defer result.Body.Close()

	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(result.Body)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// ParseCertificateFromBytes parses the CA certificate from a byte slice and returns the x509.Certificate object
func ParseCertificateFromBytes(certBytes []byte) (*x509.Certificate, error) {
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

// ParseKeyFromBytes parses the CA private key from bytes and returns the crypto.Signer interface
func ParseKeyFromBytes(der []byte, passphrase []byte) (interface{}, error) {
	block, _ := pem.Decode(der)

	if block == nil {
		return nil, fmt.Errorf("failed to decode PEM block")
	}

	switch block.Type {
	case "RSA PRIVATE KEY":
		return x509.ParsePKCS1PrivateKey(block.Bytes)
	case "PRIVATE KEY":
		key, _, err := pkcs8.ParsePrivateKey(der, passphrase)
		return key, err
	case "EC PRIVATE KEY":
		return x509.ParseECPrivateKey(block.Bytes)
	default:
		return nil, fmt.Errorf("unsupported private key type: %s", block.Type)
	}
}
