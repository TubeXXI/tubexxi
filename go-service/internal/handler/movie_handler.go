package handler

import (
	"tubexxi/video-api/internal/middleware"
	"tubexxi/video-api/internal/service"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type MovieHandler struct {
	contextMiddleware *middleware.ContextMiddleware
	service           *service.MovieService
	rateLimiter       *middleware.RateLimiterMiddleware
	logger            *zap.Logger
}

func NewMovieHandler(
	contextMiddleware *middleware.ContextMiddleware,
	service *service.MovieService,
	rateLimiter *middleware.RateLimiterMiddleware,
	logger *zap.Logger,
) *MovieHandler {
	return &MovieHandler{
		contextMiddleware: contextMiddleware,
		service:           service,
		rateLimiter:       rateLimiter,
		logger:            logger,
	}
}

func (h *MovieHandler) GetHome(c *fiber.Ctx) error {
	ctx := c.UserContext()
	
	response, err := h.service.GetHome(ctx)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"data": response,
	})
}
