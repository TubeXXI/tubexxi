package factory

import (
	"tubexxi/video-api/internal/dependencies"
	"tubexxi/video-api/internal/handler"
	"tubexxi/video-api/internal/routes"
	"tubexxi/video-api/internal/service"

	"github.com/gofiber/fiber/v2"
)

type SettingFactory struct {
	service *service.SettingService
	handler *handler.SettingHandler
	routes  *routes.SettingRoutes
}

func NewSettingFactory(cont *dependencies.Container, mw *MiddlewareFactory) *SettingFactory {
	service := service.NewSettingService(
		cont.SettingRepo,
		cont.MinioClient,
		cont.AppConfig,
		cont.Logger,
	)
	handler := handler.NewSettingHandler(
		mw.ContextMiddleware,
		mw.ScopeMiddleware,
		service,
		cont.Logger,
	)

	return &SettingFactory{
		service: service,
		handler: handler,
		routes: routes.NewSettingRoutes(
			handler,
			mw.ContextMiddleware,
			mw.AuthMiddleware,
			mw.CSRFMiddleware,
			mw.ScopeMiddleware,
		),
	}
}

func (f *SettingFactory) GetRoutes(router fiber.Router) {
	f.routes.RegisterRoutes(router)
}
