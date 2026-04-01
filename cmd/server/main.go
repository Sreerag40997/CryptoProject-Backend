package main

import (
	centralroutes "cryptox/packages/central_routes"
	"cryptox/packages/config"
	"cryptox/packages/database"
	"cryptox/packages/redis"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
)

func main() {

	//Load config
	cfg := config.LeadConfig()

	//Connect postgres
	db, err := database.NewPostgresConnection(cfg)
	if err != nil {
		log.Fatal("Postgres connection failed:", err)
	}

	// Connect Redis
	rdb, err := redis.NewRedisClient(cfg)
	if err != nil {
		log.Fatal("Redis connection failed:", err)
	}

	//fiber engine
	app := fiber.New()

	//setup routes
	centralroutes.SetUp(app, db, rdb, cfg.JWTSecret)

	port := os.Getenv("APP_PORT")

	if port == "" {
		port = "8080"
	}

	log.Println("Server running on port", port)

	log.Fatal(app.Listen(":" + port) )
}