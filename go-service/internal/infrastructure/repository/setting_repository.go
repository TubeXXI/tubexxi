package repository

import (
	"context"
	"fmt"
	"time"
	"tubexxi/video-api/internal/entity"
	"tubexxi/video-api/internal/infrastructure/contextpool"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type SettingRepository interface {
	BaseRepository
	Create(ctx context.Context, setting []entity.Setting) error
	GetAll(ctx context.Context, scope string) ([]entity.Setting, error)
	ListScopes(ctx context.Context) ([]string, error)
	GetByGroup(ctx context.Context, scope string, groupName string) ([]entity.Setting, error)
	GetByKey(ctx context.Context, scope string, key string) (*entity.Setting, error)
	UpdateByKey(ctx context.Context, scope string, key string, value string) error
	UpdateBulk(ctx context.Context, scope string, settings []entity.Setting) error
}

type settingRepository struct {
	*baseRepository
}

func NewSettingRepository(db *pgxpool.Pool, logger *zap.Logger) SettingRepository {
	return &settingRepository{
		baseRepository: NewBaseRepository(
			db,
			logger,
		).(*baseRepository),
	}
}

func (r *settingRepository) Create(ctx context.Context, setting []entity.Setting) error {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	for _, item := range setting {
		if item.Scope == "" {
			item.Scope = "default"
		}
	}

	tx, err := r.db.Begin(subCtx)
	if err != nil {
		return err
	}
	defer tx.Rollback(subCtx)

	query := `INSERT INTO settings (key, scope, value, description, group_name) VALUES ($1, $2, $3, $4, $5) ON CONFLICT (scope, key) DO NOTHING`
	for _, item := range setting {
		_, err = tx.Exec(
			subCtx,
			query,
			item.Key,
			item.Scope,
			item.Value,
			item.Description,
			item.GroupName,
		)
		if err != nil {
			r.logger.Error("[SettingRepository.Create]", zap.Error(err))
			return err
		}
	}
	if err := tx.Commit(subCtx); err != nil {
		return err
	}
	return nil
}

func (r *settingRepository) GetAll(ctx context.Context, scope string) ([]entity.Setting, error) {

	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	if scope == "" {
		scope = "default"
	}

	query := `SELECT id, scope, key, value, description, group_name, created_at, updated_at FROM settings WHERE scope = $1 ORDER BY group_name, key`
	rows, err := r.db.Query(subCtx, query, scope)
	if err != nil {
		r.logger.Error("[SettingRepository.GetAll]", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	var settings []entity.Setting
	for rows.Next() {
		var s entity.Setting
		if err := rows.Scan(&s.ID, &s.Scope, &s.Key, &s.Value, &s.Description, &s.GroupName, &s.CreatedAt, &s.UpdatedAt); err != nil {
			r.logger.Error("[SettingRepository.GetAll]", zap.Error(err))
			return nil, err
		}
		settings = append(settings, s)
	}
	return settings, nil
}
func (r *settingRepository) ListScopes(ctx context.Context) ([]string, error) {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	query := `SELECT DISTINCT scope FROM settings ORDER BY scope`
	rows, err := r.db.Query(subCtx, query)
	if err != nil {
		r.logger.Error("[SettingRepository.ListScopes]", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	var scopes []string
	for rows.Next() {
		var scope string
		if err := rows.Scan(&scope); err != nil {
			r.logger.Error("[SettingRepository.ListScopes]", zap.Error(err))
			return nil, err
		}
		scopes = append(scopes, scope)
	}
	return scopes, nil
}
func (r *settingRepository) GetByGroup(ctx context.Context, scope string, groupName string) ([]entity.Setting, error) {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	if scope == "" {
		scope = "default"
	}

	query := `SELECT id, scope, key, value, description, group_name, created_at, updated_at FROM settings WHERE scope = $1 AND group_name = $2 ORDER BY key`
	rows, err := r.db.Query(subCtx, query, scope, groupName)
	if err != nil {
		r.logger.Error("[SettingRepository.GetByGroup]", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	var settings []entity.Setting
	for rows.Next() {
		var s entity.Setting
		if err := rows.Scan(&s.ID, &s.Scope, &s.Key, &s.Value, &s.Description, &s.GroupName, &s.CreatedAt, &s.UpdatedAt); err != nil {
			r.logger.Error("[SettingRepository.GetByGroup]", zap.Error(err))
			return nil, err
		}
		settings = append(settings, s)
	}
	return settings, nil
}
func (r *settingRepository) GetByKey(ctx context.Context, scope string, key string) (*entity.Setting, error) {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	if scope == "" {
		scope = "default"
	}

	query := `SELECT id, scope, key, value, description, group_name, created_at, updated_at FROM settings WHERE scope = $1 AND key = $2`
	var s entity.Setting
	err := r.db.QueryRow(subCtx, query, scope, key).Scan(&s.ID, &s.Scope, &s.Key, &s.Value, &s.Description, &s.GroupName, &s.CreatedAt, &s.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			r.logger.Error("[SettingRepository.GetByKey]", zap.Error(err))
			return nil, nil
		}
		r.logger.Error("[SettingRepository.GetByKey]", zap.Error(err))
		return nil, err
	}
	return &s, nil
}
func (r *settingRepository) UpdateByKey(ctx context.Context, scope string, key string, value string) error {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	if scope == "" {
		scope = "default"
	}

	query := `UPDATE settings SET value = $1, updated_at = NOW() WHERE scope = $2 AND key = $3`
	_, err := r.db.Exec(subCtx, query, value, scope, key)
	return err
}
func (r *settingRepository) UpdateBulk(ctx context.Context, scope string, settings []entity.Setting) error {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	if scope == "" {
		scope = "default"
	}

	tx, err := r.db.Begin(subCtx)
	if err != nil {
		return err
	}
	defer tx.Rollback(subCtx)

	query := `UPDATE settings SET value = $1, updated_at = NOW() WHERE scope = $2 AND key = $3`
	for _, s := range settings {
		if _, err := tx.Exec(subCtx, query, s.Value, scope, s.Key); err != nil {
			r.logger.Error("[SettingRepository.UpdateBulk]", zap.Error(err))
			return fmt.Errorf("failed to update key %s: %w", s.Key, err)
		}
	}

	return tx.Commit(subCtx)
}
