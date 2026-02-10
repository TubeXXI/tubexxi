package middleware

import (
	"time"
	"tubexxi/video-api/pkg/config"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type ApiMiddleware struct {
	appConfig *config.AppConfig
	logger    *zap.Logger
}

func NewApiMiddleware(appConfig *config.AppConfig, logger *zap.Logger) *ApiMiddleware {
	return &ApiMiddleware{appConfig: appConfig, logger: logger}
}
func (m *ApiMiddleware) SetupCORS() fiber.Handler {
	return cors.New(cors.Config{
		AllowOriginsFunc: nil,
		AllowOrigins:     m.appConfig.ClientUrl,
		AllowHeaders:     "Origin, Referer, Host, Content-Type, Accept, X-Forwarded-Origin, X-Forwarded-Host, Authorization, X-Client-Platform, X-Package-ID, X-XSRF-TOKEN, X-Xsrf-Token, X-Requested-With, X-Original-Url, X-Forwarded-Referer, X-Real-Host, X-Real-IP, X-Forwarded-For, X-Forwarded-Proto, User-Agent, X-Content-Type-Options, X-Frame-Options, X-XSS-Protection, X-2FA-Session, X-Require-Confirm, X-Platform X-Api-Key",
		AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
		AllowCredentials: false,
		ExposeHeaders:    "Content-Length, X-Request-ID, X-Require-Confirm, X-2FA-Session",
		MaxAge:           86400,
	})
}
func (m *ApiMiddleware) SetupCompression() fiber.Handler {
	return compress.New(compress.Config{
		Level: compress.LevelBestSpeed,
	})
}
func (m *ApiMiddleware) SetupRequestID() fiber.Handler {
	return requestid.New(requestid.Config{
		Header: "X-Request-ID",
		Generator: func() string {
			return uuid.New().String()
		},
	})
}
func (m *ApiMiddleware) SetupMetrics(log *zap.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()
		err := c.Next()
		duration := time.Since(start)
		log.Debug("Request completed",
			zap.String("method", c.Method()),
			zap.String("path", c.Path()),
			zap.Int("status", c.Response().StatusCode()),
			zap.Duration("duration", duration),
			zap.String("request_id", c.GetRespHeader("X-Request-ID")),
		)

		return err
	}
}
