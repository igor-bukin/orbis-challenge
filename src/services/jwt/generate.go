package jwt

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/orbis-challenge/src/models"
	"github.com/orbis-challenge/src/utils"
)

func (t TokenSrv) Generate(claims *models.Claims) (string, error) {
	claims.StandardClaims.ExpiresAt = time.Now().Unix() + int64(t.expirationTime)

	tokenWithClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	token, err := tokenWithClaims.SignedString([]byte(t.secretKey))
	if err != nil {
		return "", err
	}

	return token, nil
}

func (t TokenSrv) GenerateRefresh() (string, error) {
	return utils.GenerateRandomString(models.TokenSize)
}
