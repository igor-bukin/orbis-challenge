package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/orbis-challenge/src/handlers/common"
	httperrors "github.com/orbis-challenge/src/models/http-errors"
	"github.com/orbis-challenge/src/services"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"go.uber.org/zap"
)

const (
	AccessTokenHeader = "Authorization"
	BearerSchema      = "Bearer"
)

// Auth - authenticate User by JWT token and add to context his ID, REA ID, role etc
func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var (
			authToken = r.Header.Get(AccessTokenHeader)
		)
		ctx := ParseAndProcessAuthToken(w, r, authToken)
		if ctx != nil {
			next.ServeHTTP(w, r.WithContext(ctx))
		}
	})
}

func ParseAndProcessAuthToken(w http.ResponseWriter, r *http.Request, authToken string) context.Context {
	jwtToken, err := ParseAuthorizationHeader(authToken)
	if err != nil {
		logrus.Error("Token is invalid", zap.String("token", authToken))
		common.SendHTTPError(w, httperrors.NewUnauthorizedError(err))
		return nil
	}

	_, err = services.Get().Jwt().Validate(jwtToken)
	if err != nil {
		logrus.Error("Token is invalid", zap.String("token", jwtToken))
		common.SendHTTPError(w, httperrors.NewUnauthorizedError(err))
		return nil
	}

	return context.Background()
}

func ParseAuthorizationHeader(header string) (string, error) {
	authSlice := strings.Split(strings.TrimSpace(header), " ")
	if len(authSlice) != 2 || !strings.EqualFold(authSlice[0], BearerSchema) {
		return "", errors.New("invalid authorization header")
	}

	return authSlice[1], nil
}
