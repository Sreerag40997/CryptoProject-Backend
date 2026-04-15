package kyc

import (
	"cryptox/packages/utils"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	service Service
}

func NewHandler(s Service) *Handler {
	return &Handler{service: s}
}

// submit kyc
func (h *Handler) SubmitKYC(c *fiber.Ctx) error {

	var req SubmitKYCRequest

	//  parse text fields
	if err := c.BodyParser(&req); err != nil {
		return utils.Error(c, 400, "Invalid input", err.Error())
	}

	// manually get files
	aadhaarFront, err := c.FormFile("aadhaar_front")
	if err != nil {
		return utils.Error(c, 400, "aadhaar_front is required", err.Error())
	}

	aadhaarBack, err := c.FormFile("aadhaar_back")
	if err != nil {
		return utils.Error(c, 400, "aadhaar_back is required", err.Error())
	}

	panFile, err := c.FormFile("pan_file")
	if err != nil {
		return utils.Error(c, 400, "pan_file is required", err.Error())
	}

	selfie, err := c.FormFile("selfie")
	if err != nil {
		return utils.Error(c, 400, "selfie is required", err.Error())
	}

	// assign to DTO
	req.AadhaarFront = aadhaarFront
	req.AadhaarBack = aadhaarBack
	req.PANFile = panFile
	req.Selfie = selfie

	userID := c.Locals("userID").(uint)

	err = h.service.SubmitKYC(c.UserContext(), userID, &req)
	if err != nil {
		return utils.Error(c, 500, "KYC submission failed", err.Error())
	}

	return utils.Success(c, 200, "KYC submitted successfully", nil)
}

func (h *Handler) GetKYCStatus(c *fiber.Ctx) error {

	userID := c.Locals("userID").(uint)

	result, err := h.service.GetKYCStatus(c.UserContext(), userID)
	if err != nil {
		return utils.Error(c, 500, "Failed to fetch KYC status", err.Error())
	}

	// Build proper message
	status := result["status"]

	var message string

	switch status {
	case "not_submitted":
		message = "KYC not submitted"
	case "pending":
		message = "KYC under review"
	case "approved":
		message = "KYC approved"
	case "rejected":
		message = "KYC rejected"
	}

	return utils.Success(c, 200, message, result)
}


func (h *Handler) GetMyKYC(c *fiber.Ctx) error {

	userID := c.Locals("userID").(uint)

	data, err := h.service.GetMyKYC(c.UserContext(), userID)
	if err != nil {
		return utils.Error(c, 404, "KYC not found", err.Error())
	}

	return utils.Success(c, 200, "KYC details fetched", data)
}

func (h *Handler) UpdateKYC(c *fiber.Ctx) error {

	var req UpdateKYCRequest

	if err := c.BodyParser(&req); err != nil {
		return utils.Error(c, 400, "Invalid input", err.Error())
	}

	userID := c.Locals("userID").(uint)

	err := h.service.UpdateKYC(c.UserContext(), userID, &req)
	if err != nil {
		return utils.Error(c, 400, err.Error(), nil)
	}

	return utils.Success(c, 200, "KYC updated successfully", nil)
}

func (h *Handler) GetKYCList(c *fiber.Ctx) error {

	status := c.Query("status")
	pageQuery := c.Query("page", "1")

	var page int
	fmt.Sscanf(pageQuery, "%d", &page)

	if page <= 0 {
		page = 1
	}

	limit := 10

	data, err := h.service.GetKYCList(c.UserContext(), status, page, limit)
	if err != nil {
		return utils.Error(c, 500, "Failed to fetch KYC list", err.Error())
	}

	return utils.Success(c, 200, "KYC list fetched", data)
}

func (h *Handler) GetKYCByID(c *fiber.Ctx) error {

	idParam := c.Params("id")

	var id uint
	fmt.Sscanf(idParam, "%d", &id)

	data, err := h.service.GetKYCByID(c.UserContext(), id)
	if err != nil {
		return utils.Error(c, 404, "KYC not found", err.Error())
	}

	return utils.Success(c, 200, "KYC fetched", data)
}

func (h *Handler) UpdateKYCStatus(c *fiber.Ctx) error {

	idParam := c.Params("id")

	var id uint
	fmt.Sscanf(idParam, "%d", &id)

	var req  ReviewKYCRequest

	if err := c.BodyParser(&req); err != nil {
		return utils.Error(c, 400, "Invalid input", err.Error())
	}

	err := h.service.UpdateKYCStatus(c.UserContext(), id, req.Status, req.Reason)
	if err != nil {
		return utils.Error(c, 500, "Failed to update KYC", err.Error())
	}

	return utils.Success(c, 200, "KYC updated successfully", nil)
}