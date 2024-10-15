package redis_client

import (
	"context"
	"strconv"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

func ConnectToRedis() (*redis.Client, error) {
	// Create a new Redis client
	add, pass, dbStr := viper.GetString("REDIS_URI"), viper.GetString("REDIS_PASSWORD"), viper.GetString("REDIS_DB")

	db, err := strconv.Atoi(dbStr)
	if err != nil {
		return nil, err
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr:     add,  // Update with your Redis address
		Password: pass, // Add password if Redis requires authentication
		DB:       db,   // Use default DB
	})

	// Test Redis connection
	_, err = redisClient.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}

	return redisClient, nil
}
