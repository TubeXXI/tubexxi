package handler

import (
	"tubexxi/video-api/internal/dto"
	"tubexxi/video-api/internal/middleware"
	"tubexxi/video-api/internal/service"
	"tubexxi/video-api/pkg/response"
	"tubexxi/video-api/pkg/utils"

	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	ctxinject   *middleware.ContextMiddleware
	authService *service.AuthService
}

func NewAuthHandler(ctxinject *middleware.ContextMiddleware, authService *service.AuthService) *AuthHandler {
	return &AuthHandler{ctxinject: ctxinject, authService: authService}
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	ctx := h.ctxinject.HandlerContext(c)

	var req dto.FirebaseLoginRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, err.Error(), nil)
	}
	if errs := utils.ValidateStruct(req); len(errs) > 0 {
		return response.Error(c, fiber.StatusUnprocessableEntity, response.ValidationErrors{Errors: errs}.Error(), nil)
	}

	res, err := h.authService.LoginWithIDToken(ctx, req.IDToken)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, err.Error(), nil)
	}
	return response.Success(c, "Login success", res)
}

func (h *AuthHandler) Register(c *fiber.Ctx) error {
	ctx := h.ctxinject.HandlerContext(c)

	var req dto.FirebaseRegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, err.Error(), nil)
	}
	if errs := utils.ValidateStruct(req); len(errs) > 0 {
		return response.Error(c, fiber.StatusUnprocessableEntity, response.ValidationErrors{Errors: errs}.Error(), nil)
	}

	if req.IDToken != "" {
		res, err := h.authService.LoginWithIDToken(ctx, req.IDToken)
		if err != nil {
			return response.Error(c, fiber.StatusUnauthorized, err.Error(), nil)
		}
		return response.Success(c, "Register success", res)
	}

	origin, ok := c.Locals("client_origin").(string)
	if !ok || origin == "" {
		// fallback to local host
		origin = "http://localhost:5173"
	}

	res, err := h.authService.RegisterWithEmail(ctx, &req, origin)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, err.Error(), nil)
	}
	return response.Success(c, "Register success", res)
}

func (h *AuthHandler) ResetPassword(c *fiber.Ctx) error {
	ctx := h.ctxinject.HandlerContext(c)

	var req dto.FirebaseResetPasswordRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, err.Error(), nil)
	}
	if errs := utils.ValidateStruct(req); len(errs) > 0 {
		return response.Error(c, fiber.StatusUnprocessableEntity, response.ValidationErrors{Errors: errs}.Error(), nil)
	}

	origin, ok := c.Locals("client_origin").(string)
	if !ok || origin == "" {
		// fallback to local host
		origin = "http://localhost:5173"
	}

	if err := h.authService.SendResetPassword(ctx, req.Email, origin); err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error(), nil)
	}
	return response.Success(c, "Reset password email sent", nil)
}

func (h *AuthHandler) VerifyEmail(c *fiber.Ctx) error {
	ctx := h.ctxinject.HandlerContext(c)

	var req dto.FirebaseVerifyEmailRequest
	_ = c.BodyParser(&req)
	if req.Email == "" {
		email, _ := c.Locals("email").(string)
		req.Email = email
	}
	if errs := utils.ValidateStruct(req); len(errs) > 0 {
		return response.Error(c, fiber.StatusUnprocessableEntity, response.ValidationErrors{Errors: errs}.Error(), nil)
	}

	origin, ok := c.Locals("client_origin").(string)
	if !ok || origin == "" {
		// fallback to local host
		origin = "http://localhost:5173"
	}

	if err := h.authService.SendVerifyEmail(ctx, req.Email, origin); err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error(), nil)
	}
	return response.Success(c, "Verification email sent", nil)
}

func (h *AuthHandler) ChangePassword(c *fiber.Ctx) error {
	ctx := h.ctxinject.HandlerContext(c)

	firebaseUID, _ := c.Locals("firebase_uid").(string)
	userID, _ := c.Locals("user_id").(string)
	if firebaseUID == "" || userID == "" {
		return response.Error(c, fiber.StatusUnauthorized, "User authentication required", nil)
	}

	var req dto.FirebaseChangePasswordRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, err.Error(), nil)
	}
	if errs := utils.ValidateStruct(req); len(errs) > 0 {
		return response.Error(c, fiber.StatusUnprocessableEntity, response.ValidationErrors{Errors: errs}.Error(), nil)
	}

	if err := h.authService.ChangePassword(ctx, firebaseUID, userID, req.NewPassword); err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error(), nil)
	}
	return response.Success(c, "Password changed", nil)
}
