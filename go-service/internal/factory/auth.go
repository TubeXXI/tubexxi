package factory

import (
	"tubexxi/video-api/internal/dependencies"
	"tubexxi/video-api/internal/handler"
	"tubexxi/video-api/internal/routes"
	"tubexxi/video-api/internal/service"

	"github.com/gofiber/fiber/v2"
)

type AuthFactory struct {
	service *service.AuthService
	handler *handler.AuthHandler
	routes  *routes.AuthRoutes
}

func NewAuthFactory(cont *dependencies.Container, mw *MiddlewareFactory) *AuthFactory {
	svc := service.NewAuthService(
		cont.Logger,
		cont.FirebaseClient,
		cont.UserRepo,
		cont.RoleRepo,
		cont.EmailHelper,
	)
	h := handler.NewAuthHandler(mw.ContextMiddleware, svc)
	r := routes.NewAuthRoutes(h, mw.RateLimiter, mw.AuthMiddleware)
	return &AuthFactory{service: svc, handler: h, routes: r}
}

func (f *AuthFactory) GetRoutes(router fiber.Router) {
	f.routes.RegisterRoutes(router)
}
