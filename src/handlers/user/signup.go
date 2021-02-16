package user

import (
	"net/http"

	"github.com/orbis-challenge/src/handlers/common"
	"github.com/orbis-challenge/src/models"
	"github.com/orbis-challenge/src/services"
	"github.com/sirupsen/logrus"
	"go.uber.org/zap"
)

// SignUp - creating of new user.
func SignUp(w http.ResponseWriter, r *http.Request) {
	var (
		signUpReq models.SignUpRequest
		ctx       = r.Context()
	)

	httpErr := common.UnmarshalRequestBody(r, &signUpReq)
	if httpErr != nil {
		common.SendHTTPError(w, httpErr)
		return
	}

	httpErrs := common.ValidateRequestBodyWithErrors(r, &signUpReq)
	if httpErrs != nil {
		common.SendHTTPErrors(w, http.StatusBadRequest, httpErrs)
		return
	}

	httpErr = services.Get().User().SignUp(ctx, signUpReq)
	if httpErr != nil {
		logrus.Error("signup user error",
			zap.String("email", signUpReq.Email),
			zap.Error(httpErr))
		common.SendHTTPError(w, httpErr)
		return
	}

	common.SendResponse(w, http.StatusNoContent, struct{}{})
}
