package user

import (
	"context"

	"github.com/go-pg/pg"
	"github.com/orbis-challenge/src/models"
	httperrors "github.com/orbis-challenge/src/models/http-errors"
	"github.com/pkg/errors"
)

func (s service) Login(ctx context.Context, loginReq models.LoginRequest) (models.Token, httperrors.HTTPError) {
	tx := s.pg.QueryContext(ctx)

	u, err := tx.GetUserByEmail(loginReq.Email)
	if err != nil {
		if err == pg.ErrNoRows {
			return models.Token{}, httperrors.NewNotFoundError(errors.New("user not found"), "email")
		}
		return models.Token{}, httperrors.NewInternalServerError(err)
	}

	if !s.encryptSrv.CompareHashAndPassword(loginReq.Password, u.Password) {
		return models.Token{}, httperrors.NewErrorByCode(httperrors.InvalidRequest, "")
	}

	token, err := s.createToken()
	if err != nil {
		return models.Token{}, httperrors.NewInternalServerError(errors.Wrap(err, "create token"))
	}

	return token, nil
}

func (s service) createToken() (models.Token, error) {
	claims := models.Claims{}
	jwtToken, err := s.tokenSrv.Generate(&claims)
	if err != nil {
		return models.Token{}, errors.Wrap(err, "generate jwt token")
	}

	return models.Token{
		Access: jwtToken,
	}, nil
}
