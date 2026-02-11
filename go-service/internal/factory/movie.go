package factory

import (
	"tubexxi/video-api/internal/dependencies"
	"tubexxi/video-api/internal/handler"
	"tubexxi/video-api/internal/routes"
	"tubexxi/video-api/internal/service"

	"github.com/gofiber/fiber/v2"
)

type MovieFactory struct {
	service *service.MovieService
	handler *handler.MovieHandler
	routes  *routes.MovieRoutes
}

func NewMovieFactory(cont *dependencies.Container, mw *MiddlewareFactory) *MovieFactory {
	service := service.NewMovieService(
		cont.Logger,
		cont.ScraperClient,
	)
	handler := handler.NewMovieHandler(
		mw.ContextMiddleware,
		service,
		mw.RateLimiter,
		cont.Logger,
	)
	return &MovieFactory{
		service: service,
		handler: handler,
		routes: routes.NewMovieRoutes(
			handler,
			mw.ContextMiddleware,
			mw.RateLimiter,
		),
	}
}

func (f *MovieFactory) GetRoutes(router fiber.Router) {
	f.routes.RegisterRoutes(router)
}
