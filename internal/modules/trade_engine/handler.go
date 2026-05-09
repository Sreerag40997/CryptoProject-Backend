package tradeengine

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

/////////////////////////////////////////////////////////
// USER HANDLERS
/////////////////////////////////////////////////////////

// POST /trade/order
func (h *Handler) PlaceOrder(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	var body CreateOrderReq
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	err := h.service.PlaceOrder(c.UserContext(), userID, body)
	if err != nil {
		return c.Status(400).JSON(err.Error())
	}

	return c.JSON(fiber.Map{
		"message": "order placed",
	})
}

/////////////////////////////////////////////////////////

// GET /trade/orders
func (h *Handler) GetMyOrders(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	status := c.Query("status", "")
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	page, _ := strconv.Atoi(c.Query("page", "1"))

	offset := (page - 1) * limit

	data, err := h.service.GetMyOrders(c.UserContext(), userID, status, limit, offset)
	if err != nil {
		return c.Status(500).JSON(err.Error())
	}

	return c.JSON(data)
}

/////////////////////////////////////////////////////////

// GET /trade/order/:id
func (h *Handler) GetOrder(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	id, _ := strconv.Atoi(c.Params("id"))

	data, err := h.service.GetOrderByID(c.UserContext(), userID, uint(id))
	if err != nil {
		return c.Status(400).JSON(err.Error())
	}

	return c.JSON(data)
}

/////////////////////////////////////////////////////////

// DELETE /trade/order/:id
func (h *Handler) CancelOrder(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	id, _ := strconv.Atoi(c.Params("id"))

	err := h.service.CancelOrder(c.UserContext(), userID, uint(id))
	if err != nil {
		return c.Status(400).JSON(err.Error())
	}

	return c.JSON(fiber.Map{
		"message": "order cancelled",
	})
}

/////////////////////////////////////////////////////////

// GET /trade/trades
func (h *Handler) GetMyTrades(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	page, _ := strconv.Atoi(c.Query("page", "1"))

	offset := (page - 1) * limit

	data, err := h.service.GetMyTrades(c.UserContext(), userID, limit, offset)
	if err != nil {
		return c.Status(500).JSON(err.Error())
	}

	return c.JSON(data)
}

/////////////////////////////////////////////////////////

// GET /trade/order/:id/fills
func (h *Handler) GetOrderFills(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	id, _ := strconv.Atoi(c.Params("id"))

	data, err := h.service.GetOrderFills(c.UserContext(), userID, uint(id))
	if err != nil {
		return c.Status(400).JSON(err.Error())
	}

	return c.JSON(data)
}

/////////////////////////////////////////////////////////
// PUBLIC HANDLERS (NO AUTH REQUIRED)
/////////////////////////////////////////////////////////

// GET /trade/orderbook?symbol=BTC-INR
func (h *Handler) GetOrderBook(c *fiber.Ctx) error {
	symbol := c.Query("symbol")

	bids, asks, err := h.service.GetOrderBook(c.UserContext(), symbol)
	if err != nil {
		return c.Status(500).JSON(err.Error())
	}

	return c.JSON(fiber.Map{
		"bids": bids,
		"asks": asks,
	})
}

/////////////////////////////////////////////////////////

// GET /trade/history?symbol=BTC-INR
func (h *Handler) GetTradeHistory(c *fiber.Ctx) error {
	symbol := c.Query("symbol")

	limit, _ := strconv.Atoi(c.Query("limit", "20"))

	data, err := h.service.GetTradeHistory(c.UserContext(), symbol, limit)
	if err != nil {
		return c.Status(500).JSON(err.Error())
	}

	return c.JSON(data)
}

/////////////////////////////////////////////////////////
// ADMIN HANDLERS
/////////////////////////////////////////////////////////

// GET /admin/trade/orders
func (h *Handler) GetAllOrders(c *fiber.Ctx) error {
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	page, _ := strconv.Atoi(c.Query("page", "1"))

	offset := (page - 1) * limit

	data, err := h.service.GetAllOrders(c.UserContext(), limit, offset)
	if err != nil {
		return c.Status(500).JSON(err.Error())
	}

	return c.JSON(data)
}

/////////////////////////////////////////////////////////

// GET /admin/trade/trades
func (h *Handler) GetAllTrades(c *fiber.Ctx) error {
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	page, _ := strconv.Atoi(c.Query("page", "1"))

	offset := (page - 1) * limit

	data, err := h.service.GetAllTrades(c.UserContext(), limit, offset)
	if err != nil {
		return c.Status(500).JSON(err.Error())
	}

	return c.JSON(data)
}

/////////////////////////////////////////////////////////

// GET /admin/trade/user/:userId/orders
func (h *Handler) GetOrdersByUserAdmin(c *fiber.Ctx) error {
	userID, _ := strconv.Atoi(c.Params("userId"))

	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	page, _ := strconv.Atoi(c.Query("page", "1"))

	offset := (page - 1) * limit

	data, err := h.service.GetOrdersByUserAdmin(c.UserContext(), uint(userID), limit, offset)
	if err != nil {
		return c.Status(500).JSON(err.Error())
	}

	return c.JSON(data)
}


