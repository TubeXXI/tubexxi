package handler

import (
	"errors"
	"strings"
	"tubexxi/video-api/internal/dto"
	"tubexxi/video-api/internal/middleware"
	"tubexxi/video-api/internal/service"
	"tubexxi/video-api/pkg/response"
	"tubexxi/video-api/pkg/utils"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type UserHandler struct {
	ctxinject   *middleware.ContextMiddleware
	userService *service.UserService
	rateLimiter *middleware.RateLimiterMiddleware
	logger      *zap.Logger
}

func NewUserHandler(
	ctxinject *middleware.ContextMiddleware,
	userService *service.UserService,
	rateLimiter *middleware.RateLimiterMiddleware,
	logger *zap.Logger,
) *UserHandler {
	return &UserHandler{
		ctxinject:   ctxinject,
		userService: userService,
		rateLimiter: rateLimiter,
		logger:      logger,
	}
}
func (h *UserHandler) GetCurrentUser(c *fiber.Ctx) error {
	ctx := h.ctxinject.HandlerContext(c)

	userID, ok := c.Locals("user_id").(string)
	if !ok || userID == "" {
		return response.Error(c, fiber.StatusInternalServerError, "User ID not found", nil)
	}
	user, err := h.userService.GetUserByID(ctx, userID)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error(), nil)
	}
	return response.Success(c, "User found", user)
}
func (h *UserHandler) Logout(c *fiber.Ctx) error {
	ctx := h.ctxinject.HandlerContext(c)

	userID, ok := c.Locals("user_id").(string)
	if !ok || userID == "" {
		return response.Error(c, fiber.StatusInternalServerError, "User ID not found", nil)
	}
	if err := h.userService.Logout(ctx, userID); err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error(), nil)
	}
	h.rateLimiter.ResetLimitCounters(c)
	return response.Success(c, "Logout success", nil)
}
func (h *UserHandler) ChangeAvatar(c *fiber.Ctx) error {
	ctx := h.ctxinject.HandlerContext(c)

	userID, ok := c.Locals("user_id").(string)
	if !ok || userID == "" {
		return response.Error(c, fiber.StatusUnauthorized, "User authentication required", nil)
	}

	c.Request().Header.Set("Content-Type", "multipart/form-data")

	form, err := c.MultipartForm()
	if err != nil {
		h.logger.Error("Failed to parse multipart form",
			zap.String("user_id", userID),
			zap.Error(err))
		return response.Error(c, fiber.StatusBadRequest, "Invalid form data", nil)
	}
	defer form.RemoveAll()

	files, exists := form.File["avatar"]
	if !exists || len(files) == 0 {
		return response.Error(c, fiber.StatusBadRequest, "Avatar file is required", nil)
	}

	avatarFile := files[0]

	const maxFileSize = 5 << 20 // 5MB
	if avatarFile.Size > maxFileSize {
		return response.Error(c, fiber.StatusBadRequest,
			"File too large. Maximum size is 5MB", nil)
	}

	fileHeader := avatarFile.Header.Get("Content-Type")
	if !h.isValidImageType(fileHeader) {
		return response.Error(c, fiber.StatusBadRequest,
			"Invalid image type. Only JPEG, JPG, PNG, GIF, and WebP are allowed", nil)
	}

	avatarURL, err := h.userService.ChangeAvatar(ctx, userID, avatarFile)
	if err != nil {
		h.logger.Error("Failed to change avatar",
			zap.String("user_id", userID),
			zap.Error(err))

		switch {
		case errors.Is(err, dto.ErrUserNotFound):
			return response.Error(c, fiber.StatusNotFound, "User not found", nil)
		case strings.Contains(err.Error(), "upload"):
			return response.Error(c, fiber.StatusInternalServerError, "Failed to upload avatar", nil)
		default:
			return response.Error(c, fiber.StatusInternalServerError, err.Error(), nil)
		}
	}
	return response.Success(c, "Avatar changed successfully", fiber.Map{
		"avatar_url": avatarURL,
	})
}
func (h *UserHandler) UpdateProfile(c *fiber.Ctx) error {
	ctx := h.ctxinject.HandlerContext(c)

	userID, ok := c.Locals("user_id").(string)
	if !ok || userID == "" {
		return response.Error(c, fiber.StatusUnauthorized, "User authentication required", nil)
	}

	var req dto.UpdateProfileRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request payload", nil)
	}

	if errs := utils.ValidateStruct(req); len(errs) > 0 {
		return response.Error(c, fiber.StatusUnprocessableEntity, response.ValidationErrors{Errors: errs}.Error(), nil)
	}

	userUpdate, err := h.userService.UpdateProfile(ctx, userID, &req)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error(), nil)
	}
	return response.Success(c, "Profile updated successfully", userUpdate)
}
func (h *UserHandler) isValidImageType(contentType string) bool {
	allowedTypes := map[string]bool{
		"image/jpeg": true,
		"image/jpg":  true,
		"image/png":  true,
		"image/gif":  true,
		"image/webp": true,
	}
	return allowedTypes[contentType]
}
func (h *UserHandler) ChangePassword(c *fiber.Ctx) error {
	ctx := h.ctxinject.HandlerContext(c)

	userID, ok := c.Locals("user_id").(string)
	if !ok || userID == "" {
		return response.Error(c, fiber.StatusUnauthorized, "User authentication required", nil)
	}

	var req dto.ChangePasswordRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, err.Error(), nil)
	}

	if errs := utils.ValidateStruct(req); len(errs) > 0 {
		return response.Error(c, fiber.StatusUnprocessableEntity, response.ValidationErrors{Errors: errs}.Error(), nil)
	}

	if err := h.userService.UpdatePassword(ctx, userID, &req); err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error(), nil)
	}
	return response.Success(c, "Password changed successfully", nil)
}
func (h *UserHandler) EnableTwoFactor(c *fiber.Ctx) error {
	ctx := h.ctxinject.HandlerContext(c)

	userID, ok := c.Locals("user_id").(string)
	if !ok || userID == "" {
		return response.Error(c, fiber.StatusUnauthorized, "User authentication required", nil)
	}
	var req dto.EnableTwoFactorRequest

	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, err.Error(), nil)
	}

	if errs := utils.ValidateStruct(req); len(errs) > 0 {
		return response.Error(c, fiber.StatusUnprocessableEntity, response.ValidationErrors{Errors: errs}.Error(), nil)
	}

	qr, secret, err := h.userService.EnableTwoFactor(ctx, userID, &req)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error(), nil)
	}
	return response.Success(c, "Two-factor authentication enabled successfully", fiber.Map{
		"qr_code": qr,
		"secret":  secret,
	})
}
func (h *UserHandler) VerifyTwoFactor(c *fiber.Ctx) error {
	ctx := h.ctxinject.HandlerContext(c)

	userID, ok := c.Locals("user_id").(string)
	if !ok || userID == "" {
		return response.Error(c, fiber.StatusUnauthorized, "User authentication required", nil)
	}

	var req dto.ActivateTwoFactorRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, err.Error(), nil)
	}

	if errs := utils.ValidateStruct(req); len(errs) > 0 {
		return response.Error(c, fiber.StatusUnprocessableEntity, response.ValidationErrors{Errors: errs}.Error(), nil)
	}

	if err := h.userService.ActivateTwoFactor(ctx, userID, &req); err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error(), nil)
	}
	return response.Success(c, "Two-factor authentication verified successfully", nil)
}
