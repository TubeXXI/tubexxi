package handler

import (
	"tubexxi/video-api/internal/dto"
	"tubexxi/video-api/internal/middleware"
	"tubexxi/video-api/internal/service"
	"tubexxi/video-api/pkg/response"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type ClientHandler struct {
	ctxinject     *middleware.ContextMiddleware
	clientService *service.ClientService
	logger        *zap.Logger
}

func NewClientHandler(
	ctxinject *middleware.ContextMiddleware,
	clientService *service.ClientService,
	logger *zap.Logger,
) *ClientHandler {
	return &ClientHandler{
		ctxinject:     ctxinject,
		clientService: clientService,
		logger:        logger,
	}
}

func (h *ClientHandler) SendContactEmail(c *fiber.Ctx) error {
	ctx := h.ctxinject.HandlerContext(c)

	var req dto.ContactRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, err.Error(), nil)
	}
	if err := h.clientService.SendContactEmail(ctx, &req, c.Get("Origin")); err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error(), nil)
	}
	return response.Success(c, "Contact email sent successfully", nil)
}
