package rdb

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"time"
)

var ctx = context.Background()
var rdb *redis.Client

func RedisNewClient()  {
	rdb = redis.NewClient(&redis.Options{
		Network:            "",
		Addr:               "10.41.41.178:6379",
		Dialer:             nil,
		OnConnect:          nil,
		Username:           "",
		Password:           "",
		DB:                 0,
		MaxRetries:         0,
		MinRetryBackoff:    0,
		MaxRetryBackoff:    0,
		DialTimeout:        0,
		ReadTimeout:        0,
		WriteTimeout:       0,
		PoolSize:           0,
		MinIdleConns:       0,
		MaxConnAge:         0,
		PoolTimeout:        0,
		IdleTimeout:        0,
		IdleCheckFrequency: 0,
		TLSConfig:          nil,
		Limiter:            nil,
	})
	pong, err := rdb.Ping(ctx).Result()
	if err != nil {
		fmt.Println(pong, err)
	}
}

func Set(key, value string, expire time.Duration)  {
	err := rdb.Set(ctx, key, value, expire).Err()
	if err != nil {
		panic(err)
	}
}

func Get(key string) string {
	val, err := rdb.Get(ctx, key).Result()
	if err != nil {
		panic(err)
	}
	return val
}

func MSet(values ...interface{})  {
	err := rdb.MSet(ctx, values...).Err()
	if err != nil {
		panic(err)
	}
}

func SAdd(key string, members ...interface{}) (int64, error) {
	n, err := rdb.SAdd(ctx, key, members...).Result()
	return n, err
}

func Unlink(key string) (int64, error) {
	n, err := rdb.Unlink(ctx, key).Result()
	return n, err
}