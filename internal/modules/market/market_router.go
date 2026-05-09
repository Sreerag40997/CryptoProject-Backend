package market

import (
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
)

func RegisterRoutes(r fiber.Router, redis *redis.Client) {

	hub := NewHub()
	repo := NewMarketRepo(redis)
	marketService := NewMarketService(repo, hub)
	marketController := NewMarketController(repo, hub)

	go StartBinanceStream(marketService)

	market := r.Group("/market")

	market.Get("/ticker/:symbol", marketController.GetTicker)

	market.Get("/ws/market", websocket.New(marketController.Socket))
}