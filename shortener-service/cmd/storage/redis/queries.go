package redis

import (
	"context"

	redis "github.com/redis/go-redis/v9"
)

var redisClient *redis.Client

func ConnectionRedis(conn *redis.Client) {
	redisClient = conn
}

func SetData(key, value string) error {
	err := redisClient.Set(context.Background(), key, value, 0).Err()
	if err != nil {
		return err
	}
	return nil
}

func GetData(key string) (string, error) {
	val, err := redisClient.Get(context.Background(), key).Result()
	if err != nil {
		return "", err
	}
	return val, nil
}
