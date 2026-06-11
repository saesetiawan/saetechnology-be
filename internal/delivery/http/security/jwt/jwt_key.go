package jwt

import (
	"crypto/rsa"
	"github.com/golang-jwt/jwt/v5"
	"os"
)

func LoadPrivateKey() (*rsa.PrivateKey, error) {
	keyBytes, err := os.ReadFile("private_keys/jwt_private_key.pem")
	if err != nil {
		return nil, err
	}

	return jwt.ParseRSAPrivateKeyFromPEM(keyBytes)
}

func LoadPublicKey() (*rsa.PublicKey, error) {
	keyBytes, err := os.ReadFile("private_keys/jwt_public_key.pem")
	if err != nil {
		return nil, err
	}

	return jwt.ParseRSAPublicKeyFromPEM(keyBytes)
}
