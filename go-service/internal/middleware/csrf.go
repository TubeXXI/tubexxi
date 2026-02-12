package middleware

import (
	"fmt"
	"time"
	"tubexxi/video-api/config"
	redisclient "tubexxi/video-api/internal/infrastructure/redis-client"
	"tubexxi/video-api/pkg/telegram"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
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
	cfg       *config.AppConfig
	notifier  telegram.Notifier
	ctxinject *ContextMiddleware
	redis     *redisclient.RedisClient
	logger    *zap.Logger
}

func NewCSRFMiddleware(
	cfg *config.AppConfig,
	notifier telegram.Notifier,
	ctxinject *ContextMiddleware,
	redis *redisclient.RedisClient,
	logger *zap.Logger,
) *CSRFMiddleware {
	return &CSRFMiddleware{
		cfg:       cfg,
		notifier:  notifier,
		ctxinject: ctxinject,
		redis:     redis,
		logger:    logger,
	}
}
func (m *CSRFMiddleware) CSRFProtect() fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := m.ctxinject.From(c)

		isProd := m.cfg.AppEnv == "production"

		cookieName := "csrf_session_id"

		protocol := c.Protocol()
		if protoHeader := c.Get("X-Forwarded-Proto"); protoHeader != "" {
			protocol = protoHeader
		}

		sameSite := "Lax"
		secureCookie := false

		if isProd && protocol == "https" {
			sameSite = "None"
			secureCookie = true
		}

		sessionID := c.Cookies(cookieName)
		if sessionID == "" {
			sessionID = uuid.New().String()
			c.Cookie(&fiber.Cookie{
				Name:     cookieName,
				Value:    sessionID,
				Expires:  time.Now().Add(24 * time.Hour),
				HTTPOnly: true,
				Secure:   secureCookie,
				SameSite: sameSite,
			})
		}

		redisKey := fmt.Sprintf("csrf_token:%s", sessionID)

		if c.Method() == fiber.MethodGet {
			storedToken, err := m.redis.Client().Get(ctx, redisKey).Result()
			if err == nil && storedToken != "" {
				m.redis.Client().Expire(ctx, redisKey, 1*time.Hour)
				c.Locals("csrf", storedToken)
				return c.Next()
			}

			token := uuid.New().String()

			err = m.redis.Client().Set(ctx, redisKey, token, 1*time.Hour).Err()
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"success": false,
					"message": "Failed to generate CSRF token",
				})
			}

			c.Locals("csrf", token)
			return c.Next()
		}

		clientToken := c.Get("X-XSRF-TOKEN")
		if clientToken == "" {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"success": false,
				"message": "CSRF Token missing",
			})
		}

		storedToken, err := m.redis.Client().Get(ctx, redisKey).Result()
		if err != nil {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"success": false,
				"message": "Invalid or expired CSRF session",
			})
		}

		if clientToken != storedToken {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"success": false,
				"message": "Invalid CSRF Token",
			})
		}

		return c.Next()
	}
}
