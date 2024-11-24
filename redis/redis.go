package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"
)

type RedisClient interface {
	GET(ctx context.Context, key string) (string, error)
	SET(ctx context.Context, key string, value interface{}, expiry time.Duration) error
}

type Client struct {
	redisClient *redis.Client
}

func NewRedisClient() (Client, error) {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // redis server addr
		Password: "",               // empty password
		DB:       0,                // default db
	})

	pong, err := redisClient.Ping(context.Background()).Result()
	if err != nil {
		return Client{redisClient: redisClient}, err
	}
	fmt.Println("connected to redis: ", pong)
	return Client{redisClient: redisClient}, nil
}

func (r *Client) GET(ctx context.Context, key string) (string, error) {
	val, err := r.redisClient.Get(ctx, key).Result()
	if err != nil {
		return "", err
	}
	return val, nil
}

func (r *Client) SET(ctx context.Context, key string, value interface{}, expiry time.Duration) error {
	valueJSON, err := json.Marshal(value)
	if err != nil {
		return err
	}
	err = r.redisClient.Set(context.Background(), key, valueJSON, expiry).Err()
	if err != nil {
		return err
	}
	return nil
}
