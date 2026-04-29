package cryptowallet

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


///////////////////// user ///////////
func (h *Handler) CreateWallet(c *fiber.Ctx) error {
	var body CreateWalletRequest
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	userID := c.Locals("userID").(uint)

	if err := h.service.CreateWallet(c.UserContext(), userID, body.Symbol); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	return c.JSON("wallet created")
}

func (h *Handler) GetWallets(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	data, err := h.service.GetMyWallets(c.UserContext(), userID)
	if err != nil {
		return c.Status(500).JSON(err.Error())
	}

	return c.JSON(data)
}

func (h *Handler) GetWallet(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)
	symbol := c.Params("symbol")

	data, err := h.service.GetWallet(c.UserContext(), userID, symbol)
	if err != nil {
		return c.Status(400).JSON(err.Error())
	}

	return c.JSON(data)
}

func (h *Handler) GetSummary(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	total, err := h.service.GetSummary(c.UserContext(), userID)
	if err != nil {
		return c.Status(500).JSON(err.Error())
	}

	return c.JSON(fiber.Map{"total": total})
}

func (h *Handler) GetTransactions(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	symbol := c.Query("symbol", "")
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	page, _ := strconv.Atoi(c.Query("page", "1"))

	offset := (page - 1) * limit

	data, err := h.service.GetTransactions(c.UserContext(), userID, symbol, limit, offset)
	if err != nil {
		return c.Status(500).JSON(err.Error())
	}

	return c.JSON(data)
}

func (h *Handler) GetLocks(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	data, err := h.service.GetLocks(c.UserContext(), userID)
	if err != nil {
		return c.Status(500).JSON(err.Error())
	}

	return c.JSON(data)
}

///////////////////// admin wallet //////////

func (h *Handler) GetAllWalletsAdmin(c *fiber.Ctx) error {
	data, err := h.service.GetAllWalletsAdmin(c.UserContext())
	if err != nil {
		return c.Status(500).JSON(err.Error())
	}
	return c.JSON(data)
}

func (h *Handler) GetUserWalletsAdmin(c *fiber.Ctx) error {
	userID, _ := strconv.Atoi(c.Params("userId"))

	data, err := h.service.GetUserWalletsAdmin(c.UserContext(), uint(userID))
	if err != nil {
		return c.Status(500).JSON(err.Error())
	}

	return c.JSON(data)
}

func (h *Handler) GetUserWalletBySymbolAdmin(c *fiber.Ctx) error {
	userID, _ := strconv.Atoi(c.Params("userId"))
	symbol := c.Params("symbol")

	data, err := h.service.GetUserWalletBySymbolAdmin(c.UserContext(), uint(userID), symbol)
	if err != nil {
		return c.Status(400).JSON(err.Error())
	}

	return c.JSON(data)
}

func (h *Handler) AdminCredit(c *fiber.Ctx) error {
	userID, _ := strconv.Atoi(c.Params("userId"))

	var body AdminAmountRequest
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	err := h.service.AdminCredit(c.UserContext(), uint(userID), body.Symbol, body.Amount)
	if err != nil {
		return c.Status(400).JSON(err.Error())
	}

	return c.JSON("credited")
}

func (h *Handler) AdminDebit(c *fiber.Ctx) error {
	userID, _ := strconv.Atoi(c.Params("userId"))

	var body AdminAmountRequest
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	err := h.service.AdminDebit(c.UserContext(), uint(userID), body.Symbol, body.Amount)
	if err != nil {
		return c.Status(400).JSON(err.Error())
	}

	return c.JSON("debited")
}

func (h *Handler) FreezeWallet(c *fiber.Ctx) error {
	userID, _ := strconv.Atoi(c.Params("userId"))

	if err := h.service.FreezeWallet(c.UserContext(), uint(userID)); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	return c.JSON("wallet frozen")
}

func (h *Handler) UnfreezeWallet(c *fiber.Ctx) error {
	userID, _ := strconv.Atoi(c.Params("userId"))

	if err := h.service.UnfreezeWallet(c.UserContext(), uint(userID)); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	return c.JSON("wallet active")
}

func (h *Handler) GetAllTransactionsAdmin(c *fiber.Ctx) error {
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	page, _ := strconv.Atoi(c.Query("page", "1"))

	offset := (page - 1) * limit

	data, err := h.service.GetAllTransactionsAdmin(c.UserContext(), limit, offset)
	if err != nil {
		return c.Status(500).JSON(err.Error())
	}

	return c.JSON(data)
}

func (h *Handler) GetUserTransactionsAdmin(c *fiber.Ctx) error {
	userID, _ := strconv.Atoi(c.Params("userId"))

	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	page, _ := strconv.Atoi(c.Query("page", "1"))

	offset := (page - 1) * limit

	data, err := h.service.GetUserTransactionsAdmin(c.UserContext(), uint(userID), limit, offset)
	if err != nil {
		return c.Status(500).JSON(err.Error())
	}

	return c.JSON(data)
}

////////////////////asset admin/////////////

func (h *Handler) CreateAsset(c *fiber.Ctx) error {
	var body CreateAssetRequest

	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	err := h.service.CreateAsset(c.UserContext(), body.Symbol, body.Name, body.Precision)
	if err != nil {
		return c.Status(400).JSON(err.Error())
	}

	return c.JSON("asset created")
}

func (h *Handler) GetAssets(c *fiber.Ctx) error {
	data, err := h.service.GetAssets(c.UserContext())
	if err != nil {
		return c.Status(500).JSON(err.Error())
	}

	return c.JSON(data)
}

func (h *Handler) UpdateAsset(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))

	var body CreateAssetRequest
	c.BodyParser(&body)

	err := h.service.UpdateAsset(c.UserContext(), uint(id), body.Name, body.Precision)
	if err != nil {
		return c.Status(400).JSON(err.Error())
	}

	return c.JSON("updated")
}

func (h *Handler) UpdateAssetStatus(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))

	var body struct {
		Status string `json:"status"`
	}

	c.BodyParser(&body)

	err := h.service.UpdateAssetStatus(c.UserContext(), uint(id), body.Status)
	if err != nil {
		return c.Status(400).JSON(err.Error())
	}

	return c.JSON("status updated")
}