package user

import (
	"context"
	"sync"

	"github.com/orbis-challenge/src/models"
	httperrors "github.com/orbis-challenge/src/models/http-errors"
	"github.com/orbis-challenge/src/persistence/postgres"
	"github.com/orbis-challenge/src/services/encryptor"
	"github.com/orbis-challenge/src/services/jwt"
)

type Service interface {
	Login(ctx context.Context, loginReq models.LoginRequest) (models.Token, httperrors.HTTPError)
	SignUp(ctx context.Context, signUpReq models.SignUpRequest) httperrors.HTTPError
}

type service struct {
	pg         *postgres.Client
	encryptSrv encryptor.Service
	tokenSrv   jwt.Service
}

var (
	srv  Service
	once = &sync.Once{}
)

func New(pg *postgres.Client, encryptSrv encryptor.Service, tokenSrv jwt.Service) Service {
	once.Do(func() {
		srv = service{
			pg:         pg,
			encryptSrv: encryptSrv,
			tokenSrv:   tokenSrv,
		}
	})
	return srv
}
