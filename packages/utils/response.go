package utils

import "github.com/gofiber/fiber/v2"

type APIResponse struct {
	Success bool        `json:"success"`
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Error   interface{} `json:"error"`
}

// Success Response
func Success(c *fiber.Ctx, status int, message string, data interface{}) error {
	return c.Status(status).JSON(APIResponse{
		Success: true,
		Status:  status,
		Message: message,
		Data:    data,
		Error:   nil,
	})
}

// Error Response
func Error(c *fiber.Ctx, status int, message string, err interface{}) error {
	return c.Status(status).JSON(APIResponse{
		Success: false,
		Status:  status,
		Message: message,
		Data:    nil,
		Error:   err,
	})
}