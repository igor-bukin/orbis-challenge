package models

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

const TokenSize = 32

// Token describes access and refresh tokens.
type Token struct {
	Access string `json:"access_token"  validate:"required"`
}

// Claims jwt's claims/payload
type Claims struct {
	jwt.StandardClaims
}

// TTL returns TTL in seconds
func (c *Claims) TTL() int64 {
	return c.StandardClaims.ExpiresAt - time.Now().Unix()
}
