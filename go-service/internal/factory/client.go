package factory

import (
	"tubexxi/video-api/internal/dependencies"
	"tubexxi/video-api/internal/handler"
	"tubexxi/video-api/internal/routes"
	"tubexxi/video-api/internal/service"

	"github.com/gofiber/fiber/v2"
)

type ClientFactory struct {
	service *service.ClientService
	handler *handler.ClientHandler
	routes  *routes.ClientRoutes
}

func NewClientFactory(
	cont *dependencies.Container,
	mw *MiddlewareFactory,
) *ClientFactory {
	s := service.NewClientService(
		cont.EmailHelper,
		cont.Logger,
	)
	h := handler.NewClientHandler(
		mw.ContextMiddleware,
		s,
		cont.Logger,
	)

	return &ClientFactory{
		service: s,
		handler: h,
		routes: routes.NewClientRoutes(
			h,
			mw.RateLimiter,
			mw.CSRFMiddleware,
		),
	}
}

func (f *ClientFactory) GetRoutes(parent fiber.Router) {
	f.routes.RegisterRoutes(parent)
}
