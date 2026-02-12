package repository

import (
	"context"
	"errors"
	"fmt"
	"time"
	"tubexxi/video-api/internal/entity"
	"tubexxi/video-api/internal/infrastructure/contextpool"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type RoleRepository interface {
	BaseRepository
	FindByID(ctx context.Context, id uuid.UUID) (*entity.Role, error)
	FindByName(ctx context.Context, name string) (*entity.Role, error)
	FindByLevel(ctx context.Context, level entity.RoleLevel) (*entity.Role, error)
	FindAll(ctx context.Context) ([]entity.Role, error)
}

type roleRepository struct {
	*baseRepository
	logger *zap.Logger
}

func NewRoleRepository(db *pgxpool.Pool, logger *zap.Logger) RoleRepository {
	return &roleRepository{
		baseRepository: NewBaseRepository(
			db,
			logger,
		).(*baseRepository),
	}
}
func (r *roleRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.Role, error) {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	query := `
		SELECT * FROM roles WHERE id = $1 AND deleted_at IS NULL
	`
	args := []interface{}{
		id,
	}
	var role entity.Role
	err := pgxscan.Get(subCtx, r.db, &role, query, args...)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("role not found")
		}
		return nil, fmt.Errorf("failed to find role: %w", err)
	}
	return &role, nil
}
func (r *roleRepository) FindByName(ctx context.Context, name string) (*entity.Role, error) {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	query := `
		SELECT * FROM roles WHERE name = $1 AND deleted_at IS NULL
	`
	args := []interface{}{
		name,
	}
	var role entity.Role
	err := pgxscan.Get(subCtx, r.db, &role, query, args...)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("role not found")
		}
		return nil, fmt.Errorf("failed to find role: %w", err)
	}
	return &role, nil
}
func (r *roleRepository) FindByLevel(ctx context.Context, level entity.RoleLevel) (*entity.Role, error) {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	query := `
		SELECT * FROM roles WHERE level = $1 AND deleted_at IS NULL
	`
	args := []interface{}{
		level,
	}
	var role entity.Role
	err := pgxscan.Get(subCtx, r.db, &role, query, args...)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("role not found")
		}
		return nil, fmt.Errorf("failed to find role: %w", err)
	}
	return &role, nil
}
func (r *roleRepository) FindAll(ctx context.Context) ([]entity.Role, error) {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	query := `
		SELECT * FROM roles WHERE deleted_at IS NULL
	`
	var roles []entity.Role
	err := pgxscan.Select(subCtx, r.db, &roles, query)
	if err != nil {
		return nil, fmt.Errorf("failed to find roles: %w", err)
	}
	return roles, nil
}
