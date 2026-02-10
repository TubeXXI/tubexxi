package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type LoggerMiddleware struct {
	logger *zap.Logger
}

func NewLoggerMiddleware(logger *zap.Logger) *LoggerMiddleware {
	return &LoggerMiddleware{
		logger: logger,
	}
}

func (m *LoggerMiddleware) RequestLogger() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()
		err := c.Next()
		duration := time.Since(start)
		status := c.Response().StatusCode()
		msg := "Request processed"

		fields := []zap.Field{
			zap.String("ip", c.IP()),
			zap.Duration("latency", duration),
			zap.String("user_agent", c.Get("User-Agent")),
			zap.Int("status", status),
			zap.String("method", c.Method()),
			zap.String("path", c.Path()),
		}

		if err != nil || status >= 500 {
			fields = append(fields, zap.Error(err))
			m.logger.Error(msg, fields...)
		} else if status >= 400 {
			m.logger.Warn(msg, fields...)
		} else {
			m.logger.Info(msg, fields...)
		}

		return err
	}
}
