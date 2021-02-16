package services

import (
	"github.com/orbis-challenge/src/config"
	"github.com/orbis-challenge/src/persistence/postgres"
	"github.com/orbis-challenge/src/persistence/redis"
	"github.com/orbis-challenge/src/services/encryptor"
	"github.com/orbis-challenge/src/services/holding"
	"github.com/orbis-challenge/src/services/jwt"
	"github.com/orbis-challenge/src/services/user"
)

// Service provides concrete services
type Service interface {
	Jwt() jwt.Service
	User() user.Service
	Holding() holding.Service
}

// serviceRepo combines all services
type serviceRepo struct {
	jwtSrv     jwt.Service
	userSrv    user.Service
	encryptSrv encryptor.Service
	holdingSrv holding.Service
}

var srv serviceRepo

// Load initializes services.
func Load(redisCli redis.RedisCli, postgresCli *postgres.Client, cfg *config.Configuration) error {
	jwtSrv, err := jwt.New(redisCli, cfg)
	if err != nil {
		return err
	}

	holdingSrv := holding.New(postgresCli, &cfg.SSGA)
	encryptSrv := encryptor.New()
	userSrv := user.New(postgresCli, encryptSrv, jwtSrv)

	srv.jwtSrv = jwtSrv
	srv.userSrv = userSrv
	srv.encryptSrv = encryptSrv
	srv.holdingSrv = holdingSrv

	return nil
}

// Get returns service repository
func Get() Service {
	return &srv
}

func (s serviceRepo) Jwt() jwt.Service {
	return s.jwtSrv
}

func (s serviceRepo) User() user.Service {
	return s.userSrv
}

func (s serviceRepo) Holding() holding.Service {
	return s.holdingSrv
}
