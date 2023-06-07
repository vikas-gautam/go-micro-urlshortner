package redis

import (
	"context"
	"fmt"
	"os"

	redis "github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

var ctx = context.Background()

func ConnectToRedis() (*redis.Client, error) {
	fmt.Println("Go Redis Client")
	redisEndpoint := os.Getenv("REDIS_ENDPOINT")

	redisClient := redis.NewClient(&redis.Options{
		Addr:     redisEndpoint + ":6379",
		Password: "",
		DB:       0,
	})

	err := redisClient.Ping(ctx).Err()
	if err != nil {
		logrus.Errorf("Failed to connect to redis:", err)
		return redisClient, err
	}
	

	logrus.Info("Connected to Redis")
	return redisClient, nil
}
