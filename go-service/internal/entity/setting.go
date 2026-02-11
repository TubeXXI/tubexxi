package entity

import (
	"time"

	"github.com/google/uuid"
)

type Setting struct {
	ID          uuid.UUID `json:"id" db:"id"`
	Scope       string    `json:"scope" db:"scope"`
	Key         string    `json:"key" db:"key"`
	Value       string    `json:"value" db:"value"`
	Description string    `json:"description" db:"description"`
	GroupName   string    `json:"group_name" db:"group_name"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

type SettingsResponse struct {
	WEBSITE  SettingWeb      `json:"WEBSITE"`
	EMAIL    SettingEmail    `json:"EMAIL"`
	SYSTEM   SettingSystem   `json:"SYSTEM"`
	MONETIZE SettingMonetize `json:"MONETIZE"`
}

type SettingWeb struct {
	SiteName        string    `json:"site_name"`
	SiteTagline     string    `json:"site_tagline"`
	SiteDescription string    `json:"site_description"`
	SiteKeywords    string    `json:"site_keywords"`
	SiteLogo        string    `json:"site_logo"`
	SiteFavicon     string    `json:"site_favicon"`
	SiteEmail       string    `json:"site_email"`
	SitePhone       string    `json:"site_phone"`
	SiteURL         string    `json:"site_url"`
	SiteCreatedAt   time.Time `json:"site_created_at"`
}

type SettingEmail struct {
	SMTPEnabled  bool   `json:"smtp_enabled"`
	SMTPService  string `json:"smtp_service"`
	SMTPHost     string `json:"smtp_host"`
	SMTPPort     int    `json:"smtp_port"`
	SMTPUser     string `json:"smtp_user"`
	SMTPPassword string `json:"smtp_password"`
	FromEmail    string `json:"from_email"`
	FromName     string `json:"from_name"`
}

type SettingSystem struct {
	ApiKey              string `json:"api_key"`
	EnableDocumentation bool   `json:"enable_documentation"`
	MaintenanceMode     bool   `json:"maintenance_mode"`
	MaintenanceMessage  string `json:"maintenance_message"`
	SourceLogoFavicon   string `json:"source_logo_favicon"` // 'local' | 'remote'
	HistatsTrackingCode string `json:"histats_tracking_code"`
	GoogleAnalyticsCode string `json:"google_analytics_code"`
	PlayStoreAppURL     string `json:"play_store_app_url"`
	AppStoreAppURL      string `json:"app_store_app_url"`
}

type SettingMonetize struct {
	EnableMonetize         bool   `json:"enable_monetize"`
	TypeMonetize           string `json:"type_monetize"` // 'adsense' | 'revenuecat' | 'adsterra'
	PublisherID            string `json:"publisher_id"`
	EnablePopupAd          bool   `json:"enable_popup_ad"`
	EnableSocialbarAd      bool   `json:"enable_socialbar_ad"`
	AutoAdCode             string `json:"auto_ad_code"`
	PopupAdCode            string `json:"popup_ad_code"`
	SocialbarAdCode        string `json:"socialbar_ad_code"`
	BannerRectangleAdCode  string `json:"banner_rectangle_ad_code"`
	BannerHorizontalAdCode string `json:"banner_horizontal_ad_code"`
	BannerVerticalAdCode   string `json:"banner_vertical_ad_code"`
	NativeAdCode           string `json:"native_ad_code"`
	DirectLinkAdCode       string `json:"direct_link_ad_code"`
}

// DTOs
type UpdateSettingsBulkRequest struct {
	Key         string `json:"key" validate:"required"`
	Value       string `json:"value" validate:"required"`
	Description string `json:"description" validate:"omitempty"`
	GroupName   string `json:"group_name" validate:"required"`
}

type UploadFileRequest struct {
	Key string `form:"key" validate:"required,oneof=site_logo site_favicon"`
}
