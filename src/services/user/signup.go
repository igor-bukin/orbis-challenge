package user

import (
	"context"

	"github.com/go-pg/pg"
	"github.com/orbis-challenge/src/models"
	httperrors "github.com/orbis-challenge/src/models/http-errors"
	"github.com/pkg/errors"
)

func (s service) SignUp(ctx context.Context, signUpReq models.SignUpRequest) httperrors.HTTPError {
	tx, err := s.pg.NewTXContext(ctx)
	if err != nil {
		return httperrors.NewInternalServerError(errors.Wrap(err, "start transaction"))
	}

	defer tx.Rollback() // nolint

	password, encErr := s.encryptSrv.EncryptPassword(signUpReq.Password)
	if encErr != nil {
		return httperrors.NewInternalServerError(err)
	}

	user := models.User{
		Email:    signUpReq.Email,
		Password: password,
	}
	_, err = tx.CreateUser(&user)
	if err != nil {
		pgErr, ok := err.(pg.Error)
		if ok && pgErr.IntegrityViolation() {
			return httperrors.NewAlreadyExistsError(err, "email")
		}
		return httperrors.NewInternalServerError(err)
	}

	if err := tx.Commit(); err != nil {
		return httperrors.NewInternalServerError(err)
	}

	return nil
}
