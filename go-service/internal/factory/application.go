package factory

import (
	"tubexxi/video-api/internal/dependencies"
	"tubexxi/video-api/internal/handler"
	"tubexxi/video-api/internal/routes"
	"tubexxi/video-api/internal/service"

	"github.com/gofiber/fiber/v2"
)

type ApplicationFactory struct {
	service *service.ApplicationService
	handler *handler.ApplicationHandler
	routes  *routes.ApplicationRoutes
}

func NewApplicationFactory(
	con *dependencies.Container,
	mw *MiddlewareFactory,
) *ApplicationFactory {
	service := service.NewApplicationService(
		con.ApplicationRepo,
		con.RedisClient,
		con.Logger,
	)
	handler := handler.NewApplicationHandler(
		mw.ContextMiddleware,
		service,
	)
	routes := routes.NewApplicationRoutes(
		handler,
		mw.AuthMiddleware,
		mw.AdminMiddleware,
		mw.CSRFMiddleware,
	)
	return &ApplicationFactory{
		service: service,
		handler: handler,
		routes:  routes,
	}
}

func (f *ApplicationFactory) GetRoutes(parent fiber.Router) {
	f.routes.RegisterRoutes(parent)
}
