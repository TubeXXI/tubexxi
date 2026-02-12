package factory

import (
	"tubexxi/video-api/internal/dependencies"
	"tubexxi/video-api/internal/middleware"
)

type MiddlewareFactory struct {
	ContextMiddleware  *middleware.ContextMiddleware
	Recovery           *middleware.RecoveryMiddleware
	ScopeMiddleware    *middleware.ScopeMiddleware
	ApiMiddleware      *middleware.ApiMiddleware
	RateLimiter        *middleware.RateLimiterMiddleware
	LoggerMiddleware   *middleware.LoggerMiddleware
	AuthMiddleware     *middleware.AuthMiddleware
	CSRFMiddleware     *middleware.CSRFMiddleware
	AdminMiddleware    *middleware.AdminMiddleware
	PlatformMiddleware *middleware.PlatformMiddleware
}

func NewMiddlewareFactory(cont *dependencies.Container) *MiddlewareFactory {
	ctxinject := middleware.NewContextMiddleware(cont.Logger)
	scope := middleware.NewScopeMiddleware(
		ctxinject,
		cont.Logger,
	)
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
			ctxinject,
			&cont.AppConfig.App,
			cont.Logger,
		),
		LoggerMiddleware: middleware.NewLoggerMiddleware(
			cont.Logger,
		),
		AuthMiddleware: middleware.NewAuthMiddleware(
			cont.Notifier,
			ctxinject,
			cont.FirebaseClient,
			cont.RoleRepo,
			cont.UserRepo,
			cont.Logger,
			!cont.AppConfig.App.IsDevelopment(),
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
		ScopeMiddleware: scope,
		PlatformMiddleware: middleware.NewPlatformMiddleware(
			ctxinject,
			scope,
			cont.Logger,
			cont.CacheHelper,
		),
	}
}
