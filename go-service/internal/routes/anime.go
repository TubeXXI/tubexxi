package routes

import (
	"tubexxi/video-api/internal/handler"

	"github.com/gofiber/fiber/v2"
)

type AnimeRoutes struct {
	path    string
	handler *handler.AnimeHandler
}

func NewAnimeRoutes(handler *handler.AnimeHandler) *AnimeRoutes {
	return &AnimeRoutes{path: "/anime", handler: handler}
}

func (r *AnimeRoutes) RegisterRoutes(parent fiber.Router) {
	router := parent.Group(r.path)

	router.Get("/latest", r.handler.GetLatest)
	router.Get("/search", r.handler.Search)
	router.Get("/ongoing", r.handler.GetOngoing)
	router.Get("/genres", r.handler.GetGenres)
	router.Get("/detail", r.handler.GetDetail)
	router.Get("/detail/:slug", r.handler.GetDetail)
	router.Get("/episode", r.handler.GetEpisode)
	router.Get("/episode/:slug", r.handler.GetEpisode)
}

