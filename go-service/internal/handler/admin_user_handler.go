package handler

import (
	"tubexxi/video-api/internal/dto"
	"tubexxi/video-api/internal/middleware"
	"tubexxi/video-api/internal/service"
	"tubexxi/video-api/pkg/response"
	"tubexxi/video-api/pkg/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
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
func (h *AdminUserHandler) SearchUser(c *fiber.Ctx) error {
	ctx := h.ctxinject.HandlerContext(c)

	var req dto.QueryParamsRequest
	if err := c.QueryParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, err.Error(), nil)
	}
	if errs := utils.ValidateStruct(req); len(errs) > 0 {
		return response.Error(c, fiber.StatusUnprocessableEntity, response.ValidationErrors{Errors: errs}.Error(), nil)
	}

	users, pagination, err := h.svc.SearchUser(ctx, req)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, err.Error(), nil)
	}
	return response.SuccessWithMeta(c, "Users retrieved",
		users,
		pagination,
	)
}
func (h *AdminUserHandler) HardDeleteUser(c *fiber.Ctx) error {
	ctx := h.ctxinject.HandlerContext(c)

	idStr := c.Params("id")
	if idStr == "" {
		return response.Error(c, fiber.StatusBadRequest, "User ID is required", nil)
	}

	userID, err := uuid.Parse(idStr)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid user ID", err.Error())
	}

	if err := h.svc.HardDeleteUser(ctx, userID); err != nil {
		return response.Error(c, fiber.StatusBadRequest, err.Error(), nil)
	}
	return response.Success(c, "User deleted", nil)
}
func (h *AdminUserHandler) BulkDeleteUser(c *fiber.Ctx) error {
	ctx := h.ctxinject.HandlerContext(c)

	var req struct {
		IDs []string `json:"ids" validate:"required,dive,uuid4"`
	}
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, err.Error(), nil)
	}
	if errs := utils.ValidateStruct(req); len(errs) > 0 {
		return response.Error(c, fiber.StatusUnprocessableEntity, response.ValidationErrors{Errors: errs}.Error(), nil)
	}

	userIDs := make([]uuid.UUID, 0, len(req.IDs))
	for _, idStr := range req.IDs {
		userID, err := uuid.Parse(idStr)
		if err != nil {
			return response.Error(c, fiber.StatusBadRequest, "Invalid user ID", err.Error())
		}
		userIDs = append(userIDs, userID)
	}

	if err := h.svc.BulkDeleteUser(ctx, userIDs); err != nil {
		return response.Error(c, fiber.StatusBadRequest, err.Error(), nil)
	}
	return response.Success(c, "Users deleted", nil)
}
