package user

import (
	"net/http"

	"github.com/orbis-challenge/src/handlers/common"
	"github.com/orbis-challenge/src/models"
	"github.com/orbis-challenge/src/services"
	"github.com/sirupsen/logrus"
	"go.uber.org/zap"
)

// Login - authenticates user.
func Login(w http.ResponseWriter, r *http.Request) {
	var (
		loginRequest models.LoginRequest
		ctx          = r.Context()
	)

	httpErr := common.UnmarshalRequestBody(r, &loginRequest)
	if httpErr != nil {
		common.SendHTTPError(w, httpErr)
		return
	}

	httpErrs := common.ValidateRequestBodyWithErrors(r, &loginRequest)
	if httpErrs != nil {
		common.SendHTTPErrors(w, http.StatusBadRequest, httpErrs)
		return
	}

	token, httpErr := services.Get().User().Login(ctx, loginRequest)
	if httpErr != nil {
		logrus.Error("login user error",
			zap.String("email", loginRequest.Email),
			zap.Error(httpErr))
		common.SendHTTPError(w, httpErr)
		return
	}

	common.SendResponse(w, http.StatusOK, &token)
}
