package validator

import (
	"reflect"
	"strings"
	"sync"

	v "github.com/go-playground/validator/v10"
	httperrors "github.com/orbis-challenge/src/models/http-errors"
)

var (
	validatorInstance *v.Validate
	once              sync.Once
)

// Get - initialize once and returns validator instance
func Get() *v.Validate {
	return validatorInstance
}

func FormatErrors(errs v.ValidationErrors) httperrors.HTTPErrors {
	var validationErrs = &httperrors.Errors{
		Errs: make([]httperrors.HTTPError, 0, len(errs)),
	}

	for i := range errs {
		validationErrs.Add(httperrors.Error{
			Field:       errs[i].Field(),
			Description: errs[i].Tag(),
			Code:        httperrors.JSONValidationErr,
		})
	}

	return validationErrs
}

// Load initializes validator
func Load() (err error) {
	once.Do(func() {
		validatorInstance = v.New()
		validatorInstance.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]

			if name == "-" {
				return ""
			}

			return name
		})
	})
	return err
}
