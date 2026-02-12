package middleware

import (
	"context"
	"encoding/json"
	"strings"
	"time"
	"tubexxi/video-api/config"
	"tubexxi/video-api/internal/entity"
	"tubexxi/video-api/internal/infrastructure/contextpool"
	redisclient "tubexxi/video-api/internal/infrastructure/redis-client"
	"tubexxi/video-api/internal/infrastructure/repository"
	"tubexxi/video-api/pkg/response"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

const (
	AppCacheKeyPrefix = "app:apikey:"
	AppCacheTTL       = 24 * time.Hour
)

type ApiMiddleware struct {
	ctxinject       *ContextMiddleware
	appConfig       *config.AppConfig
	logger          *zap.Logger
	redis           *redisclient.RedisClient
	applicationRepo repository.ApplicationRepository
}

func NewApiMiddleware(
	ctxinject *ContextMiddleware,
	appConfig *config.AppConfig,
	logger *zap.Logger,
	redis *redisclient.RedisClient,
	applicationRepo repository.ApplicationRepository,
) *ApiMiddleware {
	return &ApiMiddleware{
		ctxinject:       ctxinject,
		appConfig:       appConfig,
		logger:          logger,
		redis:           redis,
		applicationRepo: applicationRepo,
	}
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

func (m *ApiMiddleware) ApiKeyMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := m.ctxinject.From(c)

		path := c.Path()
		if strings.HasPrefix(path, "/.well-known/") {
			return c.Next()
		}
		if path == "/metrics" || strings.HasPrefix(path, "/api/ws") {
			return c.Next()
		}

		if strings.HasPrefix(path, "/api/v1/mobile-client") && path != "/api/mobile-client/bootstrap" {
			if c.Get("X-Session-Id") != "" {
				return c.Next()
			}
		}

		apiKey := c.Get("X-API-Key")
		if apiKey == "" {
			return response.Error(c, fiber.StatusUnauthorized, "API Key is missing", nil)
		}
		app, err := m.GetAppByAPIKey(ctx, apiKey)
		if err != nil {
			return response.Error(c, fiber.StatusInternalServerError, "Failed to validate API Key", nil)
		}

		if app == nil || !app.IsActive {
			return response.Error(c, fiber.StatusUnauthorized, "Invalid or inactive API Key", nil)
		}

		c.Locals("app_package_name", app.PackageName)
		c.Locals("app_api_key", app.ApiKey)

		return c.Next()
	}
}
func (m *ApiMiddleware) LoadAllAppsToCache(ctx context.Context) error {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	apps, err := m.applicationRepo.GetAll(subCtx, "")
	if err != nil {
		return err
	}

	count := 0
	pipeline := m.redis.Client().Pipeline()

	response := &entity.ApplicationResponse{}

	for _, app := range apps {
		m.mapAppConfigSetting(&response.CONFIG, app)
		appConfig := response.CONFIG
		if !appConfig.IsActive {
			continue
		}

		data, err := json.Marshal(appConfig)
		if err != nil {
			m.logger.Error("failed to marshal app config", zap.Error(err), zap.String("package_name", appConfig.PackageName))
			continue
		}
		pipeline.Set(subCtx, AppCacheKeyPrefix+appConfig.ApiKey, data, AppCacheTTL)

		count++
	}
	return nil
}
func (m *ApiMiddleware) GetAppByAPIKey(ctx context.Context, apiKey string) (*entity.ApplicationConfig, error) {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	data, err := m.redis.Client().Get(subCtx, AppCacheKeyPrefix+apiKey).Bytes()
	if err == nil {
		var app entity.ApplicationConfig
		if err := json.Unmarshal(data, &app); err != nil {
			return &app, err
		}

	}
	app, err := m.applicationRepo.FindByAPIKey(subCtx, apiKey)
	if err != nil {
		return nil, err
	}

	response := &entity.ApplicationResponse{}
	if app != nil {
		m.mapAppConfigSetting(&response.CONFIG, *app)
		_ = m.CacheApp(subCtx, &response.CONFIG)
	}
	return &response.CONFIG, nil
}

func (m *ApiMiddleware) CacheApp(ctx context.Context, app *entity.ApplicationConfig) error {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	data, err := json.Marshal(app)
	if err != nil {
		return err
	}
	return m.redis.Client().Set(subCtx, AppCacheKeyPrefix+app.ApiKey, data, AppCacheTTL).Err()
}
func (m *ApiMiddleware) InvalidateApp(ctx context.Context, apiKey string) error {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	return m.redis.Client().Del(subCtx, AppCacheKeyPrefix+apiKey).Err()
}
func (m *ApiMiddleware) mapAppConfigSetting(target *entity.ApplicationConfig, parent entity.Application) {

	switch parent.Key {
	case "name":
		target.Name = parent.Value
	case "api_key":
		target.ApiKey = parent.Value
	case "package_name":
		target.PackageName = parent.Value
	case "version":
		target.Version = parent.Value
	case "type":
		target.Type = parent.Value
	case "store_url":
		target.StoreURL = &parent.Value
	case "is_active":
		target.IsActive = parent.Value == "1"
	}
}
