package routes

import (
	"tubexxi/video-api/internal/handler"
	"tubexxi/video-api/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

type HealthRoutes struct {
	path    string
	handler *handler.HealthHandler
	auth    *middleware.AuthMiddleware
	admin   *middleware.AdminMiddleware
	csrf    *middleware.CSRFMiddleware
}

func NewHealthRoutes(
	handler *handler.HealthHandler,
	auth *middleware.AuthMiddleware,
	admin *middleware.AdminMiddleware,
	csrf *middleware.CSRFMiddleware,
) *HealthRoutes {
	return &HealthRoutes{
		path:    "/health",
		handler: handler,
		auth:    auth,
		admin:   admin,
		csrf:    csrf,
	}
}

func (r *HealthRoutes) RegisterRoutes(parent fiber.Router) {
	router := parent.Group(r.path)

	protected := router.Group("/protected")
	protected.Use(r.auth.FirebaseAuth(), r.admin.Handler())

	protected.Get("/check", r.handler.CheckHealth)
	protected.Get("/log", r.handler.GetLogger)
	protected.Post("/log", r.csrf.CSRFProtect(), r.handler.ClearLogs)
}
