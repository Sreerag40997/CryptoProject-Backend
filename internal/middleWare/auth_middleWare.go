package middleware

import (
	"context"
	redisClient "cryptox/packages/redis"
	"cryptox/packages/utils"

	"github.com/gofiber/fiber/v2"
)

var (
	ctx = context.Background()
)

func AuthMiddleWare(jwtSecret string) fiber.Handler {
	return func(c *fiber.Ctx) error {

		tokenStr := c.Cookies("access")
		if tokenStr == "" {
			return utils.Error(c, 401, "Missing Token", nil)
		}
		
		exists, err := redisClient.Redis.Exists(ctx, "blacklist:"+tokenStr).Result()
		if err != nil || exists == 1 {
			return utils.Error(c, 401, "Blacklist Token", nil)
		}

		claims, err := utils.VerifyToken(tokenStr, jwtSecret)
		if err != nil {
			return utils.Error(c, 401, "Authentication Failed", err)
		}

		c.Locals("userID", claims.UserID)
		c.Locals("email", claims.Email)
		c.Locals("role", claims.Role)

		return c.Next()
	}
}