package market

import (
	"cryptox/packages/utils"
	"encoding/json"
	"log"
	"strings"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

type MarketController struct {
	repo *MarketRepository
	hub *Hub
}

func NewMarketController(repo *MarketRepository, hub *Hub) *MarketController {
	return &MarketController{repo: repo, hub: hub}
}

func (c *MarketController) GetTicker(ctx *fiber.Ctx) error {

	symbol := ctx.Params("symbol")

	data, err := c.repo.GetTicker(ctx.Context(), symbol)
	if err != nil {
		return utils.Error(ctx, 404, "ticker not found", err)
	}

	return utils.Success(ctx, 200, "success", data)
}

func (c *MarketController) Socket(conn *websocket.Conn) {
	log.Println("client connected")

	var joined []string

	// read first message (symbols list)
	_, msg, err := conn.ReadMessage()
	if err != nil {
		log.Println("read error:", err)
		return
	}

	var req struct {
		Symbols []string `json:"symbols"`
	}

	if err := json.Unmarshal(msg, &req); err != nil {
		log.Println("invalid json:", err)
		return
	}

	// subscribe client to symbol rooms
	for _, s := range req.Symbols {
		s = strings.ToUpper(strings.TrimSpace(s))
		if s == "" {
			continue
		}

		c.hub.Add(s, conn)
		joined = append(joined, s)

		log.Println("subscribed:", s)
	}

	// cleanup on disconnect
	defer func() {
		for _, s := range joined {
			c.hub.Remove(s, conn)
		}
		log.Println("client disconnected")
	}()

	// keep connection alive
	for {
		if _, _, err := conn.ReadMessage(); err != nil {
			break
		}
	} 
}