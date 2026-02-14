package app

import (
	"tubexxi/video-api/internal/dependencies"
	"tubexxi/video-api/internal/factory"

	"github.com/gofiber/fiber/v2"
)

func RegisterApiRoutes(router fiber.Router, cont *dependencies.Container, mw *factory.MiddlewareFactory) {

	tokenFactory := factory.NewTokenFactory(cont, mw)
	tokenFactory.GetRoutes(router)

	settingFactory := factory.NewSettingFactory(cont, mw)
	settingFactory.GetRoutes(router)

	userFactory := factory.NewUserFactory(cont, mw)
	userFactory.GetRoutes(router)

	movieFactory := factory.NewMovieFactory(cont, mw)
	movieFactory.GetRoutes(router)

	authFactory := factory.NewAuthFactory(cont, mw)
	authFactory.GetRoutes(router)

	adminUserFactory := factory.NewAdminUserFactory(cont, mw)
	adminUserFactory.GetRoutes(router)

	animeFactory := factory.NewAnimeFactory(cont, mw)
	animeFactory.GetRoutes(router)

	applicationFactory := factory.NewApplicationFactory(cont, mw)
	applicationFactory.GetRoutes(router)

	healthFactory := factory.NewHealthFactory(cont, mw)
	healthFactory.GetRoutes(router)

	clientFactory := factory.NewClientFactory(cont, mw)
	clientFactory.GetRoutes(router)
}
