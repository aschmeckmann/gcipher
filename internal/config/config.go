package config

import (
	"context"
	"crypto/x509"
	"fmt"
	"gcipher/internal/util"
	"os"
	"strconv"
	"sync"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Port                       int    `yaml:"port"`
	DatabaseURL                string `yaml:"database_url"`
	CertificateLifetimeDefault int    `yaml:"certificate_lifetime_default"`
	CACertPath                 string `yaml:"ca_cert_path"`
	CAKeyPath                  string `yaml:"ca_key_path"`
	CACert                     *x509.Certificate
	CAKey                      interface{}
}

// Default values
const (
	DefaultPort                       = 8080
	DefaultCertificateLifetimeDefault = 365 // Days
	DefaultCACertPath                 = "ca.crt"
	DefaultCAKeyPath                  = "ca.key"
	DefaultDatabaseURL                = "mongodb://localhost:27017"
)

var (
	configOnce sync.Once
	configKey  = "config"
	cfg        *Config
)

func NewConfig() (*Config, error) {
	var cfg Config

	configFile, err := findConfigFile()
	if err != nil {
		fmt.Println("No config file found, using defaults...")

		cfg.Port = DefaultPort
		cfg.CertificateLifetimeDefault = DefaultCertificateLifetimeDefault
		cfg.CACertPath = DefaultCACertPath
		cfg.CAKeyPath = DefaultCAKeyPath
		cfg.DatabaseURL = DefaultDatabaseURL
	} else {
		content, err := os.ReadFile(configFile)
		if err != nil {
			return nil, fmt.Errorf("failed to read config file: %v", err)
		}

		if err := yaml.Unmarshal(content, &cfg); err != nil {
			return nil, fmt.Errorf("failed to unmarshal YAML: %v", err)
		}
	}

	// Set environment variable overrides
	if portStr := os.Getenv("GCIPHER_PORT"); portStr != "" {
		cfg.Port, err = strconv.Atoi(portStr)
		if err != nil {
			return nil, fmt.Errorf("invalid GCIPHER_PORT value: %s", portStr)
		}
	}

	if dbURL := os.Getenv("GCIPHER_DATABASE_URL"); dbURL != "" {
		cfg.DatabaseURL = dbURL
	}

	if lifetimeStr := os.Getenv("GCIPHER_CERTIFICATE_LIFETIME_DEFAULT"); lifetimeStr != "" {
		cfg.CertificateLifetimeDefault, err = strconv.Atoi(lifetimeStr)
		if err != nil {
			return nil, fmt.Errorf("invalid GCIPHER_CERTIFICATE_LIFETIME_DEFAULT value: %s", lifetimeStr)
		}
	}

	if caCertPath := os.Getenv("GCIPHER_CA_CERT_PATH"); caCertPath != "" {
		cfg.CACertPath = caCertPath
		cert, err := util.ParseCACertificate(caCertPath)
		if err != nil {
			return nil, fmt.Errorf("failed to parse CA certificate: %v", err)
		}
		cfg.CACert = cert
	}

	if caKeyPath := os.Getenv("GCIPHER_CA_KEY_PATH"); caKeyPath != "" {
		cfg.CAKeyPath = caKeyPath
		key, err := util.ParseCAKey(caKeyPath)
		if err != nil {
			return nil, fmt.Errorf("failed to parse CA private key: %v", err)
		}
		cfg.CAKey = key
	}

	return &cfg, nil
}

func GetConfig() (*Config, error) {
	var err error
	configOnce.Do(func() {
		cfg, err = NewConfig()
	})

	if err != nil {
		return nil, err
	}

	return cfg, nil
}

// GetConfigFromContext retrieves the Config object from the context
func GetConfigFromContext(ctx context.Context) (*Config, bool) {
	cfg, ok := ctx.Value(configKey).(*Config)
	return cfg, ok
}

func findConfigFile() (string, error) {
	for _, filename := range []string{"config.yml", "config.yaml"} {
		if _, err := os.Stat(filename); err == nil {
			return filename, nil
		}
	}
	return "", fmt.Errorf("no config file found")
}
