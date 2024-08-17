package configs

import (
	"context"
	"log"
	"os"
	"strconv"

	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
)

var RedisClient *redis.Client
var ctx = context.Background()
func Initialize() {
    // Load environment variables from .env file
    if err := godotenv.Load(); err != nil {
        log.Fatalf("Error loading .env file: %v", err)
    }

    // Initialize Redis client
	db, err := strconv.Atoi(os.Getenv("REDIS_DB"))
	if err != nil {
		log.Fatalf("Error parsing REDIS_DB: %v", err)
	}
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       db,
	})

    // Test Redis connection
    if _, err := RedisClient.Ping(context.Background()).Result(); err != nil {
        log.Fatalf("Error connecting to Redis: %v", err)
    }
}

func AddToQueue(queueName string, value interface{}) error {
	_, err := RedisClient.LPush(ctx, queueName, value).Result()
	if err != nil {
		return err
	}
	return nil
}
