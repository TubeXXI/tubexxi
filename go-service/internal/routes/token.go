package routes

import (
	"tubexxi/video-api/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

type TokenRoutes struct {
	path string
	csrf *middleware.CSRFMiddleware
}

func NewTokenRoutes(
	csrf *middleware.CSRFMiddleware,
) *TokenRoutes {
	return &TokenRoutes{
		path: "/token",
		csrf: csrf,
	}
}

func (r *TokenRoutes) RegisterRoutes(parent fiber.Router) {
	router := parent.Group(r.path)

	router.Get("/csrf", r.csrf.CSRFProtect(), func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"success": true,
			"data": fiber.Map{
				"csrf_token": c.Locals("csrf"),
			},
		})
	})
}
