package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"
	"tubexxi/video-api/internal/dto"
	"tubexxi/video-api/internal/entity"
	"tubexxi/video-api/internal/infrastructure/contextpool"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type UserRepository interface {
	BaseRepository
	Create(ctx context.Context, user *entity.User) error
	CreateTx(ctx context.Context, tx pgx.Tx, user *entity.User) error
	CreateWithRecovery(ctx context.Context, user *entity.User) error
	FindByID(ctx context.Context, id uuid.UUID) (*entity.User, error)
	FindByEmail(ctx context.Context, email string) (*entity.User, error)
	Update(ctx context.Context, user *entity.User) (*entity.User, error)
	UpdateTx(ctx context.Context, tx pgx.Tx, user *entity.User) (*entity.User, error)
	UpdateWithRecovery(ctx context.Context, user *entity.User) (*entity.User, error)
	UpdateTwoFaSecret(ctx context.Context, id uuid.UUID, twoFaSecret *string) error
	RemoveTwoFaSecret(ctx context.Context, id uuid.UUID) error
	SetEmailVerified(ctx context.Context, id uuid.UUID, isVerified bool) error
	UpdatePassword(ctx context.Context, id uuid.UUID, passwordHash string) error
	UpdateAvatar(ctx context.Context, id uuid.UUID, avatarURL string) (string, error)
	UpdateLastLogin(ctx context.Context, id uuid.UUID) error
	Delete(ctx context.Context, id uuid.UUID) error
	HardDelete(ctx context.Context, id uuid.UUID) error
	Restore(ctx context.Context, id uuid.UUID) error
	ExistsByEmail(ctx context.Context, email string) (bool, error)
	ExistsByPhone(ctx context.Context, phone string) (bool, error)
	BulkDelete(ctx context.Context, ids []uuid.UUID) error
	Search(ctx context.Context, params dto.QueryParamsRequest) ([]*entity.User, dto.Pagination, error)
	FindByRoleID(ctx context.Context, roleID uuid.UUID) ([]*entity.User, error)
	FindRoleByName(ctx context.Context, name string) (*entity.Role, error)
	SetRoleID(ctx context.Context, userID uuid.UUID, roleID uuid.UUID) error
}

type userRepository struct {
	*baseRepository
}

func NewUserRepository(db *pgxpool.Pool, logger *zap.Logger) UserRepository {
	return &userRepository{
		baseRepository: NewBaseRepository(
			db,
			logger,
		).(*baseRepository),
	}
}
func (r *userRepository) Create(ctx context.Context, user *entity.User) error {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	query := `
		INSERT INTO users (
			id, email, password_hash, full_name, role_id, phone, avatar_url,
			is_active, is_verified, created_at
		) VALUES (
			$1, $2, $3, $4,
			COALESCE(NULLIF($5, '00000000-0000-0000-0000-000000000000'::uuid), (SELECT id FROM roles WHERE name = 'user' LIMIT 1)),
			$6, $7, $8, $9, $10
		)
		RETURNING id, created_at, updated_at
	`
	err := r.db.QueryRow(subCtx,
		query,
		user.ID,
		user.Email,
		user.PasswordHash,
		user.FullName,
		user.RoleID,
		user.Phone,
		user.AvatarURL,
		user.IsActive,
		user.IsVerified,
		user.CreatedAt,
	).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		var pgxErr *pgconn.PgError
		if errors.As(err, &pgxErr) && pgxErr.Code == "23505" {
			switch pgxErr.ConstraintName {
			case "users_email_key":
				return fmt.Errorf("email %s is already registered", user.Email)
			case "users_name_length_check":
				return fmt.Errorf("full name %s is invalid, must be between 2 and 50 characters", user.FullName)
			default:
				return fmt.Errorf("unique constraint violation (%s): %w", pgxErr.ConstraintName, err)
			}
		}
		return fmt.Errorf("failed to create new user: %w", err)
	}
	return nil
}
func (r *userRepository) CreateTx(ctx context.Context, tx pgx.Tx, user *entity.User) error {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	query := `
		INSERT INTO users (
			id, email, password_hash, full_name, role_id, phone, avatar_url,
			is_active, is_verified, created_at
		) VALUES (
			$1, $2, $3, $4,
			COALESCE(NULLIF($5, '00000000-0000-0000-0000-000000000000'::uuid), (SELECT id FROM roles WHERE name = 'user' LIMIT 1)),
			$6, $7, $8, $9, $10
		)
		RETURNING id, created_at, updated_at
	`
	err := tx.QueryRow(subCtx,
		query,
		user.ID,
		user.Email,
		user.PasswordHash,
		user.FullName,
		user.RoleID,
		user.Phone,
		user.AvatarURL,
		user.IsActive,
		user.IsVerified,
		user.CreatedAt,
	).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		var pgxErr *pgconn.PgError
		if errors.As(err, &pgxErr) && pgxErr.Code == "23505" {
			switch pgxErr.ConstraintName {
			case "users_email_key":
				return fmt.Errorf("email %s is already registered", user.Email)
			case "users_name_length_check":
				return fmt.Errorf("full name %s is invalid, must be between 2 and 50 characters", user.FullName)
			default:
				return fmt.Errorf("unique constraint violation (%s): %w", pgxErr.ConstraintName, err)
			}
		}
		return fmt.Errorf("failed to create new user: %w", err)
	}
	return nil
}
func (r *userRepository) CreateWithRecovery(ctx context.Context, user *entity.User) error {
	return r.WithTransaction(ctx, func(tx pgx.Tx) error {
		return r.CreateTx(ctx, tx, user)
	})
}
func (r *userRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	query := `
		SELECT 
			u.id, u.email, u.password_hash, u.full_name, u.role_id, u.phone, u.avatar_url,
			u.two_fa_secret, u.is_active, u.is_verified, u.created_at, u.updated_at,
			r.id, r.name, r.level
		FROM users u
		LEFT JOIN roles r ON u.role_id = r.id
		WHERE u.id = $1 AND u.deleted_at IS NULL
	`
	var user entity.User
	var userRoleID uuid.NullUUID
	var roleID uuid.NullUUID
	var roleName sql.NullString
	var roleLevel sql.NullInt32

	err := r.db.QueryRow(subCtx, query, id).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.FullName,
		&userRoleID,
		&user.Phone,
		&user.AvatarURL,
		&user.TwoFaSecret,
		&user.IsActive,
		&user.IsVerified,
		&user.CreatedAt,
		&user.UpdatedAt,
		&roleID,
		&roleName,
		&roleLevel,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to find user by id: %w", err)
	}

	if userRoleID.Valid {
		user.RoleID = userRoleID.UUID
	} else {
		user.RoleID = uuid.Nil
	}

	if roleID.Valid {
		user.Role = &entity.Role{ID: roleID.UUID, Name: roleName.String, Level: entity.RoleLevel(roleLevel.Int32)}
	}
	return &user, nil
}
func (r *userRepository) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	query := `
		SELECT 
			u.id, u.email, u.password_hash, u.full_name, u.role_id, u.phone, u.avatar_url,
			u.two_fa_secret, u.is_active, u.is_verified, u.created_at, u.updated_at,
			r.id, r.name, r.level
		FROM users u
		LEFT JOIN roles r ON u.role_id = r.id
		WHERE u.email = $1 AND u.deleted_at IS NULL
	`
	var user entity.User
	var userRoleID uuid.NullUUID
	var roleID uuid.NullUUID
	var roleName sql.NullString
	var roleLevel sql.NullInt32

	err := r.db.QueryRow(subCtx, query, email).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.FullName,
		&userRoleID,
		&user.Phone,
		&user.AvatarURL,
		&user.TwoFaSecret,
		&user.IsActive,
		&user.IsVerified,
		&user.CreatedAt,
		&user.UpdatedAt,
		&roleID,
		&roleName,
		&roleLevel,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to find user by email: %w", err)
	}

	if userRoleID.Valid {
		user.RoleID = userRoleID.UUID
	} else {
		user.RoleID = uuid.Nil
	}

	if roleID.Valid {
		user.Role = &entity.Role{ID: roleID.UUID, Name: roleName.String, Level: entity.RoleLevel(roleLevel.Int32)}
	}
	return &user, nil
}
func (r *userRepository) Update(ctx context.Context, user *entity.User) (*entity.User, error) {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	query := `
		UPDATE users SET
			email = $1,
			full_name = $2,
			role_id = COALESCE(NULLIF($3, '00000000-0000-0000-0000-000000000000'::uuid), role_id, (SELECT id FROM roles WHERE name = 'user' LIMIT 1)),
			phone = $4,
			avatar_url = $5,
			is_active = $6,
			is_verified = $7,
			email_verified_at = $8
		WHERE id = $9 AND deleted_at IS NULL
		RETURNING id, email, full_name, role_id, phone, avatar_url, is_active, is_verified, email_verified_at, last_login_at, created_at, updated_at
	`
	args := []interface{}{
		user.Email,
		user.FullName,
		user.RoleID,
		user.Phone,
		user.AvatarURL,
		user.IsActive,
		user.IsVerified,
		user.EmailVerifiedAt,
		user.ID,
	}

	updatedUser := &entity.User{}
	err := r.db.QueryRow(
		subCtx,
		query,
		args...).Scan(
		&updatedUser.ID,
		&updatedUser.Email,
		&updatedUser.FullName,
		&updatedUser.RoleID,
		&updatedUser.Phone,
		&updatedUser.AvatarURL,
		&updatedUser.IsActive,
		&updatedUser.IsVerified,
		&updatedUser.EmailVerifiedAt,
		&updatedUser.LastLoginAt,
		&updatedUser.CreatedAt,
		&updatedUser.UpdatedAt,
	)

	if err != nil {
		var pgxErr *pgconn.PgError
		if errors.As(err, &pgxErr) && pgxErr.Code == "23505" {
			switch pgxErr.ConstraintName {
			case "users_email_key":
				return nil, fmt.Errorf("email %s is already registered", user.Email)
			case "users_name_length_check":
				return nil, fmt.Errorf("full name %s is invalid, must be between 2 and 50 characters", user.FullName)
			default:
				return nil, fmt.Errorf("unique constraint violation (%s): %w", pgxErr.ConstraintName, err)
			}
		}
		return nil, fmt.Errorf("failed to update user: %w", err)
	}
	return updatedUser, nil
}
func (r *userRepository) UpdateTx(ctx context.Context, tx pgx.Tx, user *entity.User) (*entity.User, error) {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	query := `
		UPDATE users SET
			email = $1,
			full_name = $2,
			role_id = COALESCE(NULLIF($3, '00000000-0000-0000-0000-000000000000'::uuid), role_id, (SELECT id FROM roles WHERE name = 'user' LIMIT 1)),
			phone = $4,
			avatar_url = $5,
			is_active = $6,
			is_verified = $7,
			email_verified_at = $8
		WHERE id = $9 AND deleted_at IS NULL
		RETURNING id, email, full_name, role_id, phone, avatar_url, is_active, is_verified, email_verified_at, last_login_at, created_at, updated_at
	`
	args := []interface{}{
		user.Email,
		user.FullName,
		user.RoleID,
		user.Phone,
		user.AvatarURL,
		user.IsActive,
		user.IsVerified,
		user.EmailVerifiedAt,
		user.ID,
	}

	updatedUser := &entity.User{}
	err := tx.QueryRow(subCtx, query, args...).Scan(
		&updatedUser.ID,
		&updatedUser.Email,
		&updatedUser.FullName,
		&updatedUser.RoleID,
		&updatedUser.Phone,
		&updatedUser.AvatarURL,
		&updatedUser.IsActive,
		&updatedUser.IsVerified,
		&updatedUser.EmailVerifiedAt,
		&updatedUser.LastLoginAt,
		&updatedUser.CreatedAt,
		&updatedUser.UpdatedAt,
	)

	if err != nil {
		var pgxErr *pgconn.PgError
		if errors.As(err, &pgxErr) && pgxErr.Code == "23505" {
			switch pgxErr.ConstraintName {
			case "users_email_key":
				return nil, fmt.Errorf("email %s is already registered", user.Email)
			case "users_name_length_check":
				return nil, fmt.Errorf("full name %s is invalid, must be between 2 and 50 characters", user.FullName)
			default:
				return nil, fmt.Errorf("unique constraint violation (%s): %w", pgxErr.ConstraintName, err)
			}
		}
		return nil, fmt.Errorf("failed to update user: %w", err)
	}
	return updatedUser, nil
}
func (r *userRepository) UpdateWithRecovery(ctx context.Context, user *entity.User) (*entity.User, error) {
	var updatedUser *entity.User

	err := r.WithTransaction(ctx, func(tx pgx.Tx) error {
		var innerErr error
		updatedUser, innerErr = r.UpdateTx(ctx, tx, user) // ✅ Gunakan ctx utama
		return innerErr
	})

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("user not found or already updated")
		}
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return updatedUser, nil
}
func (r *userRepository) UpdateTwoFaSecret(ctx context.Context, id uuid.UUID, twoFaSecret *string) error {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	query := `UPDATE users SET two_fa_secret = $1 WHERE id = $2 AND deleted_at IS NULL`
	args := []interface{}{
		twoFaSecret,
		id,
	}
	result, err := r.db.Exec(subCtx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to update two fa secret: %w", err)
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("user not found or already updated")
	}
	return nil
}
func (r *userRepository) RemoveTwoFaSecret(ctx context.Context, id uuid.UUID) error {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	query := `UPDATE users SET two_fa_secret = NULL WHERE id = $1 AND deleted_at IS NULL`
	args := []interface{}{
		id,
	}
	result, err := r.db.Exec(subCtx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to remove two fa secret: %w", err)
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("user not found or already updated")
	}
	return nil
}
func (r *userRepository) SetEmailVerified(ctx context.Context, id uuid.UUID, isVerified bool) error {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	query := `UPDATE users SET is_verified = $1, email_verified_at = NOW() WHERE id = $2 AND deleted_at IS NULL`
	args := []interface{}{
		isVerified,
		id,
	}
	result, err := r.db.Exec(subCtx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to set email verified: %w", err)
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("user not found or already updated")
	}
	return nil
}
func (r *userRepository) UpdatePassword(ctx context.Context, id uuid.UUID, passwordHash string) error {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	query := `UPDATE users SET password_hash = $1 WHERE id = $2 AND deleted_at IS NULL`
	args := []interface{}{
		passwordHash,
		id,
	}
	result, err := r.db.Exec(subCtx, query, args...)
	if err != nil {
		r.logger.Error("Failed to update password", zap.Error(err))
		return fmt.Errorf("failed to update password: %w", err)
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("user not found or already updated")
	}
	return nil
}
func (r *userRepository) UpdateAvatar(ctx context.Context, id uuid.UUID, avatarURL string) (string, error) {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	query := `UPDATE users SET avatar_url = $1 WHERE id = $2 AND deleted_at IS NULL RETURNING avatar_url`
	args := []interface{}{
		avatarURL,
		id,
	}

	var newAvatarURL string
	err := r.db.QueryRow(subCtx, query, args...).Scan(&newAvatarURL)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", fmt.Errorf("user not found or already updated")
		}
		return "", fmt.Errorf("failed to update avatar: %w", err)
	}
	return newAvatarURL, nil
}
func (r *userRepository) UpdateLastLogin(ctx context.Context, id uuid.UUID) error {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	query := `
		UPDATE users SET
			last_login_at = NOW(),
			updated_at = NOW()
		WHERE id = $1 AND deleted_at IS NULL
	`
	args := []interface{}{
		id,
	}
	result, err := r.db.Exec(subCtx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to update last login: %w", err)
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("user not found or already updated")
	}
	return nil
}

func (r *userRepository) Delete(ctx context.Context, id uuid.UUID) error {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	query := `
		UPDATE users SET
			deleted_at = NOW()
		WHERE id = $1 AND deleted_at IS NULL
	`
	args := []interface{}{
		id,
	}
	result, err := r.db.Exec(subCtx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("user not found or already deleted")
	}
	return nil
}
func (r *userRepository) HardDelete(ctx context.Context, id uuid.UUID) error {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	query := `
		DELETE FROM users WHERE id = $1
	`
	args := []interface{}{
		id,
	}
	result, err := r.db.Exec(subCtx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to hard delete user: %w", err)
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("user not found or already deleted")
	}
	return nil
}
func (r *userRepository) Restore(ctx context.Context, id uuid.UUID) error {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	query := `
		UPDATE users SET
			deleted_at = NULL,
			updated_at = NOW()
		WHERE id = $1 AND deleted_at IS NOT NULL
	`
	args := []interface{}{
		id,
	}
	result, err := r.db.Exec(subCtx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to restore user: %w", err)
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("user not found or already restored")
	}
	return nil
}
func (r *userRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	query := `
		SELECT EXISTS(
			SELECT 1 FROM users WHERE email = $1 AND deleted_at IS NULL
		)
	`
	args := []interface{}{
		email,
	}
	var exists bool
	err := pgxscan.Get(subCtx, r.db, &exists, query, args...)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}
		return false, fmt.Errorf("failed to check if user exists by email: %w", err)
	}
	return exists, nil
}
func (r *userRepository) ExistsByPhone(ctx context.Context, phone string) (bool, error) {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	query := `
		SELECT EXISTS(
			SELECT 1 FROM users WHERE phone = $1 AND deleted_at IS NULL
		)
	`
	args := []interface{}{
		phone,
	}
	var exists bool
	err := pgxscan.Get(subCtx, r.db, &exists, query, args...)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}
		return false, fmt.Errorf("failed to check if user exists by phone: %w", err)
	}
	return exists, nil
}
func (r *userRepository) BulkDelete(ctx context.Context, ids []uuid.UUID) error {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	query := `DELETE FROM users WHERE id = ANY($1) AND deleted_at IS NULL
	`
	args := []interface{}{
		ids,
	}
	result, err := r.db.Exec(subCtx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to bulk delete users: %w", err)
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("no users found or already deleted")
	}
	return nil
}
func (r *userRepository) Search(ctx context.Context, params dto.QueryParamsRequest) ([]*entity.User, dto.Pagination, error) {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	qb := NewQueryBuilder(`
		SELECT 
			u.id, u.email, u.full_name, u.avatar_url, u.role_id, u.phone, u.is_active, u.is_verified, u.email_verified_at, u.last_login_at, u.created_at, u.updated_at,
			r.id, r.name, r.level
		FROM users u
		LEFT JOIN roles r ON u.role_id = r.id
	`)

	if params.Search != "" {
		qb.Where("(u.email ILIKE $? OR u.full_name ILIKE $?)",
			"%"+params.Search+"%",
			"%"+params.Search+"%",
		)
	}

	if !params.DateFrom.IsZero() && !params.DateTo.IsZero() {
		qb.Where("u.created_at BETWEEN $? AND $?", params.DateFrom, params.DateTo)
	}

	if params.SortBy != "" {
		sortBy := params.SortBy
		switch sortBy {
		case "created_at", "email", "full_name", "is_active", "last_login_at":
			sortBy = "u." + sortBy
		case "role":
			sortBy = "r.name"
		case "updated_at":
			sortBy = "u.updated_at"
		default:
			sortBy = "u.created_at"
		}
		qb.OrderByField(sortBy, params.OrderBy)
	} else {
		qb.OrderByField("u.created_at", "DESC")
	}

	countQuery := `SELECT COUNT(*) FROM users`
	args := []interface{}{}
	whereClauses := []string{}
	argIdx := 1

	if params.Search != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("(email ILIKE $%d OR full_name ILIKE $%d)", argIdx, argIdx+1))
		args = append(args, "%"+params.Search+"%", "%"+params.Search+"%")
		argIdx += 2
	}

	if !params.DateFrom.IsZero() && !params.DateTo.IsZero() {
		whereClauses = append(whereClauses, fmt.Sprintf("created_at BETWEEN $%d AND $%d", argIdx, argIdx+1))
		args = append(args, params.DateFrom, params.DateTo)
	}

	if len(whereClauses) > 0 {
		countQuery += " WHERE " + strings.Join(whereClauses, " AND ")
	}

	var totalItems int64
	err := r.db.QueryRow(subCtx, countQuery, args...).Scan(&totalItems)
	if err != nil {
		return nil, dto.Pagination{}, fmt.Errorf("failed to count users: %w", err)
	}

	offset := (params.Page - 1) * params.Limit
	qb.WithLimit(params.Limit).WithOffset(offset)

	query, finalArgs := qb.Build()
	rows, err := r.db.Query(subCtx, query, finalArgs...)
	if err != nil {
		return nil, dto.Pagination{}, fmt.Errorf("failed to query users: %w", err)
	}
	defer rows.Close()

	var users []*entity.User
	for rows.Next() {
		var user entity.User
		var role entity.Role
		err := rows.Scan(
			&user.ID,
			&user.Email,
			&user.FullName,
			&user.AvatarURL,
			&user.RoleID,
			&user.Phone,
			&user.IsActive,
			&user.IsVerified,
			&user.EmailVerifiedAt,
			&user.LastLoginAt,
			&user.CreatedAt,
			&user.UpdatedAt,
			&role.ID,
			&role.Name,
			&role.Level,
		)
		if err != nil {
			return nil, dto.Pagination{}, fmt.Errorf("failed to scan user: %w", err)
		}
		user.Role = &role
		users = append(users, &user)
	}

	if err := rows.Err(); err != nil {
		return nil, dto.Pagination{}, fmt.Errorf("rows iteration error: %w", err)
	}

	totalPages := 0
	if params.Limit > 0 {
		totalPages = int((totalItems + int64(params.Limit) - 1) / int64(params.Limit))
	}

	return users, dto.Pagination{
		CurrentPage: params.Page,
		Limit:       params.Limit,
		TotalItems:  totalItems,
		TotalPages:  totalPages,
		HasNext:     params.Page < totalPages,
		HasPrev:     params.Page > 1,
	}, nil
}

func (r *userRepository) FindByRoleID(ctx context.Context, roleID uuid.UUID) ([]*entity.User, error) {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	query := `
		SELECT 
			u.id, u.email, u.password_hash, u.full_name, u.role_id, u.phone, u.avatar_url,
			u.two_fa_secret, u.is_active, u.is_verified, u.created_at, u.updated_at,
			r.id, r.name, r.level
		FROM users u
		LEFT JOIN roles r ON u.role_id = r.id
		WHERE u.role_id = $1 AND u.deleted_at IS NULL
	`

	rows, err := r.db.Query(subCtx, query, roleID)
	if err != nil {
		return nil, fmt.Errorf("failed to query users by role_id: %w", err)
	}
	defer rows.Close()

	var users []*entity.User
	for rows.Next() {
		var user entity.User
		var userRoleID uuid.NullUUID
		var roleID uuid.NullUUID
		var roleName sql.NullString
		var roleLevel sql.NullInt32

		err := rows.Scan(
			&user.ID,
			&user.Email,
			&user.PasswordHash,
			&user.FullName,
			&userRoleID,
			&user.Phone,
			&user.AvatarURL,
			&user.TwoFaSecret,
			&user.IsActive,
			&user.IsVerified,
			&user.CreatedAt,
			&user.UpdatedAt,
			&roleID,
			&roleName,
			&roleLevel,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}

		if userRoleID.Valid {
			user.RoleID = userRoleID.UUID
		} else {
			user.RoleID = uuid.Nil
		}
		if roleID.Valid {
			user.Role = &entity.Role{ID: roleID.UUID, Name: roleName.String, Level: entity.RoleLevel(roleLevel.Int32)}
		}

		users = append(users, &user)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate rows: %w", err)
	}
	return users, nil
}

func (r *userRepository) FindRoleByName(ctx context.Context, name string) (*entity.Role, error) {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	query := `
		SELECT id, name, slug, description, level, created_at, updated_at
		FROM roles
		WHERE name = $1 AND deleted_at IS NULL
		LIMIT 1
	`
	var role entity.Role
	err := pgxscan.Get(subCtx, r.db, &role, query, name)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("role not found")
		}
		return nil, fmt.Errorf("failed to find role: %w", err)
	}
	return &role, nil
}
func (r *userRepository) SetRoleID(ctx context.Context, userID uuid.UUID, roleID uuid.UUID) error {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	query := `UPDATE users SET role_id = $1 WHERE id = $2 AND deleted_at IS NULL`
	result, err := r.db.Exec(subCtx, query, roleID, userID)
	if err != nil {
		return fmt.Errorf("failed to update user role: %w", err)
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("user not found or already updated")
	}
	return nil
}
