package config

// EnvOverrider overrides own values with ENV values
type EnvOverrider interface {
	ApplyEnvOverrides() error
}
