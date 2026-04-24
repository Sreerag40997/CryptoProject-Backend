package cashwallet

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

func (h *Handler) SetPin(c *fiber.Ctx) error {

	var body struct {
		Pin string `json:"pin"`
	}

	if err := c.BodyParser(&body); err != nil {
		return utils.Error(c, 400, "invalid request", err.Error())
	}

	userID := c.Locals("userID").(uint)

	err := h.service.SetPin(c.UserContext(), userID, body.Pin)
	if err != nil {
		return utils.Error(c, 400, "failed to set pin", err.Error())
	}

	return utils.Success(c, 200, "pin set successfully", nil)
}

func (h *Handler) ChangePin(c *fiber.Ctx) error {

	var body struct {
		OldPin string `json:"old_pin"`
		NewPin string `json:"new_pin"`
	}

	if err := c.BodyParser(&body); err != nil {
		return utils.Error(c, 400, "invalid request", err.Error())
	}
	if err := utils.Validator.Struct(body); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "invalid input",
			"err":   err.Error(),
		})
	}

	userID := c.Locals("userID").(uint)

	err := h.service.ChangePin(c.UserContext(), userID, body.OldPin, body.NewPin)
	if err != nil {
		return utils.Error(c, 400, "failed to change pin", err.Error())
	}

	return utils.Success(c, 200, "pin changed", nil)
}

func (h *Handler) GetMyWallet(c *fiber.Ctx) error {

	userID := c.Locals("userID").(uint)

	wallet, err := h.service.GetMyWallet(c.UserContext(), userID)
	if err != nil {
		return utils.Error(c, 500, "failed", err.Error())
	}

	return utils.Success(c, 200, "wallet fetched", wallet)
}

func (h *Handler) GetBalance(c *fiber.Ctx) error {

	userID := c.Locals("userID").(uint)

	balance, err := h.service.GetBalance(c.UserContext(), userID)
	if err != nil {
		return utils.Error(c, 500, "failed", err.Error())
	}

	return utils.Success(c, 200, "balance", balance)
}

func (h *Handler) GetTransactions(c *fiber.Ctx) error {

	userID := c.Locals("userID").(uint)

	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	page, _ := strconv.Atoi(c.Query("page", "1"))
	offset := (page - 1) * limit

	txns, err := h.service.GetTransactions(c.UserContext(), userID, limit, offset)
	if err != nil {
		return utils.Error(c, 500, "failed", err.Error())
	}

	return utils.Success(c, 200, "transactions", txns)
}

func (h *Handler) Deposit(c *fiber.Ctx) error {

	var body struct {
		Amount int64 `json:"amount"`
	}

	if err := c.BodyParser(&body); err != nil {
		return utils.Error(c, 400, "invalid request", err.Error())
	}

	userID := c.Locals("userID").(uint)

	orderID, err := h.service.CreateDepositOrder(c.UserContext(), userID, body.Amount)
	if err != nil {
		return utils.Error(c, 500, "failed to create order", err.Error())
	}

	return utils.Success(c, 200, "order created", map[string]string{
		"order_id": orderID,
	})
}

func (h *Handler) Withdraw(c *fiber.Ctx) error {

	var body struct {
		Amount int64  `json:"amount"`
		Pin    string `json:"pin"`
	}

	if err := c.BodyParser(&body); err != nil {
		return utils.Error(c, 400, "invalid request", err.Error())
	}

	userID := c.Locals("userID").(uint)

	err := h.service.Withdraw(c.UserContext(), userID, body.Amount, body.Pin)
	if err != nil {
		return utils.Error(c, 400, "withdraw failed", err.Error())
	}

	return utils.Success(c, 200, "withdraw success", nil)
}


func (h *Handler) BlockWallet(c *fiber.Ctx) error {
	userID, _ := strconv.Atoi(c.Params("userId"))
	err := h.service.AdminBlockWallet(c.UserContext(), uint(userID))
	if err != nil {
		return utils.Error(c, 400, "failed", err.Error())
	}
	return utils.Success(c, 200, "wallet blocked", nil)
}

func (h *Handler) FreezeWallet(c *fiber.Ctx) error {
	userID, _ := strconv.Atoi(c.Params("userId"))
	err := h.service.AdminFreezeWallet(c.UserContext(), uint(userID))
	if err != nil {
		return utils.Error(c, 400, "failed", err.Error())
	}
	return utils.Success(c, 200, "wallet frozen", nil)
}

func (h *Handler) UnblockWallet(c *fiber.Ctx) error {
	userID, _ := strconv.Atoi(c.Params("userId"))
	err := h.service.AdminUnblockWallet(c.UserContext(), uint(userID))
	if err != nil {
		return utils.Error(c, 400, "failed", err.Error())
	}
	return utils.Success(c, 200, "wallet active", nil)
}

func (h *Handler) AdminCredit(c *fiber.Ctx) error {

	userID, _ := strconv.Atoi(c.Params("userId"))

	var body struct {
		Amount int64 `json:"amount"`
	}

	c.BodyParser(&body)

	err := h.service.AdminCredit(c.UserContext(), uint(userID), body.Amount)
	if err != nil {
		return utils.Error(c, 400, "failed", err.Error())
	}

	return utils.Success(c, 200, "credited", nil)
}

func (h *Handler) AdminDebit(c *fiber.Ctx) error {

	userID, _ := strconv.Atoi(c.Params("userId"))

	var body struct {
		Amount int64 `json:"amount"`
	}

	if err := c.BodyParser(&body); err != nil {
		return utils.Error(c, 400, "invalid request", err.Error())
	}

	err := h.service.AdminDebit(c.UserContext(), uint(userID), body.Amount)
	if err != nil {
		return utils.Error(c, 400, "failed", err.Error())
	}

	return utils.Success(c, 200, "debited", nil)
}

func (h *Handler) RazorpayWebhook(c *fiber.Ctx) error {

	var payload struct {
		Event string `json:"event"`
		Payload struct {
			Payment struct {
				Entity struct {
					ID      string `json:"id"`
					Amount  int64  `json:"amount"`
					OrderID string `json:"order_id"`
					Notes struct {
						UserID uint `json:"user_id"`
					} `json:"notes"`
				} `json:"entity"`
			} `json:"payment"`
		} `json:"payload"`
	}

	if err := c.BodyParser(&payload); err != nil {
		return utils.Error(c, 400, "invalid payload", err.Error())
	}

	//  only process successful payments
	if payload.Event != "payment.captured" {
		return utils.Success(c, 200, "ignored", nil)
	}

	paymentID := payload.Payload.Payment.Entity.ID
	orderID := payload.Payload.Payment.Entity.OrderID
	userID := payload.Payload.Payment.Entity.Notes.UserID
	amount := payload.Payload.Payment.Entity.Amount

	// get signature
	signature := c.Get("X-Razorpay-Signature")

	// VERIFY SIGNATURE
	if !h.service.VerifyPayment(orderID, paymentID, signature) {
		return utils.Error(c, 400, "invalid signature", nil)
	}

	// CREDIT WALLET
	err := h.service.HandleDepositSuccess(c.UserContext(), userID, amount, paymentID)
	if err != nil {
		return utils.Error(c, 500, "deposit failed", err.Error())
	}

	return utils.Success(c, 200, "payment processed", nil)
}