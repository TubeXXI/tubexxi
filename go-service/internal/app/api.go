package app

import (
	"tubexxi/video-api/internal/dependencies"
	"tubexxi/video-api/internal/factory"

	"github.com/gofiber/fiber/v2"
)

func RegisterApiRoutes(router fiber.Router, cont *dependencies.Container, mw *factory.MiddlewareFactory) {

	userFactory := factory.NewUserFactory(cont, mw)
	userFactory.GetRoutes(router)

	movieFactory := factory.NewMovieFactory(cont, mw)
	movieFactory.GetRoutes(router)

}
