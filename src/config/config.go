package config

import (
	cfg "github.com/orbis-challenge/src/config/redact"

	"log"
)

var (
	// Config is a package variable, which is populated during init() execution and shared to whole application
	Config Configuration
)

// Load initialize config
func Load(configPath string) error {
	err := cfg.Load(&Config, configPath)
	log.Println(Config.String())
	return err
}

type (
	// Configuration options
	Configuration struct {
		AppName      string   `json:"app_name"      envconfig:"APP_NAME"   default:"orbis-challenge"`
		LogPreset    string   `json:"log_preset"    envconfig:"LOG_PRESET" default:"dev"`
		ListenURL    string   `json:"listen_url"    envconfig:"LISTEN_URL" default:":8888"`
		URLPrefix    string   `json:"url_prefix"    envconfig:"URL_PREFIX" default:"/api"`
		Token        Token    `json:"token"`
		Postgres     Postgres `json:"postgres"`
		PostgresTest Postgres `json:"postgres_test"`
		Redis        Redis    `json:"redis"`
		SSGA         SSGA     `json:"ssga"`

		PaginationLimit int    `json:"pagination_limit" envconfig:"PAGINATION_LIMIT" default:"100"`
		HTTPTimeout     string `json:"http_timeout"     envconfig:"HTTP_TIMEOUT"     default:"10s"`
	}

	// Postgres options
	Postgres struct {
		Host         string `json:"host"             envconfig:"POSTGRES_HOST"           default:"localhost"`
		Port         string `json:"port"             envconfig:"POSTGRES_PORT"           default:"5432"`
		DBName       string `json:"db_name"          envconfig:"POSTGRES_DB"             default:"orbis-api-db"`
		User         string `json:"user"             envconfig:"POSTGRES_USER"           default:"postgres"`
		Password     string `json:"password"         envconfig:"POSTGRES_PASSWORD"       default:"12345"`
		PoolSize     int    `json:"pool_size"        envconfig:"POSTGRES_POOL_SIZE"      default:"10"`
		MaxRetries   int    `json:"max_retries"      envconfig:"POSTGRES_MAX_RETRIES"    default:"5"`
		ReadTimeout  string `json:"read_timeout"     envconfig:"POSTGRES_READ_TIMEOUT"   default:"10s"`
		WriteTimeout string `json:"write_timeout"    envconfig:"POSTGRES_WRITE_TIMEOUT"  default:"10s"`
	}

	// Redis defines configs for Redis Cache
	Redis struct {
		Addrs    string `json:"addrs"     envconfig:"REDIS_ADDRS"`
		PoolSize int    `json:"pool_size" envconfig:"REDIS_POOL_SIZE" default:"10"`
		Password string `json:"password"  envconfig:"REDIS_PASSWORD"  default:""`
	}

	// Token - JWT Token pair settings
	Token struct {
		Secret         string `json:"secret" envconfig:"TOKEN_SECRET" default:"qwerty123456789"`
		ExpirationTime string `json:"expiration_time" envconfig:"TOKEN_EXPIRATION_TIME" default:"24h"`
	}

	SSGA struct {
		URL string `json:"ssga_url" envconfig:"SSGA_URL" default:"https://www.ssga.com/bin/v1/ssmp/fund/fundfinder?country=us&language=en&role=individual&product=etfs&ui=fund-finder"` //nolint:lll
	}
)
