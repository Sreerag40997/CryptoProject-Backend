package redis

import (
	"context"
	"cryptox/packages/config"
	"log"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

var Ctx = context.Background()

func NewRedisClient(cfg *config.Config) (*redis.Client, error) {

	var rdb *redis.Client

	//  Check if production Redis URL exists (Upstash)
	redisURL := os.Getenv("REDIS_URL")

	if redisURL != "" {

		//  Production Mode
		opt, err := redis.ParseURL(redisURL)
		if err != nil {
			return nil, err
		}

		rdb = redis.NewClient(opt)

		log.Println("Using production Redis (Upstash / Render)")

	} else {

		// + Local Development Mode
		rdb = redis.NewClient(&redis.Options{
			Addr:         cfg.RedisAddr,
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 5 * time.Second,
			PoolSize:     10,
		})

		log.Println("Using local Redis (Docker)")
	}

	// Test connection
	_, err := rdb.Ping(Ctx).Result()
	if err != nil {
		return nil, err
	}

	log.Println("Connected to Redis ")

	return rdb, nil
}