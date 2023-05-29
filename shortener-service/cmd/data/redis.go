package data

import (
	"context"

	redis "github.com/redis/go-redis/v9"
)

var client *redis.Client

func ConnectionRedis(redisClient *redis.Client) {
	client = redisClient
}

func SetData(key, value string) error {
	err := client.Set(context.Background(), key, value, 0).Err()
	if err != nil {
		return err
	}
	return nil
}

func GetData(key string) (string, error) {
	val, err := client.Get(context.Background(), key).Result()
	if err != nil {
		return "", err
	}
	return val, nil
}
