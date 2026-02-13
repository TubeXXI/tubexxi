package factory

import (
	"tubexxi/video-api/internal/dependencies"
	"tubexxi/video-api/internal/handler"
	"tubexxi/video-api/internal/routes"

	"github.com/gofiber/fiber/v2"
)

type HealthFactory struct {
	handler *handler.HealthHandler
	routes  *routes.HealthRoutes
}

func NewHealthFactory(
	cont *dependencies.Container,
	mw *MiddlewareFactory,
) *HealthFactory {
	h := handler.NewHealthHandler(
		mw.ContextMiddleware,
		cont.DBPool.Pool,
		cont.RedisClient,
	)
	r := routes.NewHealthRoutes(
		h,
		mw.AuthMiddleware,
		mw.AdminMiddleware,
		mw.CSRFMiddleware,
	)
	return &HealthFactory{handler: h, routes: r}
}

func (f HealthFactory) GetRoutes(router fiber.Router) {
	f.routes.RegisterRoutes(router)
}
