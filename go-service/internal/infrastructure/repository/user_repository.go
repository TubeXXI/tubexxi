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
			id, email, password_hash, full_name, phone, avatar_url,
			is_active, is_verified, created_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id, created_at, updated_at
	`
	err := r.db.QueryRow(subCtx,
		query,
		user.ID,
		user.Email,
		user.PasswordHash,
		user.FullName,
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
			id, email, password_hash, full_name, phone, avatar_url,
			is_active, is_verified, created_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id, created_at, updated_at
	`
	err := tx.QueryRow(subCtx,
		query,
		user.ID,
		user.Email,
		user.PasswordHash,
		user.FullName,
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
		SELECT * FROM users WHERE id = $1 AND deleted_at IS NULL
	`
	var user entity.User
	err := pgxscan.Get(subCtx, r.db, &user, query, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to find user by id: %w", err)
	}
	return &user, nil
}
func (r *userRepository) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	query := `
		SELECT * FROM users WHERE email = $1 AND deleted_at IS NULL
	`
	var user entity.User
	err := pgxscan.Get(subCtx, r.db, &user, query, email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to find user by email: %w", err)
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
			phone = $3,
			avatar_url = $4,
			is_active = $5,
			is_verified = $6,
			email_verified_at = $7
		WHERE id = $8 AND deleted_at IS NULL
		RETURNING id, email, full_name, phone, avatar_url, is_active, is_verified, email_verified_at, last_login_at, created_at, updated_at
	`
	args := []interface{}{
		user.Email,
		user.FullName,
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
			phone = $3,
			avatar_url = $4,
			is_active = $5,
			is_verified = $6,
			email_verified_at = $7
		WHERE id = $8 AND deleted_at IS NULL
		RETURNING id, email, full_name, phone, avatar_url, is_active, is_verified, email_verified_at, last_login_at, created_at, updated_at
	`
	args := []interface{}{
		user.Email,
		user.FullName,
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
