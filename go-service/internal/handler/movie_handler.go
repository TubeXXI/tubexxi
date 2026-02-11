package handler

import (
	"strconv"
	"tubexxi/video-api/internal/middleware"
	"tubexxi/video-api/internal/service"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type MovieHandler struct {
	contextMiddleware *middleware.ContextMiddleware
	service           *service.MovieService
	rateLimiter       *middleware.RateLimiterMiddleware
	logger            *zap.Logger
}

func NewMovieHandler(
	contextMiddleware *middleware.ContextMiddleware,
	service *service.MovieService,
	rateLimiter *middleware.RateLimiterMiddleware,
	logger *zap.Logger,
) *MovieHandler {
	return &MovieHandler{
		contextMiddleware: contextMiddleware,
		service:           service,
		rateLimiter:       rateLimiter,
		logger:            logger,
	}
}

func (h *MovieHandler) GetHome(c *fiber.Ctx) error {
	ctx := c.UserContext()
	
	response, err := h.service.GetHome(ctx)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"data": response,
	})
}

func (h *MovieHandler) GetMoviesByGenre(c *fiber.Ctx) error {
	ctx := c.UserContext()
	slug := c.Params("slug")
	page, _ := strconv.Atoi(c.Query("page", "1"))

	response, err := h.service.GetMoviesByGenre(ctx, slug, int32(page))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"data": response,
	})
}

func (h *MovieHandler) SearchMovies(c *fiber.Ctx) error {
	ctx := c.UserContext()
	query := c.Query("s")
	if query == "" {
		// Try 'q' as well
		query = c.Query("q")
	}
	page, _ := strconv.Atoi(c.Query("page", "1"))

	response, err := h.service.SearchMovies(ctx, query, int32(page))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"data": response,
	})
}

func (h *MovieHandler) GetMoviesByFeature(c *fiber.Ctx) error {
	ctx := c.UserContext()
	featureType := c.Params("type")
	page, _ := strconv.Atoi(c.Query("page", "1"))

	response, err := h.service.GetMoviesByFeature(ctx, featureType, int32(page))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"data": response,
	})
}

func (h *MovieHandler) GetMoviesByCountry(c *fiber.Ctx) error {
	ctx := c.UserContext()
	country := c.Params("country")
	page, _ := strconv.Atoi(c.Query("page", "1"))

	response, err := h.service.GetMoviesByCountry(ctx, country, int32(page))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"data": response,
	})
}

func (h *MovieHandler) GetMoviesByYear(c *fiber.Ctx) error {
	ctx := c.UserContext()
	yearStr := c.Params("year")
	year, err := strconv.Atoi(yearStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid year",
		})
	}
	page, _ := strconv.Atoi(c.Query("page", "1"))

	response, err := h.service.GetMoviesByYear(ctx, int32(year), int32(page))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"data": response,
	})
}

func (h *MovieHandler) GetSpecialPage(c *fiber.Ctx) error {
	ctx := c.UserContext()
	pageName := c.Params("page_name")
	page, _ := strconv.Atoi(c.Query("page", "1"))

	response, err := h.service.GetSpecialPage(ctx, pageName, int32(page))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"data": response,
	})
}

func (h *MovieHandler) GetMovieDetail(c *fiber.Ctx) error {
	ctx := c.UserContext()
	slug := c.Params("slug")
	if slug == "" {
		slug = c.Query("url")
	}
    
    if slug == "" {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "slug or url is required",
        })
    }

	response, err := h.service.GetMovieDetail(ctx, slug)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
    
    if response == nil {
         return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
            "error": "movie not found",
        })
    }

	return c.JSON(fiber.Map{
		"data": response,
	})
}

