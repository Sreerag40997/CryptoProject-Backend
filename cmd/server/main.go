package main

import (
	"cryptox/cmd/seed"
	centralroutes "cryptox/packages/central_routes"
	"cryptox/packages/cloudinary"
	"cryptox/packages/config"
	"cryptox/packages/database"
	redisClient "cryptox/packages/redis"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
)

func main() {

	//Load config
	cfg := config.LeadConfig()

	//cloudinary
	cloudinary.InitCloudinary()

	//Connect postgres
	db, err := database.NewPostgresConnection(cfg)
	if err != nil {
		log.Fatal("Postgres connection failed:", err)
	}

	//add admin
	seed.Addadmin(db)

	// Connect Redis
	rdb, err := redisClient.NewRedisClient(cfg)
	if err != nil {
		log.Fatal("Redis connection failed:", err)
	}

	//fiber engine
	app := fiber.New(fiber.Config{
		BodyLimit: 20 * 1024 *1024, //set the limit to 20mb for image uploading
	})

	//setup routes
	centralroutes.SetUp(app, db, rdb, cfg.JWTSecret)

	port := os.Getenv("APP_PORT")

	if port == "" {
		port = "8080"
	}

	log.Println("Server running on port", port)

	log.Fatal(app.Listen(":" + port))
}
