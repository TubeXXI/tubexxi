package factory

import (
	"tubexxi/video-api/internal/dependencies"
	"tubexxi/video-api/internal/handler"
	"tubexxi/video-api/internal/routes"
	"tubexxi/video-api/internal/service"

	"github.com/gofiber/fiber/v2"
)

type AnimeFactory struct {
	service *service.AnimeService
	handler *handler.AnimeHandler
	routes  *routes.AnimeRoutes
}

func NewAnimeFactory(cont *dependencies.Container, mw *MiddlewareFactory) *AnimeFactory {
	_ = mw
	svc := service.NewAnimeService(cont.Logger, cont.ScraperClient)
	h := handler.NewAnimeHandler(cont.Logger, svc)
	r := routes.NewAnimeRoutes(h)

	return &AnimeFactory{service: svc, handler: h, routes: r}
}

func (f *AnimeFactory) GetRoutes(router fiber.Router) {
	f.routes.RegisterRoutes(router)
}
