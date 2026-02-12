package handler

import (
	"tubexxi/video-api/internal/entity"
	"tubexxi/video-api/internal/middleware"
	"tubexxi/video-api/internal/service"
	"tubexxi/video-api/pkg/response"
	"tubexxi/video-api/pkg/utils"

	"github.com/gofiber/fiber/v2"
)

type ApplicationHandler struct {
	ctxinject *middleware.ContextMiddleware
	service   *service.ApplicationService
}

func NewApplicationHandler(ctxinject *middleware.ContextMiddleware, service *service.ApplicationService) *ApplicationHandler {
	return &ApplicationHandler{
		ctxinject: ctxinject,
		service:   service,
	}
}

func (h *ApplicationHandler) RegisterApplication(c *fiber.Ctx) error {
	ctx := h.ctxinject.HandlerContext(c)

	var req entity.RegisterNewApplicationRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, err.Error(), nil)
	}

	if errs := utils.ValidateStruct(req); errs != nil {
		return response.Error(c, fiber.StatusBadRequest, response.ValidationErrors{Errors: errs}.Error(), nil)
	}

	if err := h.service.RegisterApplication(ctx, &req); err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error(), nil)
	}
	return response.Success(c, "Application registered successful", nil)
}
func (h *ApplicationHandler) GetPublicAppConfig(c *fiber.Ctx) error {
	ctx := h.ctxinject.HandlerContext(c)

	packageName := c.Params("package_name")
	if packageName == "" {
		return response.Error(c, fiber.StatusBadRequest, "Package name is required", nil)
	}

	appConfig, err := h.service.GetPublicAppConfig(ctx, packageName)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error(), nil)
	}
	return response.Success(c, "Application config fetched successful", appConfig)
}
func (h *ApplicationHandler) UpdateAppConfigBulk(c *fiber.Ctx) error {
	ctx := h.ctxinject.HandlerContext(c)

	packageName := c.Params("package_name")
	if packageName == "" {
		return response.Error(c, fiber.StatusBadRequest, "Package name is required", nil)
	}

	var req []entity.UpdateApplicationBulkRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, err.Error(), nil)
	}

	if errs := utils.ValidateStruct(req); errs != nil {
		return response.Error(c, fiber.StatusBadRequest, response.ValidationErrors{Errors: errs}.Error(), nil)
	}

	if err := h.service.UpdateAppConfigBulk(ctx, packageName, req); err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error(), nil)
	}
	return response.Success(c, "Application config updated successful", nil)
}
