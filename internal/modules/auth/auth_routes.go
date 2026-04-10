package auth

import (
	middleware "cryptox/internal/middleWare"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func AuthRoutes(r fiber.Router, db *gorm.DB, redis *redis.Client, jwtSecret string) {

	repo := NewRepo(db)
	service := NewAuthService(repo, redis, jwtSecret)
	controller := NewAuthController(service)

	auth := r.Group("/auth")

	auth.Post("/register", controller.Register)
	auth.Post("/login", controller.Login)

	auth.Use(middleware.AuthMiddleWare(jwtSecret))
	auth.Post("/logout", controller.Logout)
	auth.Post("/refresh", controller.Refresh)
	auth.Post("/sendotp", controller.SendOTP)
	auth.Post("/verifyotp", controller.VerifyOTP)
}