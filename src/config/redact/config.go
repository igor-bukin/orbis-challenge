// nolint
package config

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"syscall"
	"time"
)

var (
	unixEpochTime      = time.Unix(0, 0)
	timeType           = reflect.TypeOf((*time.Time)(nil)).Elem()
	durationType       = reflect.TypeOf((*time.Duration)(nil)).Elem()
	durationCustomType = reflect.TypeOf((*Duration)(nil)).Elem()
)

// Load reads and loads configuration to `config` variable, which must be a reference of struct
func Load(config interface{}, filename string) error {
	v := reflect.ValueOf(config).Elem()
	if err := applyDefaults(reflect.StructField{}, v); err != nil {
		return fmt.Errorf("init config with default values: %s", err)
	}

	if err := mergeJSONConfig(config, filename); err != nil {
		return err
	}

	if err := applyEnv(config); err != nil {
		return err
	}

	return validate(config)
}

func validate(config interface{}) error {
	t := reflect.TypeOf(config).Elem()
	v := reflect.ValueOf(config).Elem()
	invalidFields := make([]string, 0)

	for i := 0; i < v.NumField(); i++ {
		invalidFields = append(invalidFields, validateField(t.Field(i), v.Field(i))...)
	}

	if len(invalidFields) > 0 {
		return fmt.Errorf("required fields: %v are not filled up. Please check configuration", invalidFields)
	}

	return nil
}

func validateField(t reflect.StructField, v reflect.Value) (invalidFields []string) {
	if v.Kind() == reflect.Struct && !isTime(v.Type()) {
		for i := 0; i < v.NumField(); i++ {
			invalidFields = append(invalidFields, validateField(v.Type().Field(i), v.Field(i))...)
		}
		return
	}

	value, hasRequiredTag := t.Tag.Lookup("required")
	if !hasRequiredTag || !isTrue(value) {
		return
	}

	if isEmpty(v) {
		invalidFields = append(invalidFields, t.Name)
	}

	return
}

func applyDefaults(t reflect.StructField, v reflect.Value) error {
	if v.Kind() == reflect.Struct && !isTime(v.Type()) {
		for i := 0; i < v.NumField(); i++ {
			if err := applyDefaults(v.Type().Field(i), v.Field(i)); err != nil {
				return err
			}
		}
		return nil
	}

	value, ok := t.Tag.Lookup("default")
	if !ok {
		return nil
	}

	return SetValue(v, value)
}

// SetValue sets value depend on type
func SetValue(v reflect.Value, value string) error {
	switch indirectType(v.Type()) {
	case timeType:
		return setTime(&v, value)
	case durationType:
		return setDuration(&v, value)
	case durationCustomType:
		return setDurationCustom(&v, value)
	}

	switch v.Kind() {
	case reflect.String:
		v.SetString(value)
	case reflect.Slice:
		setSlice(&v, value)
	case reflect.Int, reflect.Int32, reflect.Int64:
		return setInt(&v, value)
	case reflect.Uint, reflect.Uint32, reflect.Uint64:
		return setUint(&v, value)
	case reflect.Bool:
		setBool(&v, value)
	}

	return nil
}

func mergeJSONConfig(config interface{}, filename string) error {
	if len(filename) == 0 {
		return nil
	}

	configFileContents, err := ioutil.ReadFile(filepath.Clean(filename))
	if err != nil {
		log.Printf("Reading configuration from JSON (%s) failed (err: %v). SKIPPED.\n", filename, err)
		return nil
	}
	reader := bytes.NewBuffer(configFileContents)
	return json.NewDecoder(reader).Decode(config)
}

func applyEnv(config interface{}) error {
	v := reflect.ValueOf(config).Elem()

	var errs []error

	for i := 0; i < v.NumField(); i++ {
		err := applyEnvValues(v.Type().Field(i), v.Field(i))
		if err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("mergeEnvConfig: %v", errs)
	}

	return nil
}

func applyEnvValues(t reflect.StructField, v reflect.Value) error {
	if eo, ok := getAddr(v).Interface().(EnvOverrider); ok {
		return eo.ApplyEnvOverrides()
	}

	if v.Kind() == reflect.Struct && !isTime(v.Type()) {
		for i := 0; i < v.NumField(); i++ {
			err := applyEnvValues(v.Type().Field(i), v.Field(i))
			if err != nil {
				return err
			}
		}
		return nil
	}

	envKey, ok := t.Tag.Lookup("envconfig")
	if !ok {
		return nil
	}

	value, found := syscall.Getenv(envKey)
	if !found {
		return nil
	}

	return SetValue(v, value)
}

func setTime(v *reflect.Value, value string) error {
	date, err := time.Parse(time.RFC3339, value)
	if err != nil {
		return fmt.Errorf("can not parse value %q as time.Time type", value)
	}
	if v.Kind() == reflect.Ptr {
		v.Set(reflect.ValueOf(&date))
	} else {
		v.Set(reflect.ValueOf(date))
	}
	return nil
}

func setDuration(v *reflect.Value, value string) error {
	duration, err := time.ParseDuration(value)
	if err != nil {
		return fmt.Errorf("can not parse value %q as time.Duration type", value)
	}
	if v.Kind() == reflect.Ptr {
		v.Set(reflect.ValueOf(&duration))
	} else {
		v.Set(reflect.ValueOf(duration))
	}
	return nil
}

func setDurationCustom(v *reflect.Value, value string) error {
	duration, err := time.ParseDuration(value)
	if err != nil {
		return fmt.Errorf("can not parse value %q as Duration type", value)
	}
	if v.Kind() == reflect.Ptr {
		d := Duration(duration)
		v.Set(reflect.ValueOf(&d))
	} else {
		v.Set(reflect.ValueOf(Duration(duration)))
	}
	return nil
}

func setInt(v *reflect.Value, value string) error {
	intValue, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return fmt.Errorf("can not parse value %q as Int64 type", value)
	}
	v.SetInt(intValue)
	return nil
}

func setUint(v *reflect.Value, value string) error {
	uintValue, err := strconv.ParseUint(value, 10, 64)
	if err != nil {
		return fmt.Errorf("can not parse value %q as Uint64 type", value)
	}
	v.SetUint(uintValue)
	return nil
}

func setBool(v *reflect.Value, value string) {
	v.SetBool(isTrue(value))
}

func setSlice(v *reflect.Value, value string) {
	if _, ok := v.Interface().([]string); !ok {
		return
	}

	values := strings.Split(value, ",")
	slice := reflect.MakeSlice(reflect.TypeOf([]string{}), len(values), len(values))

	for i, value := range values {
		slice.Index(i).SetString(value)
	}

	v.Set(slice)
}

func isTrue(value string) bool {
	return value == "1" || strings.ToLower(value) == "true"
}

func isTime(t reflect.Type) bool {
	return indirectType(t) == timeType
}

func isZeroTime(date time.Time) bool {
	return date.IsZero() || date.Equal(unixEpochTime)
}

func isEmpty(v reflect.Value) bool {
	switch v.Type() {
	case timeType:
		return isZeroTime(v.Interface().(time.Time))
	case durationType:
		return v.Interface().(time.Duration).Nanoseconds() == 0
	case durationCustomType:
		return time.Duration(v.Interface().(Duration)).Nanoseconds() == 0
	}

	switch v.Kind() {
	case reflect.Slice, reflect.String:
		return v.Len() == 0
	case reflect.Int, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint32, reflect.Uint64:
		return v.Interface() == 0
	case reflect.Ptr:
		if v.IsNil() {
			return true
		}
		return isEmpty(v.Elem())
	}
	zero := reflect.Zero(v.Type())
	return reflect.DeepEqual(v.Interface(), zero.Interface())
}

func getAddr(v reflect.Value) reflect.Value {
	if v.Kind() != reflect.Ptr && v.CanAddr() {
		return v.Addr()
	}
	return v
}

func indirect(v reflect.Value) reflect.Value {
	for v.Kind() == reflect.Ptr {
		return v.Elem()
	}
	return v
}

func indirectType(v reflect.Type) reflect.Type {
	for v.Kind() == reflect.Ptr || v.Kind() == reflect.Slice {
		return v.Elem()
	}
	return v
}
