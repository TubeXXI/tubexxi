package service

import (
	"context"
	"fmt"
	"mime/multipart"
	"net/url"
	"strconv"
	"strings"
	"time"
	"tubexxi/video-api/config"
	"tubexxi/video-api/internal/entity"
	"tubexxi/video-api/internal/infrastructure/contextpool"
	minioclient "tubexxi/video-api/internal/infrastructure/minio-client"
	"tubexxi/video-api/internal/infrastructure/repository"
	"tubexxi/video-api/pkg/utils"

	"go.uber.org/zap"
)

type SettingService struct {
	settingRepo repository.SettingRepository
	storage     *minioclient.MinioClient
	cfg         *config.Config
	logger      *zap.Logger
}

func NewSettingService(
	settingRepo repository.SettingRepository,
	storage *minioclient.MinioClient,
	cfg *config.Config,
	logger *zap.Logger) *SettingService {
	return &SettingService{
		settingRepo: settingRepo,
		storage:     storage,
		cfg:         cfg,
		logger:      logger,
	}
}
func (s *SettingService) UploadFile(ctx context.Context, scope string, file *multipart.FileHeader, key string) (string, error) {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	contentType := file.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "image/") {
		return "", fmt.Errorf("invalid file type: %s", contentType)
	}

	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	bucketName := s.cfg.MinIO.MinioBucketName
	objectName := fmt.Sprintf("settings/%s-%s", key, file.Filename)

	if oldSetting, err := s.settingRepo.GetByKey(subCtx, scope, key); err == nil && oldSetting.Value != "" {
		oldObject := oldSetting.Value

		parsedURL, err := url.Parse(oldObject)
		if err == nil {
			path := parsedURL.Path
			path = strings.TrimPrefix(path, "/")

			if strings.HasPrefix(path, s.cfg.MinIO.MinioBucketName+"/") {
				objectName := strings.TrimPrefix(path, s.cfg.MinIO.MinioBucketName+"/")

				if err := s.storage.DeleteFile(subCtx, s.cfg.MinIO.MinioBucketName, objectName); err != nil {
					s.logger.Error("[SettingService.UploadFile]", zap.Error(err), zap.String("object", objectName))
				}
			}
		}
	}

	uploadedPath, err := s.storage.UploadFile(subCtx, bucketName, objectName, src, file.Size, contentType)
	if err != nil {
		return "", err
	}

	if err := s.settingRepo.UpdateByKey(subCtx, scope, key, uploadedPath); err != nil {
		return "", err
	}

	return uploadedPath, nil
}

func (s *SettingService) GetPublicSettings(ctx context.Context, scope string) (*entity.SettingsResponse, error) {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	settings, err := s.settingRepo.GetAll(subCtx, scope)
	if err != nil {
		return nil, err
	}

	response := &entity.SettingsResponse{}

	for _, item := range settings {
		switch item.GroupName {
		case "WEBSITE":
			s.mapWebsiteSetting(&response.WEBSITE, item)
		case "EMAIL":
			s.mapEmailSetting(&response.EMAIL, item)
		case "SYSTEM":
			s.mapSystemSetting(&response.SYSTEM, item)
		case "MONETIZE":
			s.mapMonetizeSetting(&response.MONETIZE, item)
		}
	}

	return response, nil
}

func (s *SettingService) GetAllSettings(ctx context.Context, scope string) ([]entity.Setting, error) {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	return s.settingRepo.GetAll(subCtx, scope)
}

func (s *SettingService) UpdateSetting(ctx context.Context, scope string, key string, value string) error {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()
	return s.settingRepo.UpdateByKey(subCtx, scope, key, value)
}

func (s *SettingService) UpdateSettingsBulk(ctx context.Context, scope string, settings []entity.UpdateSettingsBulkRequest) error {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	payload := make([]entity.Setting, 0, len(settings))
	for _, item := range settings {
		payload = append(payload, entity.Setting{
			Key:         item.Key,
			Value:       item.Value,
			Description: item.Description,
			GroupName:   item.GroupName,
		})
	}
	return s.settingRepo.UpdateBulk(subCtx, scope, payload)
}
func (s *SettingService) mapWebsiteSetting(target *entity.SettingWeb, parent entity.Setting) {
	if target.SiteCreatedAt.IsZero() {
		target.SiteCreatedAt = parent.CreatedAt.In(time.UTC)
	}

	switch parent.Key {
	case "site_name":
		target.SiteName = parent.Value
	case "site_tagline":
		target.SiteTagline = parent.Value
	case "site_description":
		target.SiteDescription = parent.Value
	case "site_keywords":
		target.SiteKeywords = parent.Value
	case "site_logo":
		target.SiteLogo = parent.Value
	case "site_favicon":
		target.SiteFavicon = parent.Value
	case "site_email":
		target.SiteEmail = parent.Value
	case "site_phone":
		target.SitePhone = parent.Value
	case "site_url":
		target.SiteURL = parent.Value
	case "site_created_at":
		target.SiteCreatedAt = parent.CreatedAt.In(time.UTC)
	}
}

func (s *SettingService) mapEmailSetting(target *entity.SettingEmail, parent entity.Setting) {
	switch parent.Key {
	case "smtp_enabled":
		target.SMTPEnabled = parent.Value == "true"
	case "smtp_service":
		target.SMTPService = parent.Value
	case "smtp_host":
		target.SMTPHost = parent.Value
	case "smtp_port":
		if port, err := strconv.Atoi(parent.Value); err == nil {
			target.SMTPPort = port
		}
	case "smtp_user":
		target.SMTPUser = parent.Value
	case "smtp_password":
		hash, err := utils.HashPassword(parent.Value)
		if err != nil {
			s.logger.Error("[SettingService.mapEmailSetting]", zap.Error(err))
		}
		target.SMTPPassword = hash
	case "from_email":
		target.FromEmail = parent.Value
	case "from_name":
		target.FromName = parent.Value
	}
}

func (s *SettingService) mapSystemSetting(target *entity.SettingSystem, parent entity.Setting) {
	switch parent.Key {
	case "api_key":
		target.ApiKey = parent.Value
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
func (s *SettingService) mapMonetizeSetting(target *entity.SettingMonetize, parent entity.Setting) {
	switch parent.Key {
	case "enable_monetize":
		target.EnableMonetize = parent.Value == "true"
	case "type_monetize":
		target.TypeMonetize = parent.Value
	case "publisher_id":
		target.PublisherID = parent.Value
	case "enable_popup_ad":
		target.EnablePopupAd = parent.Value == "true"
	case "enable_socialbar_ad":
		target.EnableSocialbarAd = parent.Value == "true"
	case "auto_ad_code":
		target.AutoAdCode = parent.Value
	case "popup_ad_code":
		target.PopupAdCode = parent.Value
	case "socialbar_ad_code":
		target.SocialbarAdCode = parent.Value
	case "banner_rectangle_ad_code":
		target.BannerRectangleAdCode = parent.Value
	case "banner_horizontal_ad_code":
		target.BannerHorizontalAdCode = parent.Value
	case "banner_vertical_ad_code":
		target.BannerVerticalAdCode = parent.Value
	case "native_ad_code":
		target.NativeAdCode = parent.Value
	case "direct_link_ad_code":
		target.DirectLinkAdCode = parent.Value
	}
}
