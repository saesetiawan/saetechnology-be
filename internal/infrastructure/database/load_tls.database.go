package database

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"
)

func loadTLS(path string) (*tls.Config, error) {
	caCert, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	pool := x509.NewCertPool()
	if !pool.AppendCertsFromPEM(caCert) {
		return nil, fmt.Errorf("failed append CA cert")
	}

	return &tls.Config{
		RootCAs: pool,
	}, nil
}
