package factory

import (
	"tubexxi/video-api/internal/dependencies"
	"tubexxi/video-api/internal/handler"
	"tubexxi/video-api/internal/routes"
	"tubexxi/video-api/internal/service"

	"github.com/gofiber/fiber/v2"
)

type UserFactory struct {
	service *service.UserService
	handler *handler.UserHandler
	routes  *routes.UserRoutes
}

func NewUserFactory(cont *dependencies.Container, mw *MiddlewareFactory) *UserFactory {
	service := service.NewUserService(
		cont.UserRepo,
		cont.UserHelper,
		cont.SessionHelper,
		cont.Logger,
		cont.MinioClient,
		cont.AppConfig,
	)
	handler := handler.NewUserHandler(
		mw.ContextMiddleware,
		service,
		mw.RateLimiter,
		cont.Logger,
	)
	return &UserFactory{
		service: service,
		handler: handler,
		routes: routes.NewUserRoutes(
			handler,
			mw.ContextMiddleware,
			mw.RateLimiter,
			mw.AuthMiddleware,
			mw.CSRFMiddleware,
		),
	}
}
func (f *UserFactory) GetRoutes(router fiber.Router) {
	f.routes.RegisterRoutes(router)
}
