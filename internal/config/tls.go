package config

import (
	"crypto/tls"
	"fmt"

	"golang.org/x/crypto/acme/autocert"
	"google.golang.org/grpc/credentials"
)

type TLS struct {
	CertFilePath string `env:"GOPHKEEPER_TLS_CERT_FILE_PATH" `
	KeyFilePath  string `env:"GOPHKEEPER_TLS_KEY_FILE_PATH" `
	AutoCert     bool   `env:"GOPHKEEPER_TLS_AUTOCERT" env-default:"false"`
	SkipVerify   bool   `env:"GOPHKEEPER_TLS_SKIP_VERIFY"  env-default:"true"` // true for development
}

func LoadTLSCredentials(cfg TLS) (credentials.TransportCredentials, error) {
	if cfg.AutoCert {
		manager := autocert.Manager{
			Prompt: autocert.AcceptTOS,
			Cache:  autocert.DirCache("golang-autocert"),
		}

		config := manager.TLSConfig()
		//nolint:gosec // ignore
		config.InsecureSkipVerify = cfg.SkipVerify

		return credentials.NewTLS(config), nil
	}

	// Load server's certificate and private key
	serverCert, err := tls.LoadX509KeyPair(cfg.CertFilePath, cfg.KeyFilePath)
	if err != nil {
		return nil, fmt.Errorf("could not load server certificate: %w", err)
	}

	// Create the credentials and return it
	config := &tls.Config{
		Certificates: []tls.Certificate{serverCert},
		ClientAuth:   tls.NoClientCert,
		//nolint:gosec // ignore
		InsecureSkipVerify: cfg.SkipVerify,
	}

	return credentials.NewTLS(config), nil
}
