package factory

import (
	"tubexxi/video-api/internal/dependencies"
	"tubexxi/video-api/internal/middleware"
)

type MiddlewareFactory struct {
	ContextMiddleware *middleware.ContextMiddleware
	Recovery          *middleware.RecoveryMiddleware
	ApiMiddleware     *middleware.ApiMiddleware
	RateLimiter       *middleware.RateLimiterMiddleware
	LoggerMiddleware  *middleware.LoggerMiddleware
	AuthMiddleware    *middleware.AuthMiddleware
	CSRFMiddleware    *middleware.CSRFMiddleware
	AdminMiddleware   *middleware.AdminMiddleware
}

func NewMiddlewareFactory(cont *dependencies.Container) *MiddlewareFactory {
	ctxinject := middleware.NewContextMiddleware(cont.Logger)
	return &MiddlewareFactory{
		ContextMiddleware: ctxinject,
		Recovery: middleware.NewRecoveryMiddleware(
			ctxinject,
			cont.Logger,
			cont.Notifier,
		),
		RateLimiter: middleware.NewRateLimiterMiddleware(
			ctxinject,
			cont.RedisClient,
			cont.Logger,
		),
		ApiMiddleware: middleware.NewApiMiddleware(
			&cont.AppConfig.App,
			cont.Logger,
		),
		LoggerMiddleware: middleware.NewLoggerMiddleware(
			cont.Logger,
		),
		AuthMiddleware: middleware.NewAuthMiddleware(
			cont.Notifier,
			ctxinject,
			cont.RedisClient,
			cont.SessionHelper,
			cont.Logger,
			cont.AppConfig.JWT.JwtSecret,
		),
		CSRFMiddleware: middleware.NewCSRFMiddleware(
			cont.Notifier,
			ctxinject,
			cont.SessionHelper,
			cont.Logger,
		),
		AdminMiddleware: middleware.NewAdminMiddleware(
			ctxinject,
			cont.Logger,
		),
	}
}
