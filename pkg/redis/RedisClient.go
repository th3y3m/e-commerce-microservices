package redis_client

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

func ConnectToRedis() (*redis.Client, error) {
	// Get Redis connection parameters from environment variables
	add, pass, dbStr := viper.GetString("REDIS_URI"), viper.GetString("REDIS_PASSWORD"), viper.GetString("REDIS_DB")

	db, err := strconv.Atoi(dbStr)
	if err != nil {
		return nil, fmt.Errorf("invalid Redis DB number: %w", err)
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr:     add,  // Redis address
		Password: pass, // Redis password
		DB:       db,   // Redis DB number
	})

	// Retry logic for Redis connection
	for i := 0; i < 5; i++ { // Retry up to 5 times
		_, err = redisClient.Ping(context.Background()).Result()
		if err == nil {
			log.Println("Connected to Redis")
			return redisClient, nil
		}
		log.Printf("Failed to connect to Redis (attempt %d/5): %v", i+1, err)
		time.Sleep(1 * time.Second) // Wait before retrying
	}

	return nil, fmt.Errorf("failed to connect to Redis after multiple attempts: %w", err)
}
