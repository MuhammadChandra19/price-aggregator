package config

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
)

func RedisClient() *redis.Client {
	// Initialize the global Redis client
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // Redis server address
		Password: "mypassword",     // Redis server password
		DB:       0,                // Redis database number
	})

	// Test the connection
	pong, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(pong)
	}

	return rdb
}
