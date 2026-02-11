package routes

import (
	"tubexxi/video-api/internal/handler"
	"tubexxi/video-api/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

type UserRoutes struct {
	path      string
	handler   *handler.UserHandler
	ctxinject *middleware.ContextMiddleware
	limiter   *middleware.RateLimiterMiddleware
	auth      *middleware.AuthMiddleware
	csrf      *middleware.CSRFMiddleware
}

func NewUserRoutes(
	handler *handler.UserHandler,
	ctxinject *middleware.ContextMiddleware,
	limiter *middleware.RateLimiterMiddleware,
	auth *middleware.AuthMiddleware,
	csrf *middleware.CSRFMiddleware,
) *UserRoutes {
	return &UserRoutes{
		path:      "/user",
		handler:   handler,
		ctxinject: ctxinject,
		limiter:   limiter,
		auth:      auth,
		csrf:      csrf,
	}
}
func (r *UserRoutes) RegisterRoutes(parent fiber.Router) {
	router := parent.Group(r.path)

	protected := router.Group("/protected")
	protected.Use(r.auth.FirebaseAuth())

	protected.Get("/current", r.handler.GetCurrentUser)
	protected.Post("/logout", r.handler.Logout)

	protected.Post("/avatar", r.handler.ChangeAvatar)
	protected.Put("/profile", r.handler.UpdateProfile)
	protected.Put("/password", r.handler.ChangePassword)
	protected.Post("/two-factor/enable", r.handler.EnableTwoFactor)
	protected.Post("/two-factor/verify", r.handler.VerifyTwoFactor)
}
