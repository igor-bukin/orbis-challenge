package config

import (
	"encoding/json"
	"fmt"
	"reflect"

	config "github.com/orbis-challenge/src/config/redact"
)

var (
	redactedFields = []string{"password", "token"}
	redactor       = config.NewRedactor(redactedFields)
)

// aliases used to prevent recursion call of methods
type (
	aliasConfiguration Configuration
	aliasPostgres      Postgres
	aliasRedis         Redis
)

func cloneRedact(dst, src interface{}) error {
	err := config.Clone(dst, src)
	redactor.Redact(reflect.ValueOf(dst))
	return err
}

// MarshalJSON redact sensitive data
func (c *Configuration) MarshalJSON() ([]byte, error) {
	var tmp aliasConfiguration
	if err := cloneRedact(&tmp, c); err != nil {
		return nil, err
	}
	return json.Marshal(tmp)
}

// String - %s, %v, %+v
func (c *Configuration) String() string {
	var tmp aliasConfiguration
	if err := cloneRedact(&tmp, c); err != nil {
		return ""
	}
	return fmt.Sprintf("%+v", tmp)
}

// GoString - %#v
func (c *Configuration) GoString() string {
	var tmp aliasConfiguration
	if err := cloneRedact(&tmp, c); err != nil {
		return ""
	}
	return fmt.Sprintf("%#v", tmp)
}

// MarshalJSON redact sensitive data
func (p *Postgres) MarshalJSON() ([]byte, error) {
	var tmp aliasPostgres
	if err := cloneRedact(&tmp, p); err != nil {
		return nil, err
	}
	return json.Marshal(tmp)
}

// String - %s, %v, %+v
func (p *Postgres) String() string {
	var tmp aliasPostgres
	if err := cloneRedact(&tmp, p); err != nil {
		return ""
	}
	return fmt.Sprintf("%+v", tmp)
}

// GoString - %#v
func (p *Postgres) GoString() string {
	var tmp aliasPostgres
	if err := cloneRedact(&tmp, p); err != nil {
		return ""
	}
	return fmt.Sprintf("%#v", tmp)
}

// MarshalJSON - redis redact sensitive data
func (rds *Redis) MarshalJSON() ([]byte, error) {
	var tmp aliasRedis
	if err := cloneRedact(&tmp, rds); err != nil {
		return nil, err
	}
	return json.Marshal(tmp)
}

// String for redis - %s, %v, %+v
func (rds *Redis) String() string {
	var tmp aliasRedis
	if err := cloneRedact(&tmp, rds); err != nil {
		return ""
	}
	return fmt.Sprintf("%+v", tmp)
}

// GoString for redis - %#v
func (rds *Redis) GoString() string {
	var tmp aliasRedis
	if err := cloneRedact(&tmp, rds); err != nil {
		return ""
	}
	return fmt.Sprintf("%#v", tmp)
}
