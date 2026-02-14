package main

import (
	"context"
	"log"
	"time"
	"tubexxi/video-api/internal/dependencies"
	"tubexxi/video-api/internal/entity"
	"tubexxi/video-api/pkg/utils"

	"github.com/google/uuid"
)

func main() {
	ctxTimeout, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	cont, err := dependencies.NewContainer(ctxTimeout)
	if err != nil {
		log.Fatalf("Failed to initialize dependencies: %v", err)
	}
	defer cont.Close()

	role := []entity.Role{
		{
			ID:          uuid.New(),
			Name:        entity.RoleSuperAdmin,
			Slug:        entity.RoleSuperAdminSlug,
			Level:       entity.RoleLevelSuperAdmin,
			Description: entity.NewNullString("Super admin role"),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          uuid.New(),
			Name:        entity.RoleAdmin,
			Slug:        entity.RoleAdminSlug,
			Level:       entity.RoleLevelAdmin,
			Description: entity.NewNullString("Admin role"),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          uuid.New(),
			Name:        entity.RoleUser,
			Slug:        entity.RoleUserSlug,
			Level:       entity.RoleLevelUser,
			Description: entity.NewNullString("User role"),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}

	for _, r := range role {
		query := `
			INSERT INTO roles (id, name, slug, description, level, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7)
			ON CONFLICT (name) DO NOTHING
		`
		_, err = cont.DBPool.Exec(ctxTimeout, query, r.ID, r.Name, r.Slug, r.Description, r.Level, r.CreatedAt, r.UpdatedAt)
		if err != nil {
			log.Fatalf("Failed to seed role %s: %v", r.Name, err)
		} else {
			log.Printf("Successfully seeded role %s", r.Name)
		}
	}
	hash, err := utils.HashPassword("admin1234")
	if err != nil {
		log.Fatalf("Failed to hash password: %v", err)
	}

	adminUser := entity.User{
		ID:           uuid.New(),
		FullName:     "Admin",
		Email:        "id.tubexxi@gmail.com",
		PasswordHash: hash,
		RoleID:       role[0].ID,
		IsActive:     true,
		IsVerified:   true,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	insertAdminQuery := `
		INSERT INTO users (id, email, password_hash, full_name, role_id, is_active, is_verified, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		ON CONFLICT (email) DO NOTHING
	`
	_, err = cont.DBPool.Exec(ctxTimeout, insertAdminQuery,
		adminUser.ID,
		adminUser.Email,
		adminUser.PasswordHash,
		adminUser.FullName,
		adminUser.RoleID,
		adminUser.IsActive,
		adminUser.IsVerified,
		adminUser.CreatedAt,
		adminUser.UpdatedAt,
	)
	if err != nil {
		log.Fatalf("Failed to seed admin user: %v", err)
	}
	log.Printf("Successfully seeded admin user %s", adminUser.Email)

	settingQuery := `
	INSERT INTO settings (key, scope, value, description, group_name) VALUES
	('site_name', 'default', 'TubeXXI', 'Website Name', 'WEBSITE'),
	('site_tagline', 'default', 'Watch Movie, Series, and Anime Full HD Free', 'Website Tagline', 'WEBSITE'),
	('site_description', 'default', 'TubeXXI helps you download videos from any site without watermarks for free in MP4 or MP3 online. Fast, HD quality, and just enter the link. Try it now!', 'Website Meta Description', 'WEBSITE'),
	('site_keywords', 'default', 'TubeXXI, movies, series, imdb, tmdb, anime, watch stream movies', 'Website Meta Keywords', 'WEBSITE'),
	('site_logo', 'default', '', 'URL to Website Logo', 'WEBSITE'),
	('site_favicon', 'default', '', 'URL to Website Favicon', 'WEBSITE'),
	('site_email', 'default', 'support@agcforge.com', 'Contact Email', 'WEBSITE'),
	('site_phone', 'default', '', 'Contact Phone', 'WEBSITE'),
	('site_url', 'default', 'http://localhost:5173', 'Website Public URL', 'WEBSITE'),

	('smtp_enabled', 'default', 'false', 'Enable SMTP', 'EMAIL'),
	('smtp_service', 'default', 'gmail', 'SMTP Service Provider', 'EMAIL'),
	('smtp_host', 'default', 'smtp.gmail.com', 'SMTP Host', 'EMAIL'),
	('smtp_port', 'default', '587', 'SMTP Port', 'EMAIL'),
	('smtp_user', 'default', '', 'SMTP Username', 'EMAIL'),
	('smtp_password', 'default', '', 'SMTP Password', 'EMAIL'),
	('from_email', 'default', 'noreply@example.com', 'From Email Address', 'EMAIL'),
	('from_name', 'default', 'Simontok', 'From Name', 'EMAIL'),

	('api_key', 'default', 'd87eead46e14ec1df9857525be7c7d89c25309f543365200827d08c60e489910', 'TubeXXI API Key', 'SYSTEM'),
	('theme', 'default', 'default', 'Website Theme (light/dark)', 'SYSTEM'),
	('enable_documentation', 'default', 'false', 'Enable Documentation', 'SYSTEM'),
	('maintenance_mode', 'default', 'false', 'Enable Maintenance Mode', 'SYSTEM'),
	('maintenance_message', 'default', 'We are currently performing maintenance. Please check back later.', 'Maintenance Message', 'SYSTEM'),
	('source_logo_favicon', 'default', 'local', 'Source of Logo/Favicon (local/remote)', 'SYSTEM'),
	('histats_tracking_code', 'default', '', 'Histats Tracking Code', 'SYSTEM'),
	('google_analytics_code', 'default', '', 'Google Analytics Code', 'SYSTEM'),
	('play_store_app_url', 'default', '', 'Google Play Store App URL', 'SYSTEM'),
	('app_store_app_url', 'default', '', 'Apple App Store App URL', 'SYSTEM'),

	('enable_monetize', 'default', 'false', 'Enable Monetization Features', 'MONETIZE'),
	('type_monetize', 'default', 'adsense', 'Monetization Type (adsense, revenuecat, adsterra)', 'MONETIZE'),
	('publisher_id', 'default', '', 'Adsense Publisher ID', 'MONETIZE'),
	('enable_popup_ad', 'default', 'false', 'Enable Popup Ad', 'MONETIZE'),
	('enable_socialbar_ad', 'default', 'false', 'Enable Social Bar Ad', 'MONETIZE'),
	('popup_ad_code', 'default', '', 'Popup Ad Code', 'MONETIZE'),
	('socialbar_ad_code', 'default', '', 'Social Bar Ad Code', 'MONETIZE'),
	('auto_ad_code', 'default', '', 'Auto Ad Code', 'MONETIZE'),
	('banner_rectangle_ad_code', 'default', '', 'Banner Rectangle Ad Code', 'MONETIZE'),
	('banner_horizontal_ad_code', 'default', '', 'Banner Horizontal Ad Code', 'MONETIZE'),
	('auto_ad_code', 'default', '', 'Auto Ad Code', 'MONETIZE'),
	('banner_rectangle_ad_code', 'default', '', 'Banner Rectangle Ad Code', 'MONETIZE'),
	('banner_horizontal_ad_code', 'default', '', 'Banner Horizontal Ad Code', 'MONETIZE'),
	('banner_vertical_ad_code', 'default', '', 'Banner Vertical Ad Code', 'MONETIZE'),
	('native_ad_code', 'default', '', 'Native Ad Code', 'MONETIZE'),
	('direct_link_ad_code', 'default', '', 'Direct Link Ad Code', 'MONETIZE')
	ON CONFLICT (scope, key) DO NOTHING;`

	_, err = cont.DBPool.Exec(ctxTimeout, settingQuery)
	if err != nil {
		log.Fatalf("Failed to seed settings: %v", err)
	}
	log.Printf("Successfully seeded settings")

	appQuery := `
    INSERT INTO applications (key, package_name, value, description, group_name)
    VALUES
        ('name', 'com.agcforge.tubexxi', 'TubeXXI', 'TubeXXI Application Name', 'CONFIG'),
        ('api_key', 'com.agcforge.tubexxi', '298736dbe1b8bbfe3cc5b6d34e8f12d9856ba68289be1781252fdd5930596e14', 'TubeXXI API Key', 'CONFIG'),
        ('package_name', 'com.agcforge.tubexxi', 'com.agcforge.tubexxi', 'TubeXXI Package Name', 'CONFIG'),
        ('version', 'com.agcforge.tubexxi', '1.0.0', 'TubeXXI Application Version', 'CONFIG'),
        ('type', 'com.agcforge.tubexxi', 'android', 'TubeXXI Application Type', 'CONFIG'),
        ('store_url', 'com.agcforge.tubexxi', '', 'TubeXXI Application Store URL', 'CONFIG'),
        ('is_active', 'com.agcforge.tubexxi', 'true', 'TubeXXI Application Is Active', 'CONFIG'),
        
        ('enable_monetize', 'com.agcforge.tubexxi', 'false', 'TubeXXI Application Enable Monetize', 'MONETIZE'),
        ('enable_admob', 'com.agcforge.tubexxi', 'false', 'TubeXXI Application Enable Admob', 'MONETIZE'),
        ('enable_unity_ad', 'com.agcforge.tubexxi', 'false', 'TubeXXI Application Enable Unity Ad', 'MONETIZE'),
        ('enable_star_io_ad', 'com.agcforge.tubexxi', 'false', 'TubeXXI Application Enable Star IO Ad', 'MONETIZE'),
        ('enable_in_app_purchase', 'com.agcforge.tubexxi', 'false', 'TubeXXI Application Enable In App Purchase', 'MONETIZE'),
        ('admob_id', 'com.agcforge.tubexxi', '', 'TubeXXI Application Admob ID', 'MONETIZE'),
        ('unity_ad_id', 'com.agcforge.tubexxi', '', 'TubeXXI Application Unity Ad ID', 'MONETIZE'),
        ('star_io_ad_id', 'com.agcforge.tubexxi', '', 'TubeXXI Application Star IO Ad ID', 'MONETIZE'),
        ('admob_auto_ad', 'com.agcforge.tubexxi', '', 'TubeXXI Application Admob Auto Ad', 'MONETIZE'),
        ('admob_banner_ad', 'com.agcforge.tubexxi', '', 'TubeXXI Application Admob Banner Ad', 'MONETIZE'),
        ('admob_interstitial_ad', 'com.agcforge.tubexxi', '', 'TubeXXI Application Admob Interstitial Ad', 'MONETIZE'),
        ('admob_rewarded_ad', 'com.agcforge.tubexxi', '', 'TubeXXI Application Admob Rewarded Ad', 'MONETIZE'),
        ('admob_native_ad', 'com.agcforge.tubexxi', '', 'TubeXXI Application Admob Native Ad', 'MONETIZE'),
        ('unity_banner_ad', 'com.agcforge.tubexxi', '', 'TubeXXI Application Unity Banner Ad', 'MONETIZE'),
        ('unity_interstitial_ad', 'com.agcforge.tubexxi', '', 'TubeXXI Application Unity Interstitial Ad', 'MONETIZE'),
        ('unity_rewarded_ad', 'com.agcforge.tubexxi', '', 'TubeXXI Application Unity Rewarded Ad', 'MONETIZE'),
        ('one_signal_id', 'com.agcforge.tubexxi', '', 'TubeXXI Application One Signal ID', 'MONETIZE')
    ON CONFLICT (package_name, key) DO NOTHING;`

	_, err = cont.DBPool.Exec(ctxTimeout, appQuery)
	if err != nil {
		log.Fatalf("Failed to seed applications: %v", err)
	}
	log.Printf("Successfully seeded applications")
}
