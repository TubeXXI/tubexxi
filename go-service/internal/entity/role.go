package entity

import (
	"time"

	"github.com/google/uuid"
)

const (
	RoleLevelUser       = 0
	RoleLevelSuperAdmin = 1
	RoleLevelAdmin      = 2
)
const (
	RoleSuperAdmin = "superadmin"
	RoleAdmin      = "admin"
	RoleUser       = "user"
)

type Role struct {
	ID          uuid.UUID  `json:"id" db:"id"`
	Name        string     `json:"name" db:"name" validate:"required,max=50"`
	Slug        string     `json:"slug" db:"slug" validate:"required,max=50"`
	Description NullString `json:"description,omitempty" db:"description"`
	Level       int        `json:"level" db:"level" validate:"required,min=1,max=2"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt   NullTime   `json:"deleted_at,omitempty" db:"deleted_at"`
}

// TableName returns the table name for Role
func (Role) TableName() string {
	return "roles"
}

// IsSuperAdmin checks if role is superadmin
func (r *Role) IsSuperAdmin() bool {
	return r.Level == RoleLevelSuperAdmin
}

// IsAdmin checks if role is admin
func (r *Role) IsAdmin() bool {
	return r.Level == RoleLevelAdmin
}

// IsUser checks if role is user
func (r *Role) IsUser() bool {
	return r.Level == RoleLevelUser
}
