package tradeengine

import (
	middleware "cryptox/internal/middleWare"

	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(app fiber.Router, service Service, jwtSecret string) {

	h := NewHandler(service)

	/////////////////////////////////////////////////////////
	// USER ROUTES (AUTH REQUIRED)
	/////////////////////////////////////////////////////////

	trade := app.Group(
		"/trade",
		middleware.AuthMiddleWare(jwtSecret),
	)

	trade.Post("/order", h.PlaceOrder)
	trade.Get("/orders", h.GetMyOrders)
	trade.Get("/order/:id", h.GetOrder)
	trade.Delete("/order/:id", h.CancelOrder)

	trade.Get("/trades", h.GetMyTrades)
	trade.Get("/order/:id/fills", h.GetOrderFills)

	/////////////////////////////////////////////////////////
	// PUBLIC ROUTES (NO AUTH)
	/////////////////////////////////////////////////////////

	public := app.Group("/trade")

	public.Get("/orderbook", h.GetOrderBook)
	public.Get("/history", h.GetTradeHistory)

	/////////////////////////////////////////////////////////
	// ADMIN ROUTES (AUTH + ROLE)
	/////////////////////////////////////////////////////////

	admin := app.Group(
		"/admin/trade",
		middleware.AuthMiddleWare(jwtSecret),
		middleware.RequireRole("admin"),
	)

	admin.Get("/orders", h.GetAllOrders)
	admin.Get("/trades", h.GetAllTrades)
	admin.Get("/user/:userId/orders", h.GetOrdersByUserAdmin)
}


