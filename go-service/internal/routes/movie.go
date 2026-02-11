package routes

import (
	"tubexxi/video-api/internal/handler"
	"tubexxi/video-api/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

type MovieRoutes struct {
	path      string
	handler   *handler.MovieHandler
	ctxinject *middleware.ContextMiddleware
	limiter   *middleware.RateLimiterMiddleware
}

func NewMovieRoutes(
	handler *handler.MovieHandler,
	ctxinject *middleware.ContextMiddleware,
	limiter *middleware.RateLimiterMiddleware,
) *MovieRoutes {
	return &MovieRoutes{
		path:      "/movies",
		handler:   handler,
		ctxinject: ctxinject,
		limiter:   limiter,
	}
}

func (r *MovieRoutes) RegisterRoutes(parent fiber.Router) {
	router := parent.Group(r.path)

	// Public Routes
	router.Get("/home", r.handler.GetHome)
}
