package dependencies

import (
	"context"
	"fmt"
	"strings"
	"time"
	"tubexxi/video-api/config"
	"tubexxi/video-api/database"
	helpers "tubexxi/video-api/internal/helper"
	asynqclient "tubexxi/video-api/internal/infrastructure/asynq-client"
	"tubexxi/video-api/internal/infrastructure/centrifugo"
	"tubexxi/video-api/internal/infrastructure/contextpool"
	firebaseclient "tubexxi/video-api/internal/infrastructure/firebase-client"
	"tubexxi/video-api/internal/infrastructure/metrics"
	minioclient "tubexxi/video-api/internal/infrastructure/minio-client"
	redisclient "tubexxi/video-api/internal/infrastructure/redis-client"
	"tubexxi/video-api/internal/infrastructure/repository"
	scraper_client "tubexxi/video-api/internal/infrastructure/scraper-client"
	"tubexxi/video-api/pkg/logger"
	"tubexxi/video-api/pkg/telegram"

	"go.uber.org/zap"
)

type Container struct {
	AppConfig        *config.Config
	Logger           *zap.Logger
	Notifier         telegram.Notifier
	DBPool           *database.Database
	AppMetrics       *metrics.AppMetrics
	RedisMetrics     *metrics.RedisMetrics
	RedisClient      *redisclient.RedisClient
	AsynqClient      *asynqclient.AsynqClientWrapper
	CentrifugoClient *centrifugo.CentrifugoClient
	MinioClient      *minioclient.MinioClient
	ScraperClient    *scraper_client.ScraperClient
	FirebaseClient   *firebaseclient.FirebaseClient
	RoleRepo         repository.RoleRepository
	UserRepo         repository.UserRepository
	SettingRepo      repository.SettingRepository
	ApplicationRepo  repository.ApplicationRepository
	CacheHelper      *helpers.CacheHelper
	UserHelper       *helpers.UserHelper
	SessionHelper    *helpers.SessionHelper
	EmailHelper      *helpers.MailHelper
}

func NewContainer(ctx context.Context) (*Container, error) {
	init, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	logger := logger.GetLogger(&init.App)

	notifier, err := telegram.NewUnifiedNotifier(5, 100, 3*time.Second, &init.Telegram, logger)
	if err != nil {
		return nil, fmt.Errorf("notifier worker initialization failed: %w", err)
	}

	metrics.InitMetrics()

	// Initialize client
	dbPool, err := database.NewDatabase(ctx, &init.Database, &init.App, notifier, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to create database pool: %w", err)
	}

	redis, err := initRedis(ctx, &init.Redis, metrics.GetRedisMetrics(), notifier, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to create redis client: %w", err)
	}

	centrifugoClient, err := centrifugo.NewCentrifugoClient(ctx, &init.Centrifugo, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to create centrifugo client: %w", err)
	}

	asynqClient, err := initAsynqClient(ctx, &init.Redis, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to create asynq client: %w", err)
	}

	minio, err := minioclient.NewMinioClient(ctx, &init.MinIO, logger)
	if err != nil {
		return nil, fmt.Errorf("minio initialization failed: %w", err)
	}

	scraperClient, err := scraper_client.NewScraperClient(init.Scraper)
	if err != nil {
		return nil, fmt.Errorf("failed to create scraper client: %w", err)
	}
	firebaseClient, err := firebaseclient.NewFirebaseClient(ctx, logger, &init.App)
	if err != nil {
		return nil, fmt.Errorf("failed to create firebase client: %w", err)
	}

	// Initialize repository
	roleRepo := repository.NewRoleRepository(dbPool.Pool, logger)
	userRepo := repository.NewUserRepository(dbPool.Pool, logger)
	settingRepo := repository.NewSettingRepository(dbPool.Pool, logger)
	applicationRepo := repository.NewApplicationRepository(dbPool.Pool, logger)

	// Initialize helper
	cacheHelper := helpers.NewCacheHelper(logger, redis, applicationRepo, settingRepo)
	if err := cacheHelper.LoadAllClientToCache(ctx); err != nil {
		logger.Error("[CacheHelper.LoadAllClientToCache]", zap.Error(err))
	}
	userHelper := helpers.NewUserHelper(redis, userRepo)
	sessionHelper := helpers.NewSessionHelper(redis, logger)
	emailHelper := helpers.NewMailHelper(userHelper, sessionHelper, &init.App, &init.Email, logger)

	return &Container{
		AppConfig:        init,
		Logger:           logger,
		Notifier:         notifier,
		DBPool:           dbPool,
		AppMetrics:       metrics.GetAppMetrics(),
		RedisMetrics:     metrics.GetRedisMetrics(),
		RedisClient:      redis,
		AsynqClient:      asynqClient,
		CentrifugoClient: centrifugoClient,
		MinioClient:      minio,
		ScraperClient:    scraperClient,
		FirebaseClient:   firebaseClient,
		RoleRepo:         roleRepo,
		UserRepo:         userRepo,
		SettingRepo:      settingRepo,
		ApplicationRepo:  applicationRepo,
		CacheHelper:      cacheHelper,
		UserHelper:       userHelper,
		SessionHelper:    sessionHelper,
		EmailHelper:      emailHelper,
	}, nil

}

func initRedis(
	ctx context.Context,
	cfg *config.RedisConfig,
	redisMetric *metrics.RedisMetrics,
	notifier telegram.Notifier,
	logger *zap.Logger,
) (*redisclient.RedisClient, error) {
	ctx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	_, err := redisclient.NewRedisClient(ctx, cfg, redisMetric, logger)
	if err != nil {
		return nil, fmt.Errorf("redis initialization failed: %w", err)
	}

	instance, err := redisclient.GetRedis()
	if err != nil {
		notifier.SendAlert(telegram.AlertRequest{
			Subject: "Critical Redis connection failure",
			Message: err.Error(),
			Metadata: map[string]interface{}{
				"timestamp": time.Now(),
			},
		})
		return nil, fmt.Errorf("redis connection failed: %w", err)
	}

	if instance.IsClosed() {
		return nil, fmt.Errorf("redis connection is closed")
	}

	return instance, nil
}
func initAsynqClient(
	ctx context.Context,
	cfg *config.RedisConfig,
	logger *zap.Logger,
) (*asynqclient.AsynqClientWrapper, error) {
	ctx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	asynqClient, err := asynqclient.NewAsynqClient(cfg, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to create asynq client: %w", err)
	}

	if asynqClient != nil {
		if err = asynqClient.NewAsynqServer(); err != nil {
			logger.Warn("Asynq server initialization failed", zap.Error(err))
		}
	}

	return asynqClient, nil
}

func initMinio(
	ctx context.Context,
	cfg *config.MinIOConfig,
	logger *zap.Logger,
) (*minioclient.MinioClient, error) {
	ctx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	minio, err := minioclient.NewMinioClient(ctx, cfg, logger)
	if err != nil {
		return nil, fmt.Errorf("minio initialization failed: %w", err)
	} else {
		logger.Info("Connected to MinIO server",
			zap.String("endpoint", cfg.MinioEndpoint),
			zap.String("bucket", cfg.MinioBucketName),
			zap.Bool("ssl", cfg.MinioUseSSL),
		)
		if err := minio.CreateBucket(ctx); err != nil {
			return nil, fmt.Errorf("minio bucket creation failed: %w", err)
		} else {
			logger.Info("MinIO bucket created successfully",
				zap.String("bucket", cfg.MinioBucketName),
			)
		}
	}

	return minio, nil
}

func (cont *Container) Close() error {
	var errs []error

	cont.Logger.Info("🔄 Starting graceful shutdown...")

	if cont.DBPool != nil {
		if err := cont.DBPool.Close(); err != nil {
			cont.Logger.Error("Database shutdown error", zap.Error(err))
			errs = append(errs, fmt.Errorf("database shutdown error: %w", err))
		} else {
			cont.Logger.Info("✅ Database connection closed successfully")
		}
	}
	if cont.RedisClient != nil {
		if err := cont.RedisClient.Close(); err != nil {
			cont.Logger.Error("Redis shutdown error", zap.Error(err))
			errs = append(errs, fmt.Errorf("redis shutdown error: %w", err))
		} else {
			cont.Logger.Info("✅ Redis connection closed successfully")
		}
	}
	if cont.AsynqClient != nil {
		cont.AsynqClient.ShutdownServer()
		if err := cont.AsynqClient.Close(); err != nil {
			cont.Logger.Error("Asynq shutdown error", zap.Error(err))
		} else {
			cont.Logger.Info("Asynq closed successfully")
		}
	}

	if cont.FirebaseClient != nil {
		if err := cont.FirebaseClient.Close(); err != nil {
			cont.Logger.Error("Firebase shutdown error", zap.Error(err))
			errs = append(errs, fmt.Errorf("firebase shutdown error: %w", err))
		} else {
			cont.Logger.Info("✅ Firebase connection closed successfully")
		}
	}
	if cont.Logger != nil {
		if err := cont.Logger.Sync(); err != nil {
			// Handle known harmless errors (e.g., Windows file handle)
			if !isHarmlessSyncError(err) {
				errs = append(errs, fmt.Errorf("logger sync error: %w", err))
			}
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("errors during shutdown: %v", errs)
	}

	cont.Logger.Info("✅ All services shut down gracefully")
	return nil
}
func isHarmlessSyncError(err error) bool {
	errMsg := err.Error()
	return strings.Contains(errMsg, "The handle is invalid") ||
		strings.Contains(errMsg, "invalid argument") ||
		strings.Contains(errMsg, "bad file descriptor")
}
func (cont *Container) GetDB() *database.Database {
	return cont.DBPool
}
func (cont *Container) GetRedis() *redisclient.RedisClient {
	return cont.RedisClient
}
func (cont *Container) GetLogger() *zap.Logger {
	return cont.Logger
}
func (cont *Container) GetConfig() *config.Config {
	return cont.AppConfig
}
