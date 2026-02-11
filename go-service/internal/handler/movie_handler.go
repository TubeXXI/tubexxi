package handler

import (
	"strconv"
	"tubexxi/video-api/internal/middleware"
	"tubexxi/video-api/internal/service"
	"tubexxi/video-api/pkg/response"

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

	result, err := h.service.GetHome(ctx)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error(), nil)
	}

	return response.Success(c, "Success fetch movies",
		result,
	)
}

func (h *MovieHandler) GetMoviesByGenre(c *fiber.Ctx) error {
	ctx := c.UserContext()
	slug := c.Params("slug")
	page, _ := strconv.Atoi(c.Query("page", "1"))

	result, err := h.service.GetMoviesByGenre(ctx, slug, int32(page))
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error(), nil)
	}

	return response.SuccessWithMeta(c, "Success fetch movies",
		result.Movies,
		result.Pagination,
	)
}

func (h *MovieHandler) SearchMovies(c *fiber.Ctx) error {
	ctx := c.UserContext()
	query := c.Query("s")
	if query == "" {
		// Try 'q' as well
		query = c.Query("q")
	}
	page, _ := strconv.Atoi(c.Query("page", "1"))

	result, err := h.service.SearchMovies(ctx, query, int32(page))
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error(), nil)
	}

	return response.SuccessWithMeta(c, "Success fetch movies",
		result.Movies,
		result.Pagination)
}

func (h *MovieHandler) GetMoviesByFeature(c *fiber.Ctx) error {
	ctx := c.UserContext()
	featureType := c.Params("type")
	page, _ := strconv.Atoi(c.Query("page", "1"))

	result, err := h.service.GetMoviesByFeature(ctx, featureType, int32(page))
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error(), nil)
	}

	return response.SuccessWithMeta(c, "Success fetch movies",
		result.Movies,
		result.Pagination,
	)
}

func (h *MovieHandler) GetMoviesByCountry(c *fiber.Ctx) error {
	ctx := c.UserContext()
	country := c.Params("country")
	page, _ := strconv.Atoi(c.Query("page", "1"))

	result, err := h.service.GetMoviesByCountry(ctx, country, int32(page))
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error(), nil)
	}

	return response.SuccessWithMeta(c, "Success fetch movies",
		result.Movies,
		result.Pagination,
	)
}

func (h *MovieHandler) GetMoviesByYear(c *fiber.Ctx) error {
	ctx := c.UserContext()
	yearStr := c.Params("year")
	year, err := strconv.Atoi(yearStr)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "invalid year", nil)
	}
	page, _ := strconv.Atoi(c.Query("page", "1"))

	result, err := h.service.GetMoviesByYear(ctx, int32(year), int32(page))
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error(), nil)
	}

	return response.SuccessWithMeta(c, "Success fetch movies",
		result.Movies,
		result.Pagination,
	)
}

func (h *MovieHandler) GetSpecialPage(c *fiber.Ctx) error {
	ctx := c.UserContext()
	pageName := c.Params("page_name")
	page, _ := strconv.Atoi(c.Query("page", "1"))

	result, err := h.service.GetSpecialPage(ctx, pageName, int32(page))
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error(), nil)
	}

	return response.SuccessWithMeta(c, "Success fetch movies",
		result.Movies,
		result.Pagination,
	)
}

func (h *MovieHandler) GetMovieDetail(c *fiber.Ctx) error {
	ctx := c.UserContext()
	slug := c.Params("slug")
	if slug == "" {
		slug = c.Query("url")
	}

	if slug == "" {
		return response.Error(c, fiber.StatusBadRequest, "slug or url is required", nil)
	}

	result, err := h.service.GetMovieDetail(ctx, slug)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if result == nil {
		return response.Error(c, fiber.StatusNotFound, "movie not found", nil)
	}

	return response.Success(c, "Success fetch movie", result)
}

// Series Handlers

func (h *MovieHandler) GetSeriesHome(c *fiber.Ctx) error {
	ctx := c.UserContext()

	result, err := h.service.GetSeriesHome(ctx)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error(), nil)
	}

	return response.Success(c, "Success fetch series",
		result,
	)
}

func (h *MovieHandler) GetSeriesByGenre(c *fiber.Ctx) error {
	ctx := c.UserContext()
	slug := c.Params("slug")
	page, _ := strconv.Atoi(c.Query("page", "1"))

	result, err := h.service.GetSeriesByGenre(ctx, slug, int32(page))
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error(), nil)
	}

	return response.SuccessWithMeta(c, "Success fetch series",
		result.Movies,
		result.Pagination,
	)
}

func (h *MovieHandler) SearchSeries(c *fiber.Ctx) error {
	ctx := c.UserContext()
	query := c.Query("s")
	if query == "" {
		query = c.Query("q")
	}
	page, _ := strconv.Atoi(c.Query("page", "1"))

	result, err := h.service.SearchSeries(ctx, query, int32(page))
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error(), nil)
	}

	return response.SuccessWithMeta(c, "Success fetch series",
		result.Movies,
		result.Pagination)
}

func (h *MovieHandler) GetSeriesByFeature(c *fiber.Ctx) error {
	ctx := c.UserContext()
	featureType := c.Params("type")
	page, _ := strconv.Atoi(c.Query("page", "1"))

	result, err := h.service.GetSeriesByFeature(ctx, featureType, int32(page))
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error(), nil)
	}

	return response.SuccessWithMeta(c, "Success fetch series",
		result.Movies,
		result.Pagination,
	)
}

func (h *MovieHandler) GetSeriesByCountry(c *fiber.Ctx) error {
	ctx := c.UserContext()
	country := c.Params("country")
	page, _ := strconv.Atoi(c.Query("page", "1"))

	result, err := h.service.GetSeriesByCountry(ctx, country, int32(page))
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error(), nil)
	}

	return response.SuccessWithMeta(c, "Success fetch series",
		result.Movies,
		result.Pagination,
	)
}

func (h *MovieHandler) GetSeriesByYear(c *fiber.Ctx) error {
	ctx := c.UserContext()
	yearStr := c.Params("year")
	year, err := strconv.Atoi(yearStr)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "invalid year", nil)
	}
	page, _ := strconv.Atoi(c.Query("page", "1"))

	result, err := h.service.GetSeriesByYear(ctx, int32(year), int32(page))
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error(), nil)
	}

	return response.SuccessWithMeta(c, "Success fetch series",
		result.Movies,
		result.Pagination,
	)
}

func (h *MovieHandler) GetSeriesSpecialPage(c *fiber.Ctx) error {
	ctx := c.UserContext()
	pageName := c.Params("page_name")
	page, _ := strconv.Atoi(c.Query("page", "1"))

	result, err := h.service.GetSeriesSpecialPage(ctx, pageName, int32(page))
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error(), nil)
	}

	return response.SuccessWithMeta(c, "Success fetch series",
		result.Movies,
		result.Pagination,
	)
}

func (h *MovieHandler) GetSeriesDetail(c *fiber.Ctx) error {
	ctx := c.UserContext()
	slug := c.Params("slug")
	if slug == "" {
		slug = c.Query("url")
	}

	if slug == "" {
		return response.Error(c, fiber.StatusBadRequest, "slug or url is required", nil)
	}

	result, err := h.service.GetSeriesDetail(ctx, slug)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if result == nil {
		return response.Error(c, fiber.StatusNotFound, "series not found", nil)
	}

	return response.Success(c, "Success fetch series", result)
}
