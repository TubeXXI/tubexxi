package entity

import (
	"time"

	"github.com/google/uuid"
)

type Application struct {
	ID          uuid.UUID `json:"id" db:"id"`
	PackageName string    `json:"package_name" db:"package_name"`
	Key         string    `json:"key" db:"key"`
	Value       string    `json:"value" db:"value"`
	Description string    `json:"description" db:"description"`
	GroupName   string    `json:"group_name" db:"group_name"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

type ApplicationResponse struct {
	CONFIG   ApplicationConfig   `json:"config"`
	MONETIZE ApplicationMonetize `json:"monetize"`
}

type ApplicationConfig struct {
	Name        string  `json:"name"`
	ApiKey      string  `json:"api_key"`
	PackageName string  `json:"package_name"`
	Version     string  `json:"version"`
	Type        string  `json:"type"` // android or ios
	StoreURL    *string `json:"store_url,omitempty"`
	IsActive    bool    `json:"is_active"`
}
type ApplicationMonetize struct {
	EnableMonetize      bool    `json:"enable_monetize"`
	EnableAdmob         bool    `json:"enable_admob"`
	EnableUnityAd       bool    `json:"enable_unity_ad"`
	EnableStarIoAd      bool    `json:"enable_star_io_ad"`
	AdmobID             *string `json:"admob_id,omitempty"`
	UnityAdID           *string `json:"unity_ad_id,omitempty"`
	StarIoAdID          *string `json:"star_io_ad_id,omitempty"`
	AdmobAutoAd         *string `json:"admob_auto_ad,omitempty"`
	AdmobBannerAd       *string `json:"admob_banner_ad,omitempty"`
	AdmobInterstitialAd *string `json:"admob_interstitial_ad,omitempty"`
	AdmobRewardedAd     *string `json:"admob_rewarded_ad,omitempty"`
	AdmobNativeAd       *string `json:"admob_native_ad,omitempty"`
	UnityBannerAd       *string `json:"unity_banner_ad,omitempty"`
	UnityInterstitialAd *string `json:"unity_interstitial_ad,omitempty"`
	UnityRewardedAd     *string `json:"unity_rewarded_ad,omitempty"`
}

type UpdateApplicationBulkRequest struct {
	Key         string `json:"key" validate:"required"`
	Value       string `json:"value" validate:"required"`
	Description string `json:"description" validate:"omitempty"`
	GroupName   string `json:"group_name" validate:"required"`
}
