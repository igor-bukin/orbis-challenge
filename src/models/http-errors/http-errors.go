package httperrors

import (
	"fmt"
	"net/http"
	"strings"
)

type HTTPError interface {
	error

	HTTPCode() int
	ErrField() string
	ErrCode() ErrorCode
	Instance() Error
}

type HTTPErrors interface {
	error

	Add(...HTTPError)
	Errors() []HTTPError
	Instance() Errors
	Len() int
}

func NewHTTPErrors() HTTPErrors {
	return &Errors{Errs: make([]HTTPError, 0, 1)}
}

// Error - protocol defined structure for Server Error description
type Error struct {
	Code        ErrorCode `json:"code"`
	Description string    `json:"description,omitempty"`
	Field       string    `json:"field,omitempty"`
}

func (err Error) ErrCode() ErrorCode {
	return err.Code
}

func (err Error) ErrField() string {
	return err.Field
}

func (err Error) HTTPCode() int {
	return ErrorHTTPCodes[err.Code]
}

func (err Error) Instance() Error {
	return err
}

// Error unifying models.Error with Go' error interface
func (err Error) Error() string {
	return err.Description
}

type Errors struct {
	Errs []HTTPError `json:"errors"`
}

func (errs *Errors) Add(err ...HTTPError) { // nolint
	errs.Errs = append(errs.Errs, err...)
}

func (err Errors) Len() int { // nolint
	return len(err.Errs)
}

func (errs Errors) Instance() Errors {
	return errs
}

func (errs Errors) Error() string {
	var errors = make([]string, len(errs.Errs))

	for i := range errs.Errs {
		errors[i] = errs.Errs[i].Error()
	}

	return strings.Join(errors, " ")
}

func (errs Errors) Errors() []HTTPError {
	return errs.Errs
}

// ErrorCode error models codes type
type ErrorCode int

const (
	JSONValidationErr      ErrorCode = 4000
	JSONReadErr            ErrorCode = 4001
	JSONParseErr           ErrorCode = 4002
	ForbiddenErr           ErrorCode = 4003
	RecordNotFoundErr      ErrorCode = 4004
	RecordAlreadyExistsErr ErrorCode = 4005
	UnauthorizedErr        ErrorCode = 4010
	UUIDParseErr           ErrorCode = 4026
	InternalServerErr      ErrorCode = 5000
	InvalidRequest         ErrorCode = 6000
)

// ErrorDescriptions - error codes and description for UI.
var ErrorDescriptions = map[ErrorCode]string{
	JSONValidationErr:      "Invalid request. JSON field validation error.",
	JSONReadErr:            "Invalid request. Can't read request body.",
	JSONParseErr:           "Invalid request. Can't parse request body.",
	ForbiddenErr:           "Action is not authorized.",
	RecordNotFoundErr:      "Record not found: ",
	RecordAlreadyExistsErr: "Record already exists: ",
	UnauthorizedErr:        "Unauthorized: ",
	InternalServerErr:      "Internal Server Error: ",
	InvalidRequest:         "Invalid request: ",
	UUIDParseErr:           "Failed to parse UUID: ",
}

// ErrorHTTPCodes represents HTTP errors corresponding to ErrorCode
var ErrorHTTPCodes = map[ErrorCode]int{
	JSONValidationErr: http.StatusBadRequest,
	JSONReadErr:       http.StatusBadRequest,
	JSONParseErr:      http.StatusBadRequest,
	InvalidRequest:    http.StatusBadRequest,
	UUIDParseErr:      http.StatusBadRequest,
	ForbiddenErr:      http.StatusForbidden,
	UnauthorizedErr:   http.StatusUnauthorized,

	RecordNotFoundErr:      http.StatusNotFound,
	RecordAlreadyExistsErr: http.StatusBadRequest,

	InternalServerErr: http.StatusInternalServerError,
}

// NewError returns not nil error wrapped with models.Error
func NewError(err error, code ErrorCode, field string) HTTPError {
	if err == nil {
		return nil
	}

	return &Error{
		Code:        code,
		Description: fmt.Sprint(ErrorDescriptions[code], err),
		Field:       field,
	}
}

// NewErrorByCode returns new error by provided `code`
func NewErrorByCode(code ErrorCode, field string) HTTPError {
	return &Error{
		Code:        code,
		Description: ErrorDescriptions[code],
		Field:       field,
	}
}

func NewBadRequestError(err error) HTTPError {
	if err == nil {
		return nil
	}

	return &Error{
		Code:        InvalidRequest,
		Description: fmt.Sprint(ErrorDescriptions[InvalidRequest], err.Error()),
	}
}

func NewUnauthorizedError(err error) HTTPError {
	if err == nil {
		return nil
	}

	return &Error{
		Code:        UnauthorizedErr,
		Description: fmt.Sprint(ErrorDescriptions[UnauthorizedErr], err.Error()),
	}
}

func NewForbiddenError() HTTPError {
	return &Error{
		Code:        ForbiddenErr,
		Description: ErrorDescriptions[ForbiddenErr],
	}
}

func NewNotFoundError(err error, field string) HTTPError {
	if err == nil {
		return nil
	}

	return &Error{
		Code:        RecordNotFoundErr,
		Description: fmt.Sprint(ErrorDescriptions[RecordNotFoundErr], err.Error()),
		Field:       field,
	}
}

func NewInternalServerError(err error) HTTPError {
	if err == nil {
		return nil
	}

	return &Error{
		Code:        InternalServerErr,
		Description: fmt.Sprint(ErrorDescriptions[InternalServerErr], err.Error()),
	}
}

func NewAlreadyExistsError(err error, field string) HTTPError {
	if err == nil {
		return nil
	}

	return &Error{
		Code:        RecordAlreadyExistsErr,
		Description: fmt.Sprint(ErrorDescriptions[RecordAlreadyExistsErr], err.Error()),
		Field:       field,
	}
}
