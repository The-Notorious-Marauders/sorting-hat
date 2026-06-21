package usecases

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"time"

	"github.com/The-Notorious-Marauders/sorting-hat/core/entities"
	"github.com/golang-jwt/jwt/v4"
)

const (
	accessTokenTTL  = 15 * time.Minute
	refreshTokenTTL = 7 * 24 * time.Hour
)

func signToken(
	privateKey *ecdsa.PrivateKey,
	userID string,
	lastLoginAtInSeconds int,
	tokenType string,
	issuedAt time.Time,
	ttl time.Duration,
) (string, error) {
	claims := entities.SortingHatClaims{
		UserID:               userID,
		LastLoginAtInSeconds: lastLoginAtInSeconds,
		TokenType:            tokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(issuedAt),
			ExpiresAt: jwt.NewNumericDate(issuedAt.Add(ttl)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	return token.SignedString(privateKey)
}

func parseECDSAPrivateKey(pemKey string) (*ecdsa.PrivateKey, error) {
	block, _ := pem.Decode([]byte(pemKey))
	if block == nil {
		return nil, errors.New("invalid PEM block for ECDSA private key")
	}

	return x509.ParseECPrivateKey(block.Bytes)
}

func parseECDSAPublicKey(pemKey string) (*ecdsa.PublicKey, error) {
	block, _ := pem.Decode([]byte(pemKey))
	if block == nil {
		return nil, errors.New("invalid PEM block for ECDSA public key")
	}

	key, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	publicKey, ok := key.(*ecdsa.PublicKey)
	if !ok {
		return nil, errors.New("PEM block does not contain an ECDSA public key")
	}

	return publicKey, nil
}
