package database

import (
	"context"
	"os"
	"time"

	"github.com/Nowap83/FrameRate/backend/internal/utils"
	"github.com/redis/go-redis/v9"
)

func ConnectRedis() (*redis.Client, error) {
	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		redisURL = "localhost:6379"
	}

	password := os.Getenv("REDIS_PASSWORD")

	rdb := redis.NewClient(&redis.Options{
		Addr:     redisURL,
		Password: password,
		DB:       0, // base par défaut
	})

	// test de connexion
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if _, err := rdb.Ping(ctx).Result(); err != nil {
		return nil, err
	}

	utils.Log.Info("Redis connected successfully")
	return rdb, nil
}
