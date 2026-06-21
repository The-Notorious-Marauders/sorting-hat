package entities

import (
	"github.com/golang-jwt/jwt/v4"
)

const (
	TokenTypeAccess  = "access"
	TokenTypeRefresh = "refresh"
)

type SortingHatClaims struct {
	UserID               string `json:"user_id,omitempty"`
	LastLoginAtInSeconds int    `json:"last_login_at_in_seconds,omitempty"`
	TokenType            string `json:"token_type,omitempty"`

	jwt.RegisteredClaims
}
