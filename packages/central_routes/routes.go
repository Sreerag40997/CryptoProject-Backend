package centralroutes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func SetUp(
	app *fiber.App,
	db *gorm.DB,
	rdb *redis.Client,
	jwtSecret string,
) {

}
