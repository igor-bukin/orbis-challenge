package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis"
)

const (
	blackListPrefix = "inv_token_id"
)

// RedisCli - interface for different Redis client implementation
// nolint:golint
type RedisCli interface {
	SetInt(key string, value int, ttl int64) error
	GetInt(key string) (int, error)
	Get(key string) ([]byte, error)
	GetString(key string) string
	Set(key string, value []byte, ttl int64) error
	Del(keys ...string) error
	HMSet(key string, fields map[string]interface{}) error
	HMGet(key string, fields ...string) []interface{}
	HGetAll(key string) map[string]string
	HIncrBy(key, field string, incr int64) error
	Expire(key string, exp time.Duration) error
	Ping(ctx context.Context) error
	TTL(key string) (time.Duration, error)
	Incr(key string) error
	Exists(key string) (bool, error)
	GetCli() *redis.Client
	Subscribe(channels ...string) *redis.PubSub
	Publish(channel string, message interface{}) error
	SIsMember(key string, member interface{}) (bool, error)
	SRem(key string, member interface{}) (int64, error)
	SAdd(key string, member interface{}) (int64, error)
}

// Client for accessing to Redis
type Client struct {
	cli *redis.Client
}

func (r *Client) GetCli() *redis.Client {
	return r.cli
}

func (r *Client) Incr(key string) error {
	return r.cli.Incr(key).Err()
}

func (r *Client) Exists(key string) (bool, error) {
	res, err := r.cli.Exists(key).Result()
	if err != nil {
		return false, err
	}

	if res == 1 {
		return true, nil
	}

	return false, nil
}

func (r *Client) GetString(key string) string {
	return r.cli.Get(key).String()
}

// TTL - returns TTL by the KEY
func (r *Client) TTL(key string) (time.Duration, error) {
	return r.cli.TTL(key).Result()
}

// Get - implementation get value from redis by key
func (r *Client) Get(key string) ([]byte, error) {
	res, err := r.cli.Get(key).Bytes()

	if err == redis.Nil {
		// it means that key does not exist in redis
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return res, nil
}

// GetInt - implementation get int value from redis by key
func (r *Client) GetInt(key string) (int, error) {
	res, err := r.cli.Get(key).Int()

	if err == redis.Nil {
		// it means that key does not exist in redis
		return 0, nil
	}

	return res, err
}

// Set - implementation set int value to redis by key and set ttl
func (r *Client) SetInt(key string, value int, ttl int64) error {
	return r.cli.Set(key, value, time.Second*time.Duration(ttl)).Err()
}

// Set - implementation set value to redis by key and set ttl
func (r *Client) Set(key string, value []byte, ttl int64) error {
	return r.cli.Set(key, value, time.Second*time.Duration(ttl)).Err()
}

// Del - implementation delete keys from redis
func (r *Client) Del(keys ...string) error {
	return r.cli.Del(keys...).Err()
}

// HMSet - sets map(hash) to redis by key
func (r *Client) HMSet(key string, fields map[string]interface{}) error {
	return r.cli.HMSet(key, fields).Err()
}

// HMGet - returns `fields` from map(hash), found by `key`
func (r *Client) HMGet(key string, fields ...string) []interface{} {
	return r.cli.HMGet(key, fields...).Val()
}

// HGetAll - returns map(hash) by `key`
func (r *Client) HGetAll(key string) map[string]string {
	return r.cli.HGetAll(key).Val()
}

// HIncrBy - finds map(hash) by `key` and increments its `field`
func (r *Client) HIncrBy(key, field string, incr int64) error {
	return r.cli.HIncrBy(key, field, incr).Err()
}

// Expire - adds `exp` - expiration time for the `key`
func (r *Client) Expire(key string, exp time.Duration) error {
	return r.cli.Expire(key, exp).Err()
}

// Ping - pings the redis client connection
func (r *Client) Ping(ctx context.Context) error {
	return r.cli.WithContext(ctx).Ping().Err()
}

// Subscribe - subscribes to the pubSub topics
func (r *Client) Subscribe(channels ...string) *redis.PubSub {
	return r.cli.Subscribe(channels...)
}

func (r *Client) Publish(channel string, message interface{}) error {
	return r.cli.Publish(channel, message).Err()
}

func (r *Client) SIsMember(key string, member interface{}) (bool, error) {
	return r.cli.SIsMember(key, member).Result()
}
func (r *Client) SRem(key string, member interface{}) (int64, error) {
	return r.cli.SRem(key, member).Result()
}

func (r *Client) SAdd(key string, member interface{}) (int64, error) {
	return r.cli.SAdd(key, member).Result()
}

// TokenBlackListKey builds token black list key.
func TokenBlackListKey(tokenID fmt.Stringer) string {
	return fmt.Sprintf("%s:%s", blackListPrefix, tokenID.String())
}
