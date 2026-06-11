package jwt

import (
	"github.com/google/uuid"
	"go-platform-core/internal/config"
	"time"

	"crypto/rsa"
	"github.com/golang-jwt/jwt/v5"
)

type JWT interface {
	Generate(id uuid.UUID) (string, error)
	GenerateAccessToken(id uuid.UUID, role string) (string, error)
	GenerateRefreshToken(id uuid.UUID) (string, error)
	GenerateActivationToken(id uuid.UUID) (string, error)
	Verify(token string) (*jwt.Token, error)
}

type jwtService struct {
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
}

func NewJWT(config *config.Config) JWT {
	return &jwtService{
		privateKey: config.JwtPrivateKey,
		publicKey:  config.JwtPublicKey,
	}
}

func (j *jwtService) Generate(id uuid.UUID) (string, error) {
	return j.GenerateAccessToken(id, "")
}

func (j *jwtService) GenerateAccessToken(id uuid.UUID, role string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"id":         id.String(),
		"role":       role,
		"token_type": "access",
		"exp":        time.Now().Add(15 * time.Minute).Unix(),
	})

	return token.SignedString(j.privateKey)
}

func (j *jwtService) GenerateRefreshToken(id uuid.UUID) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"id":         id.String(),
		"token_type": "refresh",
		"exp":        time.Now().Add(365 * 24 * time.Hour).Unix(),
	})

	return token.SignedString(j.privateKey)
}

func (j *jwtService) GenerateActivationToken(id uuid.UUID) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"id":         id.String(),
		"token_type": "activation",
		"exp":        time.Now().Add(24 * time.Hour).Unix(),
	})

	return token.SignedString(j.privateKey)
}

func (j *jwtService) Verify(tokenStr string) (*jwt.Token, error) {
	return jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return j.publicKey, nil
	})
}
