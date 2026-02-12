package middleware

import (
	"strings"
	helpers "tubexxi/video-api/internal/helper"
	"tubexxi/video-api/pkg/response"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type PlatformMiddleware struct {
	ctxinject   *ContextMiddleware
	scope       *ScopeMiddleware
	logger      *zap.Logger
	cacheHelper *helpers.CacheHelper
}

var platformMiddlewareExcludedExact = map[string]struct{}{
	"/api/settings/public": {},
	"/metrics":             {},
}

var platformMiddlewareExcludedPrefixes = []string{
	"/.well-known/",
	"/api/applications/public",
	"/api/auth",
	"/api/docs",
	"/api/openapi",
	"/api/swagger",
	"/api/ws",
	"/api/token/csrf",
}

func NewPlatformMiddleware(
	ctxinject *ContextMiddleware,
	scope *ScopeMiddleware,
	logger *zap.Logger,
	cacheHelper *helpers.CacheHelper,
) *PlatformMiddleware {
	return &PlatformMiddleware{
		ctxinject:   ctxinject,
		scope:       scope,
		logger:      logger,
		cacheHelper: cacheHelper,
	}
}
func (m *PlatformMiddleware) ClientPlatformMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		method := c.Method()
		path := c.Path()
		if m.shouldSkipPlatformCheck(method, path) {
			return c.Next()
		}

		ctx := m.ctxinject.From(c)

		apiKey := strings.TrimSpace(c.Get("X-API-Key"))
		if apiKey == "" {
			return response.Error(c, fiber.StatusUnauthorized, "API Key is missing", nil)
		}
		platform := strings.ToLower(strings.TrimSpace(c.Get("X-Platform")))
		if platform == "" {
			return response.Error(c, fiber.StatusUnauthorized, "Platform is missing", nil)
		}

		if platform != "mobile" && platform != "web" {
			return response.Error(c, fiber.StatusUnauthorized, "Invalid platform", nil)
		}

		c.Locals("platform", platform)
		c.Locals("api_key", apiKey)

		if platform == "mobile" {
			app, err := m.cacheHelper.GetAppByAPIKey(ctx, apiKey)
			if err != nil {
				m.logger.Error("failed to validate app api key", zap.Error(err))
				return response.Error(c, fiber.StatusInternalServerError, "Failed to validate API Key", nil)
			}

			if app == nil || !app.IsActive {
				return response.Error(c, fiber.StatusUnauthorized, "Invalid or inactive API Key", nil)
			}

			c.Locals("app_package_name", app.PackageName)
			c.Locals("app_api_key", app.ApiKey)
			c.Locals("app_config", app)

			return c.Next()
		} else {
			scope := strings.TrimSpace(c.Get("X-Scope"))
			if scope == "" {
				scope = m.scope.ResolveSettingsScope(c)
			}

			web, err := m.cacheHelper.GetWebByAPIKey(ctx, apiKey, scope)
			if err != nil {
				m.logger.Error("failed to validate web api key", zap.Error(err), zap.String("scope", scope))
				return response.Error(c, fiber.StatusInternalServerError, "Failed to validate API Key", nil)
			}

			if web == nil || web.ApiKey != apiKey {
				return response.Error(c, fiber.StatusUnauthorized, "Invalid or inactive API Key", nil)
			}
			if web.MaintenanceMode {
				msg := web.MaintenanceMessage
				if msg == "" {
					msg = "Maintenance"
				}
				return response.Error(c, fiber.StatusServiceUnavailable, msg, nil)
			}

			c.Locals("web_api_key", web.ApiKey)
			c.Locals("web_scope", scope)
			c.Locals("web_config", web)

			return c.Next()
		}
	}
}

func (m *PlatformMiddleware) shouldSkipPlatformCheck(method string, path string) bool {
	if method == fiber.MethodOptions {
		return true
	}
	if !strings.HasPrefix(path, "/api") {
		return true
	}
	if _, ok := platformMiddlewareExcludedExact[path]; ok {
		return true
	}
	for _, prefix := range platformMiddlewareExcludedPrefixes {
		if strings.HasPrefix(path, prefix) {
			return true
		}
	}
	return false
}
