package tradeengine

import (
	"cryptox/packages/utils"
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
		return utils.Error(c, 400, "invalid request", err.Error())
	}

	err := h.service.PlaceOrder(c.UserContext(), userID, body)
	if err != nil {
		return utils.Error(c, 400, "failed to place order", err.Error())
	}

	return utils.Success(c, 200, "order placed", nil)
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
		return utils.Error(c, 500, "failed to fetch orders", err.Error())
	}

	return utils.Success(c, 200, "orders fetched", data)
}

/////////////////////////////////////////////////////////

// GET /trade/order/:id
func (h *Handler) GetOrder(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	id, _ := strconv.Atoi(c.Params("id"))

	data, err := h.service.GetOrderByID(c.UserContext(), userID, uint(id))
	if err != nil {
		return utils.Error(c, 400, "failed to fetch order", err.Error())
	}

	return utils.Success(c, 200, "order fetched", data)
}

/////////////////////////////////////////////////////////

// DELETE /trade/order/:id
func (h *Handler) CancelOrder(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	id, _ := strconv.Atoi(c.Params("id"))

	err := h.service.CancelOrder(c.UserContext(), userID, uint(id))
	if err != nil {
		return utils.Error(c, 400, "failed to cancel order", err.Error())
	}

	return utils.Success(c, 200, "order cancelled", nil)
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
		return utils.Error(c, 500, "failed to fetch trades", err.Error())
	}

	return utils.Success(c, 200, "trades fetched", data)
}

/////////////////////////////////////////////////////////

// GET /trade/order/:id/fills
func (h *Handler) GetOrderFills(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	id, _ := strconv.Atoi(c.Params("id"))

	data, err := h.service.GetOrderFills(c.UserContext(), userID, uint(id))
	if err != nil {
		return utils.Error(c, 400, "failed to fetch fills", err.Error())
	}

	return utils.Success(c, 200, "fills fetched", data)
}

/////////////////////////////////////////////////////////
// PUBLIC HANDLERS (NO AUTH REQUIRED)
/////////////////////////////////////////////////////////

// GET /trade/orderbook?symbol=BTC-INR
func (h *Handler) GetOrderBook(c *fiber.Ctx) error {
	symbol := c.Query("symbol")

	bids, asks, err := h.service.GetOrderBook(c.UserContext(), symbol)
	if err != nil {
		return utils.Error(c, 500, "failed to get orderbook", err.Error())
	}

	return utils.Success(c, 200, "orderbook fetched", fiber.Map{
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
		return utils.Error(c, 500, "failed to fetch trade history", err.Error())
	}

	return utils.Success(c, 200, "trade history fetched", data)
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
		return utils.Error(c, 500, "failed to fetch all orders", err.Error())
	}

	return utils.Success(c, 200, "all orders fetched", data)
}

/////////////////////////////////////////////////////////

// GET /admin/trade/trades
func (h *Handler) GetAllTrades(c *fiber.Ctx) error {
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	page, _ := strconv.Atoi(c.Query("page", "1"))

	offset := (page - 1) * limit

	data, err := h.service.GetAllTrades(c.UserContext(), limit, offset)
	if err != nil {
		return utils.Error(c, 500, "failed to fetch all trades", err.Error())
	}

	return utils.Success(c, 200, "all trades fetched", data)
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
		return utils.Error(c, 500, "failed to fetch user orders", err.Error())
	}

	return utils.Success(c, 200, "user orders fetched", data)
}


