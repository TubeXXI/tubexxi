package service

import (
	"context"
	"time"
	"tubexxi/video-api/internal/entity"
	"tubexxi/video-api/internal/infrastructure/contextpool"
	redisclient "tubexxi/video-api/internal/infrastructure/redis-client"
	"tubexxi/video-api/internal/infrastructure/repository"

	"go.uber.org/zap"
)

type ApplicationService struct {
	applicationRepo repository.ApplicationRepository
	redis           *redisclient.RedisClient
	logger          *zap.Logger
}

func NewApplicationService(
	applicationRepo repository.ApplicationRepository,
	redis *redisclient.RedisClient,
	logger *zap.Logger,
) *ApplicationService {
	return &ApplicationService{
		applicationRepo: applicationRepo,
		redis:           redis,
		logger:          logger,
	}
}

func (s *ApplicationService) GetPublicAppConfig(ctx context.Context, packageName string) (*entity.ApplicationResponse, error) {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	settings, err := s.applicationRepo.GetAll(subCtx, packageName)
	if err != nil {
		return nil, err
	}

	response := &entity.ApplicationResponse{}

	for _, item := range settings {
		switch item.GroupName {
		case "CONFIG":
			s.mapAppConfigSetting(&response.CONFIG, item)
		case "MONETIZE":
			s.mapAppMonetizeSetting(&response.MONETIZE, item)
		}
	}

	return response, nil
}
func (s *ApplicationService) GetAllAppConfig(ctx context.Context, packageName string) ([]entity.Application, error) {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	return s.applicationRepo.GetAll(subCtx, packageName)
}
func (s *ApplicationService) UpdateApplication(ctx context.Context, packageName string, key string, value string) error {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()
	return s.applicationRepo.UpdateByKey(subCtx, packageName, key, value)
}
func (s *ApplicationService) UpdateAppConfigBulk(ctx context.Context, packageName string, settings []entity.UpdateApplicationBulkRequest) error {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	payload := make([]entity.Application, 0, len(settings))
	for _, item := range settings {
		payload = append(payload, entity.Application{
			Key:         item.Key,
			Value:       item.Value,
			Description: item.Description,
			GroupName:   item.GroupName,
		})
	}
	return s.applicationRepo.UpdateBulk(subCtx, packageName, payload)
}

func (s *ApplicationService) mapAppConfigSetting(target *entity.ApplicationConfig, parent entity.Application) {

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
func (s *ApplicationService) mapAppMonetizeSetting(target *entity.ApplicationMonetize, parent entity.Application) {

	switch parent.Key {
	case "enable_monetize":
		target.EnableMonetize = parent.Value == "1"
	case "enable_admob":
		target.EnableAdmob = parent.Value == "1"
	case "enable_unity_ad":
		target.EnableUnityAd = parent.Value == "1"
	case "enable_star_io_ad":
		target.EnableStarIoAd = parent.Value == "1"
	case "admob_id":
		target.AdmobID = &parent.Value
	case "unity_ad_id":
		target.UnityAdID = &parent.Value
	case "star_io_ad_id":
		target.StarIoAdID = &parent.Value
	case "admob_auto_ad":
		target.AdmobAutoAd = &parent.Value
	case "admob_banner_ad":
		target.AdmobBannerAd = &parent.Value
	case "admob_interstitial_ad":
		target.AdmobInterstitialAd = &parent.Value
	case "admob_rewarded_ad":
		target.AdmobRewardedAd = &parent.Value
	case "admob_native_ad":
		target.AdmobNativeAd = &parent.Value
	case "unity_banner_ad":
		target.UnityBannerAd = &parent.Value
	case "unity_interstitial_ad":
		target.UnityInterstitialAd = &parent.Value
	case "unity_rewarded_ad":
		target.UnityRewardedAd = &parent.Value
	}
}
