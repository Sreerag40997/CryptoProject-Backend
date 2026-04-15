package ecard

import (
	"cryptox/packages/utils"

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