package middleware

import (
	"cryptox/packages/utils"

	"github.com/gofiber/fiber/v2"
)

func RequireRole(reqRole string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		role, ok := c.Locals("role").(string)
		if !ok {
			return utils.Error(c, 401, "unauthorized",nil)
		}
		if role != reqRole {
			return utils.Error(c, 401, "unauthorized",nil)
		}
		return c.Next()
	}
}