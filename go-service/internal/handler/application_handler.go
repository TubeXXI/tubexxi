package handler

import (
	"tubexxi/video-api/internal/dto"
	"tubexxi/video-api/internal/entity"
	"tubexxi/video-api/internal/middleware"
	"tubexxi/video-api/internal/service"
	"tubexxi/video-api/pkg/logger"
	"tubexxi/video-api/pkg/response"
	"tubexxi/video-api/pkg/utils"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
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

	var req []entity.RegisterNewApplicationRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, err.Error(), nil)
	}

	if len(req) == 0 {
		return response.Error(c, fiber.StatusBadRequest, "Request body cannot be empty", nil)
	}

	if err := h.service.RegisterApplication(ctx, req); err != nil {
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

	if len(req) == 0 {
		return response.Error(c, fiber.StatusBadRequest, "Request body cannot be empty", nil)
	}

	if err := h.service.UpdateAppConfigBulk(ctx, packageName, req); err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error(), nil)
	}
	return response.Success(c, "Application config updated successful", nil)
}
func (h *ApplicationHandler) Search(c *fiber.Ctx) error {
	ctx := h.ctxinject.HandlerContext(c)

	var req dto.QueryParamsRequest
	if err := c.QueryParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, err.Error(), nil)
	}

	if errs := utils.ValidateStruct(req); errs != nil {
		logger.Logger.Error("Validation errors", zap.Error(response.ValidationErrors{Errors: errs}))
		return response.Error(c, fiber.StatusBadRequest, response.ValidationErrors{Errors: errs}.Error(), nil)
	}

	apps, pagination, err := h.service.Search(ctx, req)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error(), nil)
	}
	return response.SuccessWithMeta(c, "Applications fetched successful",
		apps,
		pagination,
	)
}
func (h *ApplicationHandler) GetByPackageName(c *fiber.Ctx) error {
	ctx := h.ctxinject.HandlerContext(c)

	packageName := c.Params("package_name")
	if packageName == "" {
		return response.Error(c, fiber.StatusBadRequest, "Package name is required", nil)
	}

	app, err := h.service.GetByPackageName(ctx, packageName)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error(), nil)
	}
	return response.Success(c, "Application fetched successful", app)
}
func (h *ApplicationHandler) DeleteApplication(c *fiber.Ctx) error {
	ctx := h.ctxinject.HandlerContext(c)

	packageName := c.Params("package_name")
	if packageName == "" {
		return response.Error(c, fiber.StatusBadRequest, "Package name is required", nil)
	}

	if err := h.service.DeleteApplication(ctx, packageName); err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error(), nil)
	}
	return response.Success(c, "Application deleted successful", nil)
}
func (h *ApplicationHandler) BulkDeleteApplication(c *fiber.Ctx) error {
	ctx := h.ctxinject.HandlerContext(c)

	var req struct {
		PackageNames []string `json:"package_names" validate:"required,dive"`
	}
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, err.Error(), nil)
	}

	if errs := utils.ValidateStruct(req); errs != nil {
		return response.Error(c, fiber.StatusBadRequest, response.ValidationErrors{Errors: errs}.Error(), nil)
	}

	if err := h.service.BulkDeleteApplication(ctx, req.PackageNames); err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error(), nil)
	}
	return response.Success(c, "Applications deleted successful", nil)
}
