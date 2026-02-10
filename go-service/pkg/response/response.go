package response

import (
	"github.com/gofiber/fiber/v2"
)

type APIResponse struct {
	Status     int         `json:"status"`
	Success    bool        `json:"success"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data,omitempty"`
	Error      interface{} `json:"error,omitempty"`
	Pagination interface{} `json:"pagination,omitempty"`
}

func Success(c *fiber.Ctx, message string, data interface{}) error {
	return c.Status(fiber.StatusOK).JSON(APIResponse{
		Status:  fiber.StatusOK,
		Success: true,
		Message: message,
		Data:    data,
	})
}

func Created(c *fiber.Ctx, message string, data interface{}) error {
	return c.Status(fiber.StatusCreated).JSON(APIResponse{
		Status:  fiber.StatusCreated,
		Success: true,
		Message: message,
		Data:    data,
	})
}

func Error(c *fiber.Ctx, status int, message string, err interface{}) error {
	return c.Status(status).JSON(APIResponse{
		Status:  status,
		Success: false,
		Message: message,
		Error:   err,
	})
}

func SuccessWithMeta(c *fiber.Ctx, message string, data interface{}, pagination interface{}) error {
	return c.Status(fiber.StatusOK).JSON(APIResponse{
		Status:     fiber.StatusOK,
		Success:    true,
		Message:    message,
		Data:       data,
		Pagination: pagination,
	})
}
