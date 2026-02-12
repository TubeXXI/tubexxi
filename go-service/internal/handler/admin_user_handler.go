package handler

import (
	"tubexxi/video-api/internal/dto"
	"tubexxi/video-api/internal/middleware"
	"tubexxi/video-api/internal/service"
	"tubexxi/video-api/pkg/response"
	"tubexxi/video-api/pkg/utils"

	"github.com/gofiber/fiber/v2"
)

type AdminUserHandler struct {
	ctxinject *middleware.ContextMiddleware
	svc       *service.AdminUserService
}

func NewAdminUserHandler(ctxinject *middleware.ContextMiddleware, svc *service.AdminUserService) *AdminUserHandler {
	return &AdminUserHandler{ctxinject: ctxinject, svc: svc}
}

func (h *AdminUserHandler) SetRole(c *fiber.Ctx) error {
	ctx := h.ctxinject.HandlerContext(c)

	var req dto.SetUserRoleRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, err.Error(), nil)
	}
	if errs := utils.ValidateStruct(req); len(errs) > 0 {
		return response.Error(c, fiber.StatusUnprocessableEntity, response.ValidationErrors{Errors: errs}.Error(), nil)
	}

	if err := h.svc.SetUserRole(ctx, &req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, err.Error(), nil)
	}
	return response.Success(c, "Role updated", nil)
}
