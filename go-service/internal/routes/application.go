package routes

import (
	"tubexxi/video-api/internal/handler"
	"tubexxi/video-api/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

type ApplicationRoutes struct {
	path    string
	handler *handler.ApplicationHandler
	auth    *middleware.AuthMiddleware
	admin   *middleware.AdminMiddleware
	csrf    *middleware.CSRFMiddleware
}

func NewApplicationRoutes(
	h *handler.ApplicationHandler,
	auth *middleware.AuthMiddleware,
	admin *middleware.AdminMiddleware,
	csrf *middleware.CSRFMiddleware,
) *ApplicationRoutes {
	return &ApplicationRoutes{
		path:    "/applications",
		handler: h,
		auth:    auth,
		admin:   admin,
		csrf:    csrf,
	}
}

func (r *ApplicationRoutes) RegisterRoutes(parent fiber.Router) {
	router := parent.Group(r.path)

	router.Get("/public/:package_name", r.handler.GetPublicAppConfig)

	protected := router.Group("/protected")
	protected.Use(r.auth.FirebaseAuth())

	protected.Post("/", r.csrf.CSRFProtect(), r.handler.RegisterApplication)
	protected.Put("/:package_name", r.csrf.CSRFProtect(), r.handler.UpdateAppConfigBulk)

}
