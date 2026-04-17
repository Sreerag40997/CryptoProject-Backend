package ecard

import (
	"cryptox/packages/utils"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	service Service
}

func NewHandler(s Service) *Handler {
	return &Handler{service: s}
}

func (h *Handler) GetMyCard(c *fiber.Ctx) error {

	userID := c.Locals("userID").(uint)

	data, err := h.service.GetMyCard(c.UserContext(), userID)
	if err != nil {
		return utils.Error(c, 404, "Card not found", err.Error())
	}

	return utils.Success(c, 200, "Card fetched", data)
}

func (h *Handler) BlockCard(c *fiber.Ctx) error {

	userID := c.Locals("userID").(uint)

	err := h.service.BlockCard(c.UserContext(), userID)
	if err != nil {
		return utils.Error(c, 500, "Failed to block card", err.Error())
	}

	return utils.Success(c, 200, "Card blocked successfully", nil)
}

func (h *Handler) UnblockCard(c *fiber.Ctx) error {

	userID := c.Locals("userID").(uint)

	err := h.service.UnblockCard(c.UserContext(), userID)
	if err != nil {
		return utils.Error(c, 500, "Failed to unblock card", err.Error())
	}

	return utils.Success(c, 200, "Card unblocked successfully", nil)
}

func (h *Handler) AdminGetCard(c *fiber.Ctx) error {

	userIDParam := c.Params("userId")

	userID, err := strconv.Atoi(userIDParam)
	if err != nil {
		return utils.Error(c, 400, "Invalid user ID", err.Error())
	}

	card, err := h.service.AdminGetCard(c.UserContext(), uint(userID))
	if err != nil {
		return utils.Error(c, 404, "Card not found", err.Error())
	}

	return utils.Success(c, 200, "Card fetched successfully", card)
}

func (h *Handler) AdminBlockCard(c *fiber.Ctx) error {

	userIDParam := c.Params("userId")

	userID, err := strconv.Atoi(userIDParam)
	if err != nil {
		return utils.Error(c, 400, "Invalid user ID", err.Error())
	}

	err = h.service.AdminBlockCard(c.UserContext(), uint(userID))
	if err != nil {
		return utils.Error(c, 500, "Failed to block card", err.Error())
	}

	return utils.Success(c, 200, "Card blocked by admin", nil)
}

func (h *Handler) AdminUnblockCard(c *fiber.Ctx) error {

	userIDParam := c.Params("userId")

	userID, err := strconv.Atoi(userIDParam)
	if err != nil {
		return utils.Error(c, 400, "Invalid user ID", err.Error())
	}

	err = h.service.AdminUnblockCard(c.UserContext(), uint(userID))
	if err != nil {
		return utils.Error(c, 500, "Failed to unblock card", err.Error())
	}

	return utils.Success(c, 200, "Card unblocked by admin", nil)
}