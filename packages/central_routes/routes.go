package centralroutes

import (
	"cryptox/internal/modules/auth"
	"cryptox/internal/modules/profile"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func SetUp(app *fiber.App, db *gorm.DB, rdb *redis.Client, jwtSecret string) {

	api := app.Group("/api")

	auth.AuthRoutes(api, db, rdb, jwtSecret)
	profile.ProfileRoutes(api, db, jwtSecret)

}
