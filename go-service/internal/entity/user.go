package entity

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type User struct {
	ID              uuid.UUID  `json:"id" db:"id"`
	Email           string     `json:"email" db:"email" validate:"required,email,max=255"`
	PasswordHash    string     `json:"-" db:"password_hash"`
	FullName        string     `json:"full_name" db:"full_name" validate:"required,max=255"`
	RoleID          uuid.UUID  `json:"role_id" db:"role_id"`
	Phone           NullString `json:"phone,omitempty" db:"phone" validate:"omitempty,max=20"`
	AvatarURL       NullString `json:"avatar_url,omitempty" db:"avatar_url"`
	TwoFaSecret     NullString `json:"two_fa_secret,omitempty" db:"two_fa_secret"`
	IsActive        bool       `json:"is_active" db:"is_active"`
	IsVerified      bool       `json:"is_verified" db:"is_verified"`
	EmailVerifiedAt NullTime   `json:"email_verified_at,omitempty" db:"email_verified_at"`
	LastLoginAt     NullTime   `json:"last_login_at,omitempty" db:"last_login_at"`
	CreatedAt       time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt       NullTime   `json:"deleted_at,omitempty" db:"deleted_at"`
	Role            *Role      `json:"role,omitempty" db:"-"`
}

type UserWithRole struct {
	User
	Role *Role `json:"role,omitempty"`
}

func (User) TableName() string {
	return "users"
}

func (u *User) ToResponse() *User {
	return &User{
		ID:              u.ID,
		Email:           u.Email,
		FullName:        u.FullName,
		Phone:           u.Phone,
		AvatarURL:       u.AvatarURL,
		IsActive:        u.IsActive,
		IsVerified:      u.IsVerified,
		EmailVerifiedAt: u.EmailVerifiedAt,
		LastLoginAt:     u.LastLoginAt,
		CreatedAt:       u.CreatedAt,
		UpdatedAt:       u.UpdatedAt,
	}
}

type JWTClaims struct {
	UserID    string `json:"uid"`
	Email     string `json:"em,omitempty"`
	RoleID    string `json:"rid"`
	RoleLevel int    `json:"role_level"`
	SessionID string `json:"sid"`
	jwt.RegisteredClaims
}

func (u *User) IsDeleted() bool {
	return u.DeletedAt.Valid
}

func (u *User) CanLogin() bool {
	return u.IsActive && !u.IsDeleted()
}

func (u *User) MarkAsVerified() {
	now := time.Now()
	u.IsVerified = true
	u.EmailVerifiedAt = NewNullTimePtr(&now)
}

func (u *User) UpdateLastLogin() {
	now := time.Now()
	u.LastLoginAt = NewNullTimePtr(&now)
}
func (u *User) IsTwoFaActive() bool {
	return u.TwoFaSecret.Valid
}
