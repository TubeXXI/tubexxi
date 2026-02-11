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
	router.Get("/genre/:slug", r.handler.GetMoviesByGenre)
	router.Get("/country/:country", r.handler.GetMoviesByCountry)
	router.Get("/year/:year", r.handler.GetMoviesByYear)
	router.Get("/search", r.handler.SearchMovies)
	router.Get("/special/:page_name", r.handler.GetSpecialPage)
	router.Get("/featured/:type", r.handler.GetMoviesByFeature)
	router.Get("/detail/:slug", r.handler.GetMovieDetail)
	router.Get("/detail", r.handler.GetMovieDetail)
}

