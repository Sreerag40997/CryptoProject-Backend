package utils

import (
	"cryptox/packages/redis"
	"time"

	"github.com/redis/go-redis/v9"
)

// -> save the otp on redis
func SaveOTP(email, otp string) error {
	return redisClient.Redis.Set(Ctx, email, otp, 5*time.Minute).Err()
}

// -> Get the otp on redis
func GetOTP(email string) (storedOTP string, err error) {

	storedOTP, err = redisClient.Redis.Get(Ctx, email).Result()
	if err != nil {
		if err == redis.Nil{
			return "", nil
		}
		return "", err
	}
	return storedOTP, err
}

// -> delete the otp from redis
func DeleteOTP(email string) error {
	return redisClient.Redis.Del(Ctx, email).Err()
}
