package routes

import (
	"time"
	"tubexxi/video-api/internal/handler"
	"tubexxi/video-api/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

type AuthRoutes struct {
	path    string
	handler *handler.AuthHandler
	limiter *middleware.RateLimiterMiddleware
	auth    *middleware.AuthMiddleware
}

func NewAuthRoutes(handler *handler.AuthHandler, limiter *middleware.RateLimiterMiddleware, auth *middleware.AuthMiddleware) *AuthRoutes {
	return &AuthRoutes{path: "/auth", handler: handler, limiter: limiter, auth: auth}
}

func (r *AuthRoutes) RegisterRoutes(parent fiber.Router) {
	router := parent.Group(r.path)

	router.Post("/login", r.limiter.ProgressDelay("login"), r.limiter.BaseLimiter("login", 10, 1*time.Minute), r.handler.Login)
	router.Post("/register", r.limiter.ProgressDelay("register"), r.limiter.BaseLimiter("register", 10, 1*time.Minute), r.handler.Register)
	router.Post("/reset-password", r.limiter.ProgressDelay("reset_password"), r.limiter.BaseLimiter("reset_password", 5, 5*time.Minute), r.handler.ResetPassword)

	protected := router.Group("/protected")
	protected.Use(r.auth.FirebaseAuth())
	protected.Post("/verify-email", r.limiter.BaseLimiter("verify_email", 3, 5*time.Minute), r.handler.VerifyEmail)
	protected.Post("/change-password", r.limiter.BlockLimiter("change_password", 10, 30*time.Minute), r.handler.ChangePassword)
}
