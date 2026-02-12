package helpers

import (
	"context"
	"encoding/json"
	"time"
	"tubexxi/video-api/internal/entity"
	"tubexxi/video-api/internal/infrastructure/contextpool"
	redisclient "tubexxi/video-api/internal/infrastructure/redis-client"
	"tubexxi/video-api/internal/infrastructure/repository"

	"go.uber.org/zap"
)

const (
	ClientCacheTTL = 24 * time.Hour

	clientWebCacheKeyPrefix = "client:web:"
	clientAppCacheKeyPrefix = "client:app:"
)

func webCacheKey(scope string, apiKey string) string {
	if scope == "" {
		scope = "default"
	}
	return clientWebCacheKeyPrefix + "scope:" + scope + ":apikey:" + apiKey
}

func appCacheKey(apiKey string) string {
	return clientAppCacheKeyPrefix + "apikey:" + apiKey
}

type CacheHelper struct {
	logger          *zap.Logger
	redis           *redisclient.RedisClient
	applicationRepo repository.ApplicationRepository
	settingRepo     repository.SettingRepository
}

func NewCacheHelper(
	logger *zap.Logger,
	redis *redisclient.RedisClient,
	applicationRepo repository.ApplicationRepository,
	settingRepo repository.SettingRepository,
) *CacheHelper {
	return &CacheHelper{
		logger:          logger,
		redis:           redis,
		applicationRepo: applicationRepo,
		settingRepo:     settingRepo,
	}
}
func (h *CacheHelper) LoadAllClientToCache(ctx context.Context) error {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	count := 0
	pipeline := h.redis.Client().Pipeline()

	scopes, err := h.settingRepo.ListScopes(subCtx)
	if err != nil {
		return err
	}
	for _, scope := range scopes {
		settings, err := h.settingRepo.GetAll(subCtx, scope)
		if err != nil {
			h.logger.Error("failed to load settings", zap.Error(err), zap.String("scope", scope))
			continue
		}
		var webConfig entity.SettingSystem
		for _, s := range settings {
			h.mapSystemSetting(&webConfig, s)
		}
		if webConfig.ApiKey == "" || webConfig.MaintenanceMode {
			continue
		}

		data, err := json.Marshal(webConfig)
		if err != nil {
			h.logger.Error("failed to marshal web config", zap.Error(err), zap.String("scope", scope))
			continue
		}
		pipeline.Set(subCtx, webCacheKey(scope, webConfig.ApiKey), data, ClientCacheTTL)
		count++
	}

	packageNames, err := h.applicationRepo.ListPackageNames(subCtx)
	if err != nil {
		return err
	}
	for _, packageName := range packageNames {
		apps, err := h.applicationRepo.GetAll(subCtx, packageName)
		if err != nil {
			h.logger.Error("failed to load applications", zap.Error(err), zap.String("package_name", packageName))
			continue
		}

		var appConfig entity.ApplicationConfig
		for _, a := range apps {
			h.mapAppConfigSetting(&appConfig, a)
		}
		if appConfig.ApiKey == "" || !appConfig.IsActive {
			continue
		}

		data, err := json.Marshal(appConfig)
		if err != nil {
			h.logger.Error("failed to marshal app config", zap.Error(err), zap.String("package_name", packageName))
			continue
		}
		pipeline.Set(subCtx, appCacheKey(appConfig.ApiKey), data, ClientCacheTTL)
		count++
	}

	if _, err := pipeline.Exec(subCtx); err != nil {
		return err
	}
	h.logger.Info("client cache loaded", zap.Int("count", count))
	return nil
}
func (h *CacheHelper) GetAppByAPIKey(ctx context.Context, apiKey string) (*entity.ApplicationConfig, error) {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	data, err := h.redis.GetByte(subCtx, appCacheKey(apiKey))
	if err == nil && len(data) > 0 {
		var app entity.ApplicationConfig
		if err := json.Unmarshal(data, &app); err == nil {
			return &app, nil
		}
		_ = h.redis.Client().Del(subCtx, appCacheKey(apiKey)).Err()
	}

	apiKeyRow, err := h.applicationRepo.FindByAPIKey(subCtx, apiKey)
	if err != nil {
		return nil, err
	}
	if apiKeyRow == nil {
		return nil, nil
	}

	apps, err := h.applicationRepo.GetAll(subCtx, apiKeyRow.PackageName)
	if err != nil {
		return nil, err
	}
	var appConfig entity.ApplicationConfig
	for _, a := range apps {
		h.mapAppConfigSetting(&appConfig, a)
	}
	if appConfig.ApiKey == "" {
		return nil, nil
	}
	_ = h.CacheApp(subCtx, &appConfig)
	return &appConfig, nil
}
func (h *CacheHelper) GetWebByAPIKey(ctx context.Context, apiKey, scope string) (*entity.SettingSystem, error) {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	if scope == "" {
		scope = "default"
	}

	cacheKey := webCacheKey(scope, apiKey)
	data, err := h.redis.GetByte(subCtx, cacheKey)
	if err == nil && len(data) > 0 {
		var web entity.SettingSystem
		if err := json.Unmarshal(data, &web); err == nil {
			return &web, nil
		}
		_ = h.redis.Client().Del(subCtx, cacheKey).Err()
	}

	settings, err := h.settingRepo.GetAll(subCtx, scope)
	if err != nil {
		return nil, err
	}

	var webConfig entity.SettingSystem
	for _, s := range settings {
		h.mapSystemSetting(&webConfig, s)
	}
	if webConfig.ApiKey == "" || webConfig.ApiKey != apiKey {
		return nil, nil
	}
	if webConfig.MaintenanceMode {
		return &webConfig, nil
	}
	_ = h.CacheWeb(subCtx, scope, &webConfig)
	return &webConfig, nil
}

func (h *CacheHelper) CacheApp(ctx context.Context, app *entity.ApplicationConfig) error {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	data, err := json.Marshal(app)
	if err != nil {
		return err
	}
	return h.redis.Setbyte(subCtx, appCacheKey(app.ApiKey), data, ClientCacheTTL)
}
func (h *CacheHelper) CacheWeb(ctx context.Context, scope string, web *entity.SettingSystem) error {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	data, err := json.Marshal(web)
	if err != nil {
		return err
	}
	return h.redis.Setbyte(subCtx, webCacheKey(scope, web.ApiKey), data, ClientCacheTTL)
}
func (h *CacheHelper) InvalidateApp(ctx context.Context, apiKey string) error {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	return h.redis.Client().Del(subCtx, appCacheKey(apiKey)).Err()
}
func (h *CacheHelper) InvalidateAllClient(ctx context.Context) error {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	var deletedCount int
	var cursor uint64
	var lastErr error
	pattern := "client:*"

	for {
		var keys []string

		keys, cursor, lastErr = h.redis.Client().Scan(subCtx, cursor, pattern, 100).Result()
		if lastErr != nil {
			break
		}

		if len(keys) > 0 {
			if err := h.redis.Client().Del(subCtx, keys...).Err(); err != nil {
				h.logger.Error("failed to delete cache keys", zap.Error(err), zap.Strings("keys", keys))
				lastErr = err
				break
			}
			deletedCount += len(keys)
		}
		if cursor == 0 {
			break
		}
	}
	if lastErr == nil {
		h.logger.Info("client cache invalidated", zap.Int("deleted", deletedCount))
	}
	return lastErr
}

func (h *CacheHelper) InvalidateWeb(ctx context.Context, scope string, apiKey string) error {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	return h.redis.Client().Del(subCtx, webCacheKey(scope, apiKey)).Err()
}

func (h *CacheHelper) mapAppConfigSetting(target *entity.ApplicationConfig, parent entity.Application) {

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
func (h *CacheHelper) mapSystemSetting(target *entity.SettingSystem, parent entity.Setting) {
	switch parent.Key {
	case "api_key":
		target.ApiKey = parent.Value
	case "theme":
		target.Theme = parent.Value
	case "enable_documentation":
		target.EnableDocumentation = parent.Value == "true"
	case "maintenance_mode":
		target.MaintenanceMode = parent.Value == "true"
	case "maintenance_message":
		target.MaintenanceMessage = parent.Value
	case "source_logo_favicon":
		target.SourceLogoFavicon = parent.Value
	case "histats_tracking_code":
		target.HistatsTrackingCode = parent.Value
	case "google_analytics_code":
		target.GoogleAnalyticsCode = parent.Value
	case "play_store_app_url":
		target.PlayStoreAppURL = parent.Value
	case "app_store_app_url":
		target.AppStoreAppURL = parent.Value
	}
}
