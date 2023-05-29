package data

import (
	"context"
	"fmt"
	"log"
	"os"

	redis "github.com/redis/go-redis/v9"
)

var ctx = context.Background()

func ConnectToRedis() *redis.Client {
	fmt.Println("Go Redis Client")
	redisEndpoint := os.Getenv("REDIS_ENDPOINT")

	redisClient := redis.NewClient(&redis.Options{
		Addr:     redisEndpoint + ":6379",
		Password: "",
		DB:       0,
	})

	err := redisClient.Ping(ctx).Err()
	if err != nil {
		log.Println("Failed to ping Redis:", err)
	}
	log.Println("Connected to Redis")

	return redisClient
}
