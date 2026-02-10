package minioclient

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
	"tubexxi/video-api/internal/infrastructure/contextpool"
	"tubexxi/video-api/pkg/config"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"go.uber.org/zap"
)

var (
	MinioStorage *MinioClient
	minioOnce    sync.Once
)

type MinioClient struct {
	client     *minio.Client
	config     *config.MinIOConfig
	logger     *zap.Logger
	bucketName string
	isUp       bool
	mu         sync.RWMutex
}

func NewMinioClient(ctx context.Context, cfg *config.MinIOConfig, logger *zap.Logger) (*MinioClient, error) {
	var initErr error
	minioOnce.Do(func() {
		if cfg.MinioEndpoint == "" {
			initErr = errors.New("minio endpoint is required")
			logger.Error("MinIO configuration missing")
			return
		}
		if cfg.MinioAccessKeyID == "" || cfg.MinioSecretAccessKey == "" {
			initErr = errors.New("minio access key and secret key are required")
			logger.Error("MinIO configuration missing")
			return
		}
		if cfg.MinioBucketName == "" {
			initErr = errors.New("minio bucket name is required")
			logger.Error("MinIO configuration missing")
			return
		}
		client, err := minio.New(cfg.MinioEndpoint, &minio.Options{
			Creds:  credentials.NewStaticV4(cfg.MinioAccessKeyID, cfg.MinioSecretAccessKey, ""),
			Secure: cfg.MinioUseSSL,
		})
		if err != nil {
			initErr = fmt.Errorf("failed to create MinIO client: %w", err)
			logger.Error("MinIO client creation failed", zap.Error(err))
			return
		}

		subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 10*time.Second)
		defer cancel()

		exists, err := client.BucketExists(subCtx, cfg.MinioBucketName)
		if err != nil {
			initErr = fmt.Errorf("failed to check bucket existence: %w", err)
			logger.Error("MinIO bucket check failed", zap.Error(err))
			return
		}

		if !exists {
			err = client.MakeBucket(subCtx, cfg.MinioBucketName, minio.MakeBucketOptions{})
			if err != nil {
				initErr = fmt.Errorf("failed to create bucket: %w", err)
				logger.Error("MinIO bucket creation failed", zap.Error(err))
				return
			}
			logger.Info("MinIO bucket created", zap.String("bucket", cfg.MinioBucketName))
		}

		MinioStorage = &MinioClient{
			client:     client,
			config:     cfg,
			logger:     logger,
			bucketName: cfg.MinioBucketName,
			isUp:       true,
		}

		logger.Info("✅ MinIO client initialized successfully",
			zap.String("endpoint", cfg.MinioEndpoint),
			zap.String("bucket", cfg.MinioBucketName),
			zap.Bool("ssl", cfg.MinioUseSSL),
		)
	})
	if initErr != nil {
		return nil, initErr
	}
	return MinioStorage, nil
}
func GetMinio() (*MinioClient, error) {
	if MinioStorage == nil {
		return nil, errors.New("minio not initialized: call NewMinioClient first")
	}
	return MinioStorage, nil
}
