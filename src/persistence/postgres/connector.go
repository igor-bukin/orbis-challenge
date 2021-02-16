package postgres

import (
	"sync"
	"time"

	"github.com/go-pg/pg"
	"github.com/orbis-challenge/src/config"
	"github.com/pkg/errors"
)

// Client postgres connection client
type Client struct {
	conn *pg.DB
}

var (
	postgres *Client
	once     = &sync.Once{}
)

// Load - loads postgres client instance.
func Load(cfg *config.Postgres, lgr DebugLogger) error {
	writeTimeout, err := time.ParseDuration(cfg.WriteTimeout)
	if err != nil {
		return errors.Wrap(err, "parse write timeout")
	}

	readTimeout, err := time.ParseDuration(cfg.ReadTimeout)
	if err != nil {
		return errors.Wrap(err, "parse read timeout")
	}

	once.Do(func() {
		db := pg.Connect(&pg.Options{
			Addr:         cfg.Host + ":" + cfg.Port,
			User:         cfg.User,
			Password:     cfg.Password,
			Database:     cfg.DBName,
			PoolSize:     cfg.PoolSize,
			WriteTimeout: writeTimeout,
			ReadTimeout:  readTimeout,
			MaxRetries:   cfg.MaxRetries,
		})
		db.AddQueryHook(dbLogger{
			logger: lgr,
		})

		postgres = &Client{
			conn: db,
		}
	})

	return postgres.Ping()
}

// GetDB returns postgres DBName client
func GetDB() *Client {
	return postgres
}

// Ping checks db connection.
func (p *Client) Ping() (err error) {
	var n int
	_, err = p.conn.QueryOne(pg.Scan(&n), "SELECT 1")
	return
}
