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

type ApplicationRepository interface {
	BaseRepository
	GetAll(ctx context.Context, packageName string) ([]entity.Application, error)
	GetByGroup(ctx context.Context, packageName string, groupName string) ([]entity.Application, error)
	GetByKey(ctx context.Context, packageName string, key string) (*entity.Application, error)
	UpdateByKey(ctx context.Context, packageName string, key string, value string) error
	UpdateBulk(ctx context.Context, packageName string, settings []entity.Application) error
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
