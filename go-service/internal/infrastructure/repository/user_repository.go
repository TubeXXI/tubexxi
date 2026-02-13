package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
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
	Search(ctx context.Context, opts *ListOptions) ([]*entity.User, int64, error)
	Count(ctx context.Context, filter *Filter) (int64, error)
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

	query := `
		UPDATE users SET
			deleted_at = NOW()
		WHERE id = ANY($1) AND deleted_at IS NULL
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

func (r *userRepository) Search(ctx context.Context, opts *ListOptions) ([]*entity.User, int64, error) {
	subCtx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	if opts == nil {
		opts = NewListOptions()
	}

	totalRows, err := r.Count(ctx, opts.Filter)
	if err != nil {
		return nil, 0, err
	}

	qb := r.buildBaseQuery("SELECT * FROM users", opts.Filter)

	if opts.OrderBy != "" {
		qb.OrderByField(opts.OrderBy, opts.OrderDir)
	} else {
		qb.OrderByField("created_at", "DESC")
	}
	if opts.Pagination != nil && opts.Pagination.Limit > 0 {
		qb.WithLimit(opts.Pagination.Limit)
		if opts.Pagination.Page > 1 {
			qb.WithOffset(opts.Pagination.GetOffset())
		}
	}

	query, args := qb.Build()
	var users []*entity.User
	err = pgxscan.Select(subCtx, r.db, &users, query, args...)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, 0, fmt.Errorf("no users found")
		}
		return nil, 0, fmt.Errorf("failed to list users: %w", err)
	}

	return users, totalRows, nil
}
func (r *userRepository) Count(ctx context.Context, filter *Filter) (int64, error) {
	subCtx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	qb := r.buildBaseQuery("SELECT COUNT(*) FROM users", filter)
	query, args := qb.Build()

	var count int64
	err := r.db.QueryRow(subCtx, query, args...).Scan(&count)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, fmt.Errorf("no users found")
		}
		return 0, fmt.Errorf("failed to count users: %w", err)
	}
	return count, nil
}
func (r *userRepository) buildBaseQuery(baseQuery string, filter *Filter) *QueryBuilder {
	qb := NewQueryBuilder(baseQuery)

	if filter == nil {
		qb.Where("deleted_at IS NULL")
		return qb
	}

	if filter.IncludeDeleted != nil && *filter.IncludeDeleted {
		qb.Where("deleted_at IS NOT NULL")
	} else {
		qb.Where("deleted_at IS NULL")
	}

	if filter.Search != "" {
		searchPattern := "%" + filter.Search + "%"
		qb.Where("(email ILIKE $? OR full_name ILIKE $?)",
			searchPattern, searchPattern, searchPattern)
	}
	if filter.IsActive != nil {
		qb.Where("is_active = $?", *filter.IsActive)
	}
	if filter.UserID != nil {
		qb.Where("id = $?", *filter.UserID)
	}
	if filter.IsVerified != nil {
		qb.Where("is_verified = $?", *filter.IsVerified)
	}
	if filter.RangeDate != nil {
		var startDate time.Time
		var endDate time.Time

		if !filter.RangeDate.StartDate.IsZero() {
			startDate = filter.RangeDate.StartDate
		} else {
			startDate = time.Now().AddDate(0, 0, -7)
		}
		if !filter.RangeDate.EndDate.IsZero() {
			endDate = filter.RangeDate.EndDate
		} else {
			endDate = time.Now()
		}
		if !startDate.IsZero() || !endDate.IsZero() {
			qb.Where("created_at BETWEEN $? AND $?", startDate, endDate)
		}
	}

	return qb
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
