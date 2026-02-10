package middleware

import (
	"time"
	helpers "tubexxi/video-api/internal/helper"
	"tubexxi/video-api/internal/infrastructure/contextpool"
	"tubexxi/video-api/pkg/response"
	"tubexxi/video-api/pkg/telegram"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

var (
	ErrNotFoundCsrfToken    = "Unauthorized - CSRF token not found"
	ErrMismatchCsrfToken    = "Unauthorized - CSRF mismatch"
	ErrExpiredCsrfToken     = "Unauthorized - Invalid or expired csrf token"
	ErrInvalidCsrfAuth      = "Unauthorized - Invalid authorization header"
	ErrInvalidOriginRequest = "Forbidden - Invalid request origin"
)

type CSRFMiddleware struct {
	notifier      telegram.Notifier
	ctxinject     *ContextMiddleware
	sessionHelper *helpers.SessionHelper
	logger        *zap.Logger
}

func NewCSRFMiddleware(notifier telegram.Notifier, ctxinject *ContextMiddleware, sessionHelper *helpers.SessionHelper, logger *zap.Logger) *CSRFMiddleware {
	return &CSRFMiddleware{
		notifier:      notifier,
		ctxinject:     ctxinject,
		sessionHelper: sessionHelper,
		logger:        logger,
	}
}
func (m *CSRFMiddleware) CSRFProtect() fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := m.ctxinject.From(c)
		m.logger.Debug("CSRF Protect", zap.String("path", c.Path()))

		provided := c.Get("X-XSRF-TOKEN")
		if provided == "" {
			m.logger.Error("CSRF token not found", zap.String("path", c.Path()))
			return response.Error(c, fiber.StatusUnauthorized, ErrNotFoundCsrfToken, nil)
		}

		if c.Method() == fiber.MethodGet || c.Method() == fiber.MethodOptions {
			return c.Next()
		}

		subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 60*time.Second)
		defer cancel()

		expected, err := m.sessionHelper.GetCSRFBySession(subCtx, provided)
		if err != nil || expected == "" {
			m.logger.Error("CSRF token not found", zap.String("path", c.Path()))
			return response.Error(c, fiber.StatusUnauthorized, ErrNotFoundCsrfToken, nil)
		}

		return c.Next()
	}
}
