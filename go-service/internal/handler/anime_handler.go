package handler

import (
	"net/url"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"

	"tubexxi/video-api/internal/service"
	"tubexxi/video-api/pkg/response"
)

type AnimeHandler struct {
	logger  *zap.Logger
	service *service.AnimeService
}

func NewAnimeHandler(logger *zap.Logger, svc *service.AnimeService) *AnimeHandler {
	return &AnimeHandler{logger: logger, service: svc}
}

func parseAnimeSlugAndPage(raw string, defaultPage int) (string, int) {
	if raw == "" {
		return raw, defaultPage
	}
	parts := strings.SplitN(raw, "&", 2)
	clean := parts[0]
	page := defaultPage
	if len(parts) == 2 {
		vals, err := url.ParseQuery(parts[1])
		if err == nil {
			if p := vals.Get("page"); p != "" {
				if n, err := strconv.Atoi(p); err == nil {
					page = n
				}
			}
		}
	}
	return clean, page
}

func (h *AnimeHandler) GetLatest(c *fiber.Ctx) error {
	ctx := c.UserContext()
	page, _ := strconv.Atoi(c.Query("page", "1"))
	if page <= 0 {
		page = 1
	}

	result, err := h.service.GetLatest(ctx, int32(page))
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error(), nil)
	}
	return response.SuccessWithMeta(c, "Success fetch anime latest",
		result.Animes,
		result.Pagination,
	)
}

func (h *AnimeHandler) Search(c *fiber.Ctx) error {
	ctx := c.UserContext()
	q := c.Query("s")
	page, _ := strconv.Atoi(c.Query("page", "1"))
	if page <= 0 {
		page = 1
	}

	result, err := h.service.Search(ctx, q, int32(page))
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error(), nil)
	}
	return response.SuccessWithMeta(c, "Success search anime",
		result.Animes,
		result.Pagination,
	)
}

func (h *AnimeHandler) GetOngoing(c *fiber.Ctx) error {
	ctx := c.UserContext()
	page, _ := strconv.Atoi(c.Query("page", "1"))
	if page <= 0 {
		page = 1
	}

	result, err := h.service.GetOngoing(ctx, int32(page))
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error(), nil)
	}
	return response.SuccessWithMeta(c, "Success fetch anime ongoing",
		result.Animes,
		result.Pagination)
}

func (h *AnimeHandler) GetGenres(c *fiber.Ctx) error {
	ctx := c.UserContext()
	result, err := h.service.GetGenres(ctx)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error(), nil)
	}
	return response.Success(c, "Success fetch anime genres", result)
}

func (h *AnimeHandler) GetDetail(c *fiber.Ctx) error {
	ctx := c.UserContext()
	urlStr := c.Query("url")
	if urlStr == "" {
		slug := c.Params("slug")
		if slug != "" {
			slug, _ = parseAnimeSlugAndPage(slug, 1)
			urlStr = "https://otakudesu.best/anime/" + slug + "/"
		}
	}
	if urlStr == "" {
		return response.Error(c, fiber.StatusBadRequest, "url is required", nil)
	}

	result, err := h.service.GetDetail(ctx, urlStr)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error(), nil)
	}
	if result == nil {
		return response.Error(c, fiber.StatusNotFound, "anime not found", nil)
	}
	return response.Success(c, "Success fetch anime detail", result)
}

func (h *AnimeHandler) GetEpisode(c *fiber.Ctx) error {
	ctx := c.UserContext()
	urlStr := c.Query("url")
	if urlStr == "" {
		slug := c.Params("slug")
		if slug != "" {
			slug, _ = parseAnimeSlugAndPage(slug, 1)
			urlStr = "https://otakudesu.best/episode/" + slug + "/"
		}
	}
	if urlStr == "" {
		return response.Error(c, fiber.StatusBadRequest, "url is required", nil)
	}

	result, err := h.service.GetEpisode(ctx, urlStr)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error(), nil)
	}
	if result == nil {
		return response.Error(c, fiber.StatusNotFound, "episode not found", nil)
	}
	return response.Success(c, "Success fetch anime episode", result)
}
