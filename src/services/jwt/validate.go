package jwt

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/orbis-challenge/src/models"
	"github.com/pkg/errors"
)

// Service error list.
var (
	ErrTokenInvalid       = errors.New("jwt is invalid")
	ErrTokenInBlackList   = errors.New("jwt is/was blocked")
	ErrTokenNotFound      = errors.New("jwt not found")
	ErrTokenClaimsInvalid = errors.New("jwt claims is invalid")
)

// Validate validates access jwt.
func (t TokenSrv) Validate(accessToken string) (models.Claims, error) {
	token, err := t.parseJWT(accessToken)
	if err != nil {
		return models.Claims{}, err
	}

	claims := t.parseClaims(token)
	if claims == nil || !token.Valid {
		return models.Claims{}, ErrTokenInvalid
	}

	return *claims, nil
}

func (t TokenSrv) parseJWT(tokenString string) (*jwt.Token, error) {
	token, err := jwt.ParseWithClaims(tokenString, &models.Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(t.secretKey), nil
	})

	if token == nil {
		return nil, ErrTokenInvalid
	}

	switch ve := err.(type) {
	case *jwt.ValidationError:
		if ve.Errors|(jwt.ValidationErrorExpired) != jwt.ValidationErrorExpired {
			return nil, ErrTokenInvalid
		}
	case nil:
	default:
		return nil, err
	}

	return token, nil
}

func (t TokenSrv) parseClaims(token *jwt.Token) *models.Claims {
	if claims, ok := token.Claims.(*models.Claims); ok {
		return claims
	}

	return nil
}
