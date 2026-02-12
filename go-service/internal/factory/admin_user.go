package factory

import (
	"tubexxi/video-api/internal/dependencies"
	"tubexxi/video-api/internal/handler"
	"tubexxi/video-api/internal/routes"
	"tubexxi/video-api/internal/service"

	"github.com/gofiber/fiber/v2"
)

type AdminUserFactory struct {
	routes *routes.AdminUserRoutes
}

func NewAdminUserFactory(cont *dependencies.Container, mw *MiddlewareFactory) *AdminUserFactory {
	svc := service.NewAdminUserService(cont.Logger, cont.UserRepo, cont.FirebaseClient)
	h := handler.NewAdminUserHandler(mw.ContextMiddleware, svc)
	r := routes.NewAdminUserRoutes(h, mw.AuthMiddleware, mw.AdminMiddleware, mw.RateLimiter)
	return &AdminUserFactory{routes: r}
}

func (f *AdminUserFactory) GetRoutes(router fiber.Router) {
	f.routes.RegisterRoutes(router)
}
