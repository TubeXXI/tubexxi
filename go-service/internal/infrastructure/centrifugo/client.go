package centrifugo

import (
	"context"
	"errors"
	"tubexxi/video-api/config"

	"sync"

	"github.com/centrifugal/gocent/v3"
	"go.uber.org/zap"
)

var (
	CentrifugoStorage *CentrifugoClient
	centrifugoOnce    sync.Once
)

type CentrifugoClient struct {
	client *gocent.Client
	config *config.CentrifugoConfig
	logger *zap.Logger
	isUp   bool
	mu     sync.RWMutex
}

func NewCentrifugoClient(ctx context.Context, cfg *config.CentrifugoConfig, logger *zap.Logger) (*CentrifugoClient, error) {
	var initErr error
	centrifugoOnce.Do(func() {
		if cfg.CentrifugoUrl == "" || cfg.CentrifugoApiKey == "" {
			initErr = errors.New("centrifugo URL and API key are required")
			logger.Error("Centrifugo configuration missing")
			return
		}

		client := gocent.New(gocent.Config{
			Addr: cfg.CentrifugoUrl,
			Key:  cfg.CentrifugoApiKey,
		})

		CentrifugoStorage = &CentrifugoClient{
			client: client,
			config: cfg,
			logger: logger,
			isUp:   true,
		}

		logger.Info("✅ Centrifugo client initialized successfully",
			zap.String("url", cfg.CentrifugoUrl),
		)
	})

	if initErr != nil {
		return nil, initErr
	}
	return CentrifugoStorage, nil
}

/**
 * GetCentrifugo returns the singleton Centrifugo client instance
 * @return {*CentrifugoClient} - The Centrifugo client instance
 * @return {error} - Error if the client is not initialized
 */
func GetCentrifugo() (*CentrifugoClient, error) {
	if CentrifugoStorage == nil {
		return nil, errors.New("centrifugo not initialized: call NewCentrifugoClient first")
	}
	return CentrifugoStorage, nil
}
