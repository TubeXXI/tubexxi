package routes

import (
	"tubexxi/video-api/internal/handler"
	"tubexxi/video-api/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

type SettingRoutes struct {
	path      string
	handler   *handler.SettingHandler
	ctxinject *middleware.ContextMiddleware
	auth      *middleware.AuthMiddleware
	csrf      *middleware.CSRFMiddleware
}

func NewSettingRoutes(
	handler *handler.SettingHandler,
	ctxinject *middleware.ContextMiddleware,
	auth *middleware.AuthMiddleware,
	csrf *middleware.CSRFMiddleware,
) *SettingRoutes {
	return &SettingRoutes{
		path:      "/settings",
		handler:   handler,
		ctxinject: ctxinject,
		auth:      auth,
		csrf:      csrf,
	}
}

func (r *SettingRoutes) RegisterRoutes(parent fiber.Router) {
	router := parent.Group(r.path)

	router.Get("/public", r.handler.GetPublicSettings)
	router.Get("/all", r.handler.GetAllSettings)

	protected := router.Group("/protected")
	protected.Use(r.auth)
	protected.Post("/update", r.csrf.CSRFProtect(), r.handler.UpdateSettingsBulk)
	protected.Put("/upload", r.csrf.CSRFProtect(), r.handler.UploadFile)
}
