package routes

import (
	"time"
	"tubexxi/video-api/internal/handler"
	"tubexxi/video-api/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

type AdminUserRoutes struct {
	path    string
	h       *handler.AdminUserHandler
	auth    *middleware.AuthMiddleware
	admin   *middleware.AdminMiddleware
	csrf    *middleware.CSRFMiddleware
	limiter *middleware.RateLimiterMiddleware
}

func NewAdminUserRoutes(
	h *handler.AdminUserHandler,
	auth *middleware.AuthMiddleware,
	admin *middleware.AdminMiddleware,
	csrf *middleware.CSRFMiddleware,
	limiter *middleware.RateLimiterMiddleware,
) *AdminUserRoutes {
	return &AdminUserRoutes{
		path:    "/admin/users",
		h:       h,
		auth:    auth,
		admin:   admin,
		csrf:    csrf,
		limiter: limiter,
	}
}

func (r *AdminUserRoutes) RegisterRoutes(parent fiber.Router) {
	router := parent.Group(r.path)
	protected := router.Group("/protected")
	protected.Use(r.auth.FirebaseAuth())
	protected.Use(r.admin.Handler())
	protected.Post("/set-role",
		r.csrf.CSRFProtect(),
		r.limiter.BlockLimiter("admin_set_role", 30, 30*time.Minute),
		r.h.SetRole)
}
