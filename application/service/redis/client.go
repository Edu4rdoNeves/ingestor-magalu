package redis

import (
	"context"
	"time"

	"github.com/go-redis/redis"
)

type IRedisClient interface {
	IncrementCounter(key string, value float64) error
	GetValue(key string) (string, error)
	DeleteKey(key string) error
	GetKeysByPattern(pattern string) ([]string, error)
}

var ctx = context.Background()

type RedisClient struct {
	Client *redis.Client
}

type ClientConfig struct {
	Addr         string
	Password     string
	Db           int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

func NewRedisClient(clientConfig ClientConfig) *RedisClient {
	rdb := redis.NewClient(&redis.Options{
		Addr:         clientConfig.Addr,
		Password:     clientConfig.Password,
		DB:           clientConfig.Db,
		ReadTimeout:  clientConfig.ReadTimeout,
		WriteTimeout: clientConfig.WriteTimeout,
	})

	return &RedisClient{Client: rdb}
}
