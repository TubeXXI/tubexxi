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
	scope     *middleware.ScopeMiddleware
	admin     *middleware.AdminMiddleware
}

func NewSettingRoutes(
	handler *handler.SettingHandler,
	ctxinject *middleware.ContextMiddleware,
	auth *middleware.AuthMiddleware,
	csrf *middleware.CSRFMiddleware,
	scope *middleware.ScopeMiddleware,
	admin *middleware.AdminMiddleware,
) *SettingRoutes {
	return &SettingRoutes{
		path:      "/settings",
		handler:   handler,
		ctxinject: ctxinject,
		auth:      auth,
		csrf:      csrf,
		scope:     scope,
		admin:     admin,
	}
}

func (r *SettingRoutes) RegisterRoutes(parent fiber.Router) {
	router := parent.Group(r.path)
	router.Use(r.scope.SettingsScopeMiddleware())

	router.Get("/public", r.handler.GetPublicSettings)

	protected := router.Group("/protected")
	protected.Post("/register", r.csrf.CSRFProtect(), r.handler.RegisterSetting)
	protected.Use(r.auth.FirebaseAuth(), r.admin.Handler())
	protected.Get("/all", r.handler.GetAllSettings)
	protected.Put("/bulk-update", r.csrf.CSRFProtect(), r.handler.UpdateSettingsBulk)
	protected.Post("/upload", r.csrf.CSRFProtect(), r.handler.UploadFile)
}
