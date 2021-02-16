package jwt

import (
	"time"

	"github.com/orbis-challenge/src/config"
	"github.com/orbis-challenge/src/models"
	"github.com/orbis-challenge/src/persistence/redis"
	"github.com/pkg/errors"
)

// Service jwt interface
type ServiceValidate interface {
	Validate(token string) (models.Claims, error)
}

type ServiceGenerate interface {
	Generate(*models.Claims) (string, error)
	GenerateRefresh() (string, error)
}

type Service interface {
	ServiceValidate
	ServiceGenerate
}

// TokenSrv Service
type TokenSrv struct {
	rds            redis.RedisCli
	secretKey      string
	expirationTime int
}

// New creates new TokenSrv Service
func New(rds redis.RedisCli, cfg *config.Configuration) (Service, error) {
	expTime, err := time.ParseDuration(cfg.Token.ExpirationTime)
	if err != nil {
		return nil, errors.Wrap(err, "parse token expiration time")
	}

	return &TokenSrv{
		rds:            rds,
		secretKey:      cfg.Token.Secret,
		expirationTime: int(expTime.Seconds()),
	}, nil
}
