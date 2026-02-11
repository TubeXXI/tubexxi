package minioclient

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"sync"
	"tubexxi/video-api/config"

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
		if len(cfg.MinioEndpoint) >= 8 && cfg.MinioEndpoint[:8] == "https://" {
			cfg.MinioUseSSL = true
		}

		parsedURL, err := url.Parse(cfg.MinioEndpoint)
		if err == nil && parsedURL.Host != "" {
			if parsedURL.Scheme == "https" {
				cfg.MinioUseSSL = true
			}
			cfg.MinioEndpoint = parsedURL.Host
		}
		if len(cfg.MinioEndpoint) > 0 {
			if len(cfg.MinioEndpoint) >= 8 && cfg.MinioEndpoint[:8] == "https://" {
				cfg.MinioEndpoint = cfg.MinioEndpoint[8:]
			} else if len(cfg.MinioEndpoint) >= 7 && cfg.MinioEndpoint[:7] == "http://" {
				cfg.MinioEndpoint = cfg.MinioEndpoint[7:]
			}
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
