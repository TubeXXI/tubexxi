package dto

import (
	"time"

	"github.com/google/uuid"
)

type AnalyticsDaily struct {
	ID         uuid.UUID `json:"id" db:"id"`
	Date       time.Time `json:"date" db:"date"`
	TotalUsers int       `json:"total_users" db:"total_users"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
	TotalViews int       `json:"total_views" db:"total_views"`
}
