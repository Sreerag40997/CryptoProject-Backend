package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort        string
	DBUrl          string
	RedisAddr      string
	JWTSecret      string
	RazorpayKey    string
	RazorpaySecret string
}

func LeadConfig() *Config {

	err := godotenv.Load(".env")
	if err != nil {
		log.Println("No env file found")
		log.Println("Error loading .env:", err)
	}

	dbURL := os.Getenv("DATABASE_URL")

	//production case (neon)
	if dbURL ==""{

	//loads .env and validate
	required := []string{
		"APP_PORT",
		"DB_HOST",
		"DB_USER",
		"DB_PASSWORD",
		"DB_NAME",
		"DB_PORT",
		"REDIS_ADDR",
		"JWT_SECRET",
		"BREVO_FROM_EMAIL",
		"BREVO_FROM_NAME",
		"BREVO_API_KEY",
	}
	for _, key := range required {
		if os.Getenv(key) == "" {
			log.Fatalf("Missing required env variable: %s", key)
		}
	}
	//creating dburl
	dbURL = fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)
	}
	//return the hole struct
	return &Config{
		AppPort:   os.Getenv("APP_PORT"),
		DBUrl:     dbURL,
		RedisAddr: os.Getenv("REDIS_ADDR"),
		JWTSecret: os.Getenv("JWT_SECRET"),
		RazorpayKey: os.Getenv("RAZORPAY_KEY_ID"),
		RazorpaySecret: os.Getenv("RAZORPAY_SECRET"),
	}
}
