package common

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"

	v "github.com/go-playground/validator/v10"
	"github.com/orbis-challenge/src/config"
	httperrors "github.com/orbis-challenge/src/models/http-errors"
	"github.com/orbis-challenge/src/validator"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"go.uber.org/zap"
)

// SendResponse - common method for encoding and writing any json response
func SendResponse(w http.ResponseWriter, statusCode int, respBody interface{}) {
	w.Header().Set("Content-Type", "application/json;charset=utf-8")

	binRespBody, err := json.Marshal(respBody)
	if err != nil {
		statusCode = http.StatusInternalServerError
	}

	w.WriteHeader(statusCode)
	//nolint
	w.Write(binRespBody)
}

// SendHTTPError sends HTTP error with code selected based on the `httpErr` ErrorCode
func SendHTTPError(w http.ResponseWriter, httpErr httperrors.HTTPError) {
	err := httperrors.Errors{Errs: []httperrors.HTTPError{httpErr.Instance()}}
	SendResponse(w, httpErr.HTTPCode(), err)
}

// SendHTTPErrors sends HTTP errors
func SendHTTPErrors(w http.ResponseWriter, code int, httpErrs httperrors.HTTPErrors) {
	SendResponse(w, code, httpErrs.Instance())
}

// HandleRequest reads all bytes from request and closes it to prevent memory leaks
func HandleRequest(r *http.Request) (data []byte, err error) {
	data, err = ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	return data, r.Body.Close()
}

func UnmarshalRequestBody(r *http.Request, body interface{}) httperrors.HTTPError {
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		serverError := httperrors.Error{
			Code:        httperrors.JSONReadErr,
			Description: httperrors.ErrorDescriptions[httperrors.JSONReadErr],
		}

		logrus.Error("Invalid request body", zap.Error(err))
		return serverError
	}
	defer r.Body.Close()

	err = json.Unmarshal(reqBody, body)
	if err != nil {
		serverError := httperrors.Error{
			Code:        httperrors.JSONParseErr,
			Description: httperrors.ErrorDescriptions[httperrors.JSONParseErr],
		}

		logrus.Error("Invalid JSON request body",
			zap.Error(err), zap.String("Corrupted JSON", string(reqBody)))

		return serverError
	}

	return nil
}

// ValidateRequestBody validates object with bad request error.
func ValidateRequestBody(r *http.Request, body interface{}) httperrors.HTTPError {
	err := validator.Get().Struct(body)
	if err != nil {
		serverError := httperrors.NewBadRequestError(err)

		logrus.Info("request body validation failed", zap.Error(err))
		return serverError
	}

	return nil
}

// ValidateRequestBodyWithErrors validates object with errors description list.
func ValidateRequestBodyWithErrors(r *http.Request, body interface{}) httperrors.HTTPErrors {
	err := validator.Get().Struct(body)
	if err != nil {
		validationErrors := err.(v.ValidationErrors)
		serverError := validator.FormatErrors(validationErrors)

		logrus.Info("request body validation failed", zap.Error(err))
		return serverError
	}

	return nil
}

// ProcessRequestBody - read and parse request body.
func ProcessRequestBody(w http.ResponseWriter, r *http.Request, body interface{}) error {
	httpErr := UnmarshalRequestBody(r, body)
	if httpErr != nil {
		SendHTTPError(w, httpErr)
		return httpErr
	}

	httpErr = ValidateRequestBody(r, body)
	if httpErr != nil {
		SendHTTPError(w, httpErr)
		return httpErr
	}

	return nil
}

// ProcessRequestBodyWithErrors - read and parse request body with errors description list.
func ProcessRequestBodyWithErrors(w http.ResponseWriter, r *http.Request, body interface{}) error {
	httpErr := UnmarshalRequestBody(r, body)
	if httpErr != nil {
		SendHTTPError(w, httpErr)
		return httpErr
	}

	httpErrs := ValidateRequestBodyWithErrors(r, body)
	if httpErrs != nil {
		SendHTTPErrors(w, http.StatusUnprocessableEntity, httpErrs)
		return httpErrs
	}

	return nil
}

// GetLimitAndOffset returns limit and offset from query parameters
func GetLimitAndOffset(urlQuery url.Values) (_, _ int, httpErr httperrors.HTTPError) {
	limitStr := urlQuery.Get("limit")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		httpErr = httperrors.NewBadRequestError(errors.Wrap(err, "failed to parse limit from url"))
		return
	}

	if limit <= 0 {
		httpErr = httperrors.NewBadRequestError(errors.New("'limit' must be greater than 0"))
		return
	}

	if limit > config.Config.PaginationLimit {
		httpErr = httperrors.NewBadRequestError(fmt.Errorf("'limit' must be less than %d",
			config.Config.PaginationLimit))
		return
	}

	offsetStr := urlQuery.Get("offset")
	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		httpErr = httperrors.NewBadRequestError(errors.Wrap(err, "failed to parse offset from url"))
		return
	}

	if offset < 0 {
		httpErr = httperrors.NewBadRequestError(errors.New("'offset' must be greater or equal to 0"))
		return
	}

	return limit, offset, nil
}
