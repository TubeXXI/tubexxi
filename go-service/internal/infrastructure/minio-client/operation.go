package minioclient

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/url"
	"time"
	"tubexxi/video-api/internal/infrastructure/contextpool"

	"github.com/minio/minio-go/v7"
	"go.uber.org/zap"
)

func (m *MinioClient) IsUp() bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.isUp
}

func (m *MinioClient) CreateBucket(ctx context.Context) error {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	exists, err := m.client.BucketExists(subCtx, m.bucketName)
	if err != nil {
		if errors.Is(subCtx.Err(), context.DeadlineExceeded) {
			m.logger.Error("MinIO CreateBucket BucketExists timed out", zap.String("bucket", m.bucketName), zap.Error(err))
		}
		return fmt.Errorf("failed to check bucket existence: %w", err)
	}
	if !exists {
		err = m.client.MakeBucket(subCtx, m.bucketName, minio.MakeBucketOptions{})
		if err != nil {
			if errors.Is(subCtx.Err(), context.DeadlineExceeded) {
				m.logger.Error("MinIO CreateBucket MakeBucket timed out", zap.String("bucket", m.bucketName), zap.Error(err))
			}
			return fmt.Errorf("failed to create bucket: %w", err)
		}
		m.logger.Info("Bucket created successfully", zap.String("bucket", m.bucketName))

		policy := fmt.Sprintf(`{
			"Version": "2012-10-17",
			"Statement": [
				{
					"Effect": "Allow",
					"Principal": {
						"AWS": ["*"]
					},
					"Action": ["s3:GetObject"],
					"Resource": ["arn:aws:s3:::%s/*"]
				}
			]
		}`, m.bucketName)

		if err := m.client.SetBucketPolicy(subCtx, m.bucketName, policy); err != nil {
			if errors.Is(subCtx.Err(), context.DeadlineExceeded) {
				m.logger.Error("MinIO CreateBucket SetBucketPolicy timed out", zap.String("bucket", m.bucketName), zap.Error(err))
			}
			m.logger.Error("Failed to set bucket policy", zap.String("bucket", m.bucketName), zap.Error(err))
		} else {
			m.logger.Info("Bucket policy set to public read", zap.String("bucket", m.bucketName))
		}
	} else {
		policy := fmt.Sprintf(`{
			"Version": "2012-10-17",
			"Statement": [
				{
					"Effect": "Allow",
					"Principal": {
						"AWS": ["*"]
					},
					"Action": ["s3:GetObject"],
					"Resource": ["arn:aws:s3:::%s/*"]
				}
			]
		}`, m.bucketName)

		if err := m.client.SetBucketPolicy(subCtx, m.bucketName, policy); err != nil {
			if errors.Is(subCtx.Err(), context.DeadlineExceeded) {
				m.logger.Error("MinIO CreateBucket ensure policy timed out", zap.String("bucket", m.bucketName), zap.Error(err))
			}
			m.logger.Error("Failed to ensure bucket policy", zap.String("bucket", m.bucketName), zap.Error(err))
		}
	}
	return nil
}

func (c *MinioClient) UploadFile(ctx context.Context, bucketName string, objectName string, reader io.Reader, objectSize int64, contentType string) (string, error) {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	info, err := c.client.PutObject(subCtx, bucketName, objectName, reader, objectSize, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		c.logger.Error("MinIO UploadFile PutObject failed", zap.String("bucket", bucketName), zap.String("object", objectName), zap.Error(err))
		return "", fmt.Errorf("failed to upload file: %w", err)
	}

	c.logger.Debug("Successfully uploaded file", zap.String("bucket", bucketName), zap.String("object", objectName), zap.Int64("size", info.Size))

	protocol := "http"
	if c.config.MinioUseSSL {
		protocol = "https"
	}
	return fmt.Sprintf("%s://%s/%s/%s", protocol, c.config.MinioEndpoint, bucketName, objectName), nil
}

func (c *MinioClient) GetFileURL(ctx context.Context, bucketName string, objectName string, expiry time.Duration) (string, error) {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	reqParams := make(url.Values)
	presignedURL, err := c.client.PresignedGetObject(subCtx, bucketName, objectName, expiry, reqParams)
	if err != nil {
		if errors.Is(subCtx.Err(), context.DeadlineExceeded) {
			c.logger.Error("MinIO GetFileURL timed out", zap.String("bucket", bucketName), zap.String("object", objectName), zap.Error(err))
		}
		return "", fmt.Errorf("failed to generate presigned url: %w", err)
	}
	return presignedURL.String(), nil
}

func (c *MinioClient) DeleteFile(ctx context.Context, bucketName string, objectName string) error {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	opts := minio.RemoveObjectOptions{
		GovernanceBypass: true,
	}
	err := c.client.RemoveObject(subCtx, bucketName, objectName, opts)
	if err != nil {
		c.logger.Error("MinIO DeleteFile RemoveObject failed", zap.String("bucket", bucketName), zap.String("object", objectName), zap.Error(err))
		return fmt.Errorf("failed to delete file: %w", err)
	}
	return nil
}

func (c *MinioClient) DeleteFolder(ctx context.Context, bucketName string, prefix string) error {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 30*time.Second)
	defer cancel()

	objectsCh := make(chan minio.ObjectInfo)
	go func() {
		defer close(objectsCh)
		for object := range c.client.ListObjects(subCtx, bucketName, minio.ListObjectsOptions{Prefix: prefix, Recursive: true}) {
			if object.Err != nil {
				c.logger.Error("MinIO DeleteFolder ListObjects failed", zap.String("bucket", bucketName), zap.String("prefix", prefix), zap.Error(object.Err))
				continue
			}
			objectsCh <- object
		}
	}()

	opts := minio.RemoveObjectsOptions{
		GovernanceBypass: true,
	}

	errorCh := c.client.RemoveObjects(subCtx, bucketName, objectsCh, opts)
	for e := range errorCh {
		c.logger.Error("MinIO DeleteFolder RemoveObjects failed", zap.String("bucket", bucketName), zap.String("object", e.ObjectName), zap.Error(e.Err))
	}

	return nil
}

func (m *MinioClient) Close() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.isUp = false
	m.logger.Info("✅ MinIO client closed")
	return nil
}
