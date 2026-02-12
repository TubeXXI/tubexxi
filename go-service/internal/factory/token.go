package factory

import (
	"tubexxi/video-api/internal/dependencies"
	"tubexxi/video-api/internal/routes"

	"github.com/gofiber/fiber/v2"
)

type TokenFactory struct {
	routes *routes.TokenRoutes
}

func NewTokenFactory(
	cont *dependencies.Container,
	mw *MiddlewareFactory,
) *TokenFactory {
	return &TokenFactory{
		routes: routes.NewTokenRoutes(
			mw.CSRFMiddleware,
		),
	}
}

func (f *TokenFactory) GetRoutes(parent fiber.Router) {
	f.routes.RegisterRoutes(parent)
}
