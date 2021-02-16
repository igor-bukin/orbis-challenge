// nolint
package config

import (
	"reflect"
	"strings"
)

const placeholder = "[REDACTED]"

// NewRedactor creates Redactor instance
func NewRedactor(fields []string) Redactor {
	return Redactor{fields: fields}
}

// Redactor redact sensitive data
type Redactor struct {
	fields []string
}

// Redact redact sensitive data
func (r Redactor) Redact(reflectValue reflect.Value) {
	reflectValue = indirect(reflectValue)

	if reflectValue.Kind() == reflect.Struct {
		l := reflectValue.NumField()
		for i := 0; i < l; i++ {
			field := indirect(reflectValue.Field(i))
			fieldKind := field.Kind()

			if fieldKind == reflect.Struct || fieldKind == reflect.Slice {
				r.Redact(field)
				continue
			}

			tt := reflectValue.Type().Field(i)
			if r.isSensitive(tt.Name) && field.CanSet() {
				field.SetString(placeholder)
			}
		}
		return
	}

	if reflectValue.Kind() == reflect.Slice {
		l := reflectValue.Len()
		for i := 0; i < l; i++ {
			r.Redact(reflectValue.Index(i))
		}
		return
	}
}

func (r Redactor) isSensitive(fieldName string) bool {
	fieldName = strings.ToLower(fieldName)
	for _, name := range r.fields {
		if strings.ToLower(name) == fieldName {
			return true
		}
	}
	return false
}
