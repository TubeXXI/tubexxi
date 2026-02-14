package repository

import (
	"context"
	"fmt"
	"strings"
	"time"
	"tubexxi/video-api/internal/dto"
	"tubexxi/video-api/internal/entity"
	"tubexxi/video-api/internal/infrastructure/contextpool"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type ApplicationRepository interface {
	BaseRepository
	Create(ctx context.Context, app []entity.Application) error
	GetAll(ctx context.Context, packageName string) ([]entity.Application, error)
	ListPackageNames(ctx context.Context) ([]string, error)
	GetByGroup(ctx context.Context, packageName string, groupName string) ([]entity.Application, error)
	GetByKey(ctx context.Context, packageName string, key string) (*entity.Application, error)
	UpdateByKey(ctx context.Context, packageName string, key string, value string) error
	UpdateBulk(ctx context.Context, packageName string, settings []entity.Application) error
	FindByPackageName(ctx context.Context, packageName string) ([]entity.Application, error)
	FindByAPIKey(ctx context.Context, apiKey string) (*entity.Application, error)
	Search(ctx context.Context, params dto.QueryParamsRequest) ([]*entity.ApplicationResponse, dto.Pagination, error)
	GetByPackageName(ctx context.Context, packageName string) (*entity.ApplicationResponse, error)
	Delete(ctx context.Context, packageName string) error
	BulkDelete(ctx context.Context, packageNames []string) error
}

type applicationRepository struct {
	*baseRepository
}

func NewApplicationRepository(db *pgxpool.Pool, logger *zap.Logger) ApplicationRepository {
	return &applicationRepository{
		baseRepository: NewBaseRepository(
			db,
			logger,
		).(*baseRepository),
	}
}
func (r *applicationRepository) Create(ctx context.Context, app []entity.Application) error {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	tx, err := r.db.Begin(subCtx)
	if err != nil {
		return err
	}
	defer tx.Rollback(subCtx)

	query := `
		INSERT INTO applications 
			(key, package_name, value, description, group_name, created_at, updated_at) 
		VALUES 
			($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (package_name, key) 
		DO UPDATE SET 
			value = EXCLUDED.value,
			description = EXCLUDED.description,
			package_name = EXCLUDED.package_name,
			group_name = EXCLUDED.group_name,
			updated_at = EXCLUDED.updated_at
	`

	for _, item := range app {
		_, err := tx.Exec(
			subCtx,
			query,
			item.Key,
			item.PackageName,
			item.Value,
			item.Description,
			item.GroupName,
			item.CreatedAt,
			item.UpdatedAt,
		)

		if err != nil {
			r.logger.Error("[ApplicationRepository.Create]", zap.Error(err))
			return fmt.Errorf("Failed to register new application : %w", err)
		}
	}
	return tx.Commit(subCtx)
}

func (r *applicationRepository) GetAll(ctx context.Context, packageName string) ([]entity.Application, error) {

	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	if packageName == "" {
		packageName = "default"
	}

	query := `SELECT id, package_name, key, value, description, group_name, created_at, updated_at FROM applications WHERE package_name = $1 ORDER BY group_name, key`
	rows, err := r.db.Query(subCtx, query, packageName)
	if err != nil {
		r.logger.Error("[ApplicationRepository.GetAll]", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	var applications []entity.Application
	for rows.Next() {
		var a entity.Application
		if err := rows.Scan(&a.ID, &a.PackageName, &a.Key, &a.Value, &a.Description, &a.GroupName, &a.CreatedAt, &a.UpdatedAt); err != nil {
			r.logger.Error("[ApplicationRepository.GetAll]", zap.Error(err))
			return nil, err
		}
		applications = append(applications, a)
	}
	return applications, nil
}

func (r *applicationRepository) ListPackageNames(ctx context.Context) ([]string, error) {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	query := `SELECT DISTINCT package_name FROM applications ORDER BY package_name`
	rows, err := r.db.Query(subCtx, query)
	if err != nil {
		r.logger.Error("[ApplicationRepository.ListPackageNames]", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	var names []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			r.logger.Error("[ApplicationRepository.ListPackageNames]", zap.Error(err))
			return nil, err
		}
		names = append(names, name)
	}
	return names, nil
}
func (r *applicationRepository) GetByGroup(ctx context.Context, packageName string, groupName string) ([]entity.Application, error) {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	if packageName == "" {
		packageName = "default"
	}

	query := `SELECT id, package_name, key, value, description, group_name, created_at, updated_at FROM applications WHERE package_name = $1 AND group_name = $2 ORDER BY key`
	rows, err := r.db.Query(subCtx, query, packageName, groupName)
	if err != nil {
		r.logger.Error("[ApplicationRepository.GetByGroup]", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	var applications []entity.Application
	for rows.Next() {
		var a entity.Application
		if err := rows.Scan(&a.ID, &a.PackageName, &a.Key, &a.Value, &a.Description, &a.GroupName, &a.CreatedAt, &a.UpdatedAt); err != nil {
			r.logger.Error("[ApplicationRepository.GetByGroup]", zap.Error(err))
			return nil, err
		}
		applications = append(applications, a)
	}
	return applications, nil
}
func (r *applicationRepository) GetByKey(ctx context.Context, packageName string, key string) (*entity.Application, error) {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	if packageName == "" {
		packageName = "default"
	}

	query := `SELECT id, package_name, key, value, description, group_name, created_at, updated_at FROM applications WHERE package_name = $1 AND key = $2`
	var a entity.Application
	err := r.db.QueryRow(subCtx, query, packageName, key).Scan(
		&a.ID,
		&a.PackageName,
		&a.Key,
		&a.Value,
		&a.Description,
		&a.GroupName,
		&a.CreatedAt,
		&a.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			r.logger.Error("[ApplicationRepository.GetByKey]", zap.Error(err))
			return nil, nil
		}
		r.logger.Error("[ApplicationRepository.GetByKey]", zap.Error(err))
		return nil, err
	}
	return &a, nil
}
func (r *applicationRepository) UpdateByKey(ctx context.Context, packageName string, key string, value string) error {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	if packageName == "" {
		packageName = "default"
	}

	query := `UPDATE applications SET value = $1, updated_at = NOW() WHERE package_name = $2 AND key = $3`
	_, err := r.db.Exec(subCtx, query, value, packageName, key)
	return err
}
func (r *applicationRepository) UpdateBulk(ctx context.Context, packageName string, settings []entity.Application) error {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	if packageName == "" {
		packageName = "default"
	}

	tx, err := r.db.Begin(subCtx)
	if err != nil {
		return err
	}
	defer tx.Rollback(subCtx)

	query := `UPDATE applications SET value = $1, updated_at = NOW() WHERE package_name = $2 AND key = $3`
	for _, s := range settings {
		if _, err := tx.Exec(subCtx, query, s.Value, packageName, s.Key); err != nil {
			r.logger.Error("[ApplicationRepository.UpdateBulk]", zap.Error(err))
			return fmt.Errorf("failed to update key %s: %w", s.Key, err)
		}
	}

	return tx.Commit(subCtx)
}

func (r *applicationRepository) FindByPackageName(ctx context.Context, packageName string) ([]entity.Application, error) {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	if packageName == "" {
		packageName = "default"
	}

	query := `SELECT id, package_name, key, value, description, group_name, created_at, updated_at FROM applications WHERE package_name = $1 ORDER BY group_name, key`
	rows, err := r.db.Query(subCtx, query, packageName)
	if err != nil {
		r.logger.Error("[ApplicationRepository.FindByPackageName]", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	var applications []entity.Application
	for rows.Next() {
		var a entity.Application
		if err := rows.Scan(&a.ID, &a.PackageName, &a.Key, &a.Value, &a.Description, &a.GroupName, &a.CreatedAt, &a.UpdatedAt); err != nil {
			r.logger.Error("[ApplicationRepository.FindByPackageName]", zap.Error(err))
			return nil, err
		}
		applications = append(applications, a)
	}
	return applications, nil
}

func (r *applicationRepository) FindByAPIKey(ctx context.Context, apiKey string) (*entity.Application, error) {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	if apiKey == "" {
		return nil, fmt.Errorf("apiKey is required")
	}

	query := `SELECT id, package_name, key, value, description, group_name, created_at, updated_at FROM applications WHERE key = 'api_key' AND value = $1 LIMIT 1`
	var a entity.Application
	err := r.db.QueryRow(subCtx, query, apiKey).Scan(
		&a.ID,
		&a.PackageName,
		&a.Key,
		&a.Value,
		&a.Description,
		&a.GroupName,
		&a.CreatedAt,
		&a.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			r.logger.Error("[ApplicationRepository.FindByAPIKey]", zap.Error(err))
			return nil, nil
		}
		r.logger.Error("[ApplicationRepository.FindByAPIKey]", zap.Error(err))
		return nil, err
	}
	return &a, nil
}

func (r *applicationRepository) Search(ctx context.Context, params dto.QueryParamsRequest) ([]*entity.ApplicationResponse, dto.Pagination, error) {

	if params.Page <= 0 {
		params.Page = 1
	}
	if params.Limit <= 0 {
		params.Limit = 10
	}
	if params.Limit > 100 {
		params.Limit = 100
	}

	packageQuery := `
		SELECT DISTINCT package_name 
		FROM applications
		WHERE 1=1
	`
	var countQuery strings.Builder
	countQuery.WriteString(`
		SELECT COUNT(DISTINCT package_name)
		FROM applications
		WHERE 1=1
	`)

	args := []interface{}{}
	argIdx := 1

	if params.Search != "" {
		packageQuery += fmt.Sprintf(" AND (package_name ILIKE $%d OR key ILIKE $%d OR value ILIKE $%d)", argIdx, argIdx, argIdx)
		countQuery.WriteString(fmt.Sprintf(" AND (package_name ILIKE $%d OR key ILIKE $%d OR value ILIKE $%d)", argIdx, argIdx, argIdx))
		args = append(args, "%"+params.Search+"%")
		argIdx++
	}

	if !params.DateFrom.IsZero() && !params.DateTo.IsZero() {
		packageQuery += fmt.Sprintf(" AND created_at BETWEEN $%d AND $%d", argIdx, argIdx+1)
		countQuery.WriteString(fmt.Sprintf(" AND created_at BETWEEN $%d AND $%d", argIdx, argIdx+1))
		args = append(args, params.DateFrom, params.DateTo)
		argIdx += 2
	}

	var totalItems int64
	err := r.db.QueryRow(ctx, countQuery.String(), args...).Scan(&totalItems)
	if err != nil {
		return nil, dto.Pagination{}, fmt.Errorf("failed to count applications: %w", err)
	}

	offset := (params.Page - 1) * params.Limit
	packageQuery += fmt.Sprintf(" ORDER BY package_name LIMIT $%d OFFSET $%d", argIdx, argIdx+1)
	args = append(args, params.Limit, offset)

	rows, err := r.db.Query(ctx, packageQuery, args...)
	if err != nil {
		return nil, dto.Pagination{}, fmt.Errorf("failed to query package names: %w", err)
	}
	defer rows.Close()

	var packageNames []string
	for rows.Next() {
		var packageName string
		if err := rows.Scan(&packageName); err != nil {
			return nil, dto.Pagination{}, fmt.Errorf("failed to scan package name: %w", err)
		}
		packageNames = append(packageNames, packageName)
	}

	applications := make([]*entity.ApplicationResponse, 0, len(packageNames))
	for _, packageName := range packageNames {
		app, err := r.GetByPackageName(ctx, packageName)
		if err != nil {
			r.logger.Error("[ApplicationRepository.Search]", zap.Error(err))
			continue
		}
		applications = append(applications, app)
	}

	totalPages := 0
	if params.Limit > 0 {
		totalPages = int((totalItems + int64(params.Limit) - 1) / int64(params.Limit))
	}

	return applications, dto.Pagination{
		CurrentPage: params.Page,
		Limit:       params.Limit,
		TotalItems:  totalItems,
		TotalPages:  totalPages,
		HasNext:     params.Page < totalPages,
		HasPrev:     params.Page > 1,
	}, nil
}

func (r *applicationRepository) GetByPackageName(ctx context.Context, packageName string) (*entity.ApplicationResponse, error) {
	query := `
		SELECT id, package_name, key, value, description, group_name, created_at, updated_at
		FROM applications
		WHERE package_name = $1
		ORDER BY group_name, key
	`

	rows, err := r.db.Query(ctx, query, packageName)
	if err != nil {
		return nil, fmt.Errorf("failed to query application: %w", err)
	}
	defer rows.Close()

	response := &entity.ApplicationResponse{
		CONFIG:   entity.ApplicationConfig{},
		MONETIZE: entity.ApplicationMonetize{},
	}

	response.CONFIG.PackageName = packageName
	response.CONFIG.IsActive = true

	for rows.Next() {
		var app entity.Application
		err := rows.Scan(
			&app.ID, &app.PackageName, &app.Key, &app.Value,
			&app.Description, &app.GroupName, &app.CreatedAt, &app.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan application: %w", err)
		}

		r.mapApplicationToResponse(&app, response)
	}

	return response, nil
}

func (r *applicationRepository) Delete(ctx context.Context, packageName string) error {
	query := `
		DELETE FROM applications
		WHERE package_name = $1
	`
	_, err := r.db.Exec(ctx, query, packageName)
	if err != nil {
		return fmt.Errorf("failed to delete application: %w", err)
	}

	return nil
}

func (r *applicationRepository) BulkDelete(ctx context.Context, packageNames []string) error {
	query := `
		DELETE FROM applications
		WHERE package_name = ANY($1)
	`
	_, err := r.db.Exec(ctx, query, packageNames)
	if err != nil {
		return fmt.Errorf("failed to delete applications: %w", err)
	}

	return nil
}

func (r *applicationRepository) mapApplicationToResponse(app *entity.Application, resp *entity.ApplicationResponse) {
	groupName := app.GroupName
	if groupName == "" {
		switch app.Key {
		case "name", "api_key", "version", "type", "store_url", "is_active", "package_name":
			groupName = "config"
		case "enable_monetize", "enable_admob", "enable_unity_ad", "enable_star_io_ad",
			"admob_id", "unity_ad_id", "star_io_ad_id", "admob_auto_ad", "admob_banner_ad",
			"admob_interstitial_ad", "admob_rewarded_ad", "admob_native_ad",
			"unity_banner_ad", "unity_interstitial_ad", "unity_rewarded_ad", "one_signal_id":
			groupName = "monetize"
		}
	}

	switch groupName {
	case "config", "CONFIG", "Config":
		switch app.Key {
		case "name", "app_name", "application_name":
			resp.CONFIG.Name = app.Value
		case "api_key", "apikey", "api-key":
			resp.CONFIG.ApiKey = app.Value
		case "version", "app_version", "version_code":
			resp.CONFIG.Version = app.Value
		case "type", "app_type", "platform":
			resp.CONFIG.Type = app.Value
		case "store_url", "storeurl", "playstore_url":
			resp.CONFIG.StoreURL = &app.Value
		case "is_active", "active", "enabled":
			resp.CONFIG.IsActive = app.Value == "true" || app.Value == "1" || app.Value == "yes"
		}
		resp.CONFIG.PackageName = app.PackageName

	case "monetize", "MONETIZE", "Monetize", "ads", "ADS":
		switch app.Key {
		case "enable_monetize", "monetize_enabled", "ads_enabled":
			resp.MONETIZE.EnableMonetize = app.Value == "true" || app.Value == "1"
		case "enable_admob", "admob_enabled":
			resp.MONETIZE.EnableAdmob = app.Value == "true" || app.Value == "1"
		case "enable_unity_ad", "unity_enabled":
			resp.MONETIZE.EnableUnityAd = app.Value == "true" || app.Value == "1"
		case "enable_star_io_ad", "star_io_enabled":
			resp.MONETIZE.EnableStarIoAd = app.Value == "true" || app.Value == "1"
		case "enable_in_app_purchase", "in_app_purchase_enabled":
			resp.MONETIZE.EnableInAppPurchase = app.Value == "true" || app.Value == "1"
		case "admob_id", "admob_app_id":
			resp.MONETIZE.AdmobID = &app.Value
		case "unity_ad_id", "unity_game_id":
			resp.MONETIZE.UnityAdID = &app.Value
		case "star_io_ad_id", "star_io_app_id":
			resp.MONETIZE.StarIoAdID = &app.Value
		case "admob_auto_ad":
			resp.MONETIZE.AdmobAutoAd = &app.Value
		case "admob_banner_ad":
			resp.MONETIZE.AdmobBannerAd = &app.Value
		case "admob_interstitial_ad":
			resp.MONETIZE.AdmobInterstitialAd = &app.Value
		case "admob_rewarded_ad":
			resp.MONETIZE.AdmobRewardedAd = &app.Value
		case "admob_native_ad":
			resp.MONETIZE.AdmobNativeAd = &app.Value
		case "unity_banner_ad":
			resp.MONETIZE.UnityBannerAd = &app.Value
		case "unity_interstitial_ad":
			resp.MONETIZE.UnityInterstitialAd = &app.Value
		case "unity_rewarded_ad":
			resp.MONETIZE.UnityRewardedAd = &app.Value
		case "one_signal_id":
			resp.MONETIZE.OneSignalID = &app.Value
		}
	}
}
