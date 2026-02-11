package handler

import (
	"tubexxi/video-api/internal/entity"
	"tubexxi/video-api/internal/middleware"
	"tubexxi/video-api/internal/service"
	"tubexxi/video-api/pkg/response"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type SettingHandler struct {
	ctxinject       *middleware.ContextMiddleware
	scopeMiddleware *middleware.ScopeMiddleware
	service         *service.SettingService
	logger          *zap.Logger
}

func NewSettingHandler(
	ctxinject *middleware.ContextMiddleware,
	scopeMiddleware *middleware.ScopeMiddleware,
	service *service.SettingService,
	logger *zap.Logger,
) *SettingHandler {
	return &SettingHandler{
		ctxinject:       ctxinject,
		scopeMiddleware: scopeMiddleware,
		service:         service,
		logger:          logger,
	}
}

func (h *SettingHandler) UpdateSettingsBulk(c *fiber.Ctx) error {
	ctx := h.ctxinject.HandlerContext(c)
	scope := c.Query("scope")
	if scope == "" {
		scope = h.scopeMiddleware.GetSettingsScope(c)
	}

	var settings []entity.UpdateSettingsBulkRequest
	if err := c.BodyParser(&settings); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	if err := h.service.UpdateSettingsBulk(ctx, scope, settings); err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to update settings", err.Error())
	}

	return response.Success(c, "Settings updated successfully", nil)
}
func (h *SettingHandler) UploadFile(c *fiber.Ctx) error {
	ctx := h.ctxinject.HandlerContext(c)
	scope := c.Query("scope")
	if scope == "" {
		scope = h.scopeMiddleware.GetSettingsScope(c)
	}

	req := new(entity.UploadFileRequest)
	if err := c.BodyParser(req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	if req.Key == "" {
		req.Key = c.FormValue("key")
	}

	if req.Key != "site_logo" && req.Key != "site_favicon" {
		return response.Error(c, fiber.StatusBadRequest, "Invalid key. Must be 'site_logo' or 'site_favicon'", nil)
	}

	file, err := c.FormFile("file")
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "File is required", err.Error())
	}

	url, err := h.service.UploadFile(ctx, scope, file, req.Key)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to upload file", err.Error())
	}

	return response.Success(c, "File uploaded successfully", fiber.Map{"url": url})
}
func (h *SettingHandler) GetPublicSettings(c *fiber.Ctx) error {
	ctx := h.ctxinject.HandlerContext(c)
	scope := c.Query("scope")
	if scope == "" {
		scope = h.scopeMiddleware.GetSettingsScope(c)
	}

	settings, err := h.service.GetPublicSettings(ctx, scope)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to fetch settings", err.Error())
	}
	return response.Success(c, "Settings fetched successfully", settings)
}
func (h *SettingHandler) GetAllSettings(c *fiber.Ctx) error {
	ctx := h.ctxinject.HandlerContext(c)
	scope := c.Query("scope")
	if scope == "" {
		scope = h.scopeMiddleware.GetSettingsScope(c)
	}

	settings, err := h.service.GetAllSettings(ctx, scope)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to fetch settings", err.Error())
	}
	return response.Success(c, "All settings fetched", settings)
}
