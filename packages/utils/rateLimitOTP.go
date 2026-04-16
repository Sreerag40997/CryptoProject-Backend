package utils

import (
	"context"
	 "cryptox/packages/redis"
	"time"

)

var (
	Ctx = context.Background()
)

// -> Rate limiting the request it will avoid multiple requests
func RateLimitOTP(email string) (bool, error) {
	key := "OTP:RateLimit:" + email
	count, err := redisClient.Redis.Incr(Ctx, key).Result()
	if err != nil {
		return false, err
	}
	
	if count == 1 {
		redisClient.Redis.Expire(Ctx, key, 5*time.Minute)
	}

	if count > 3 {
		return false, nil
	}
	return true, nil
}