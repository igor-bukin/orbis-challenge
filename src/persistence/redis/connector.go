package redis

import (
	"sync"

	"github.com/go-redis/redis"
	"github.com/orbis-challenge/src/config"
)

var (
	client *Client
	once   = &sync.Once{}
)

// Load - create redis client
func Load(cfg config.Redis) (err error) {
	once.Do(func() {
		cli := redis.NewClient(
			&redis.Options{
				Password: cfg.Password,
				PoolSize: cfg.PoolSize,
				Addr:     cfg.Addrs,
			})

		err = cli.Ping().Err()
		if err != nil {
			return
		}

		client = &Client{cli: cli}
	})

	return
}

// GetRedis returns redis client
func GetRedis() RedisCli {
	return client
}
