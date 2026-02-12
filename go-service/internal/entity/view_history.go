package entity

import (
	"time"

	"github.com/google/uuid"
)

const (
	ViewHistoryTypeMovie  = "movies"
	ViewHistoryTypeSeries = "series"
	ViewHistoryTypeAnime  = "anime"
)

type ViewHistoryType string

type ViewHistory struct {
	ID              uuid.UUID       `json:"id" db:"id"`
	UserID          uuid.UUID       `json:"user_id" db:"user_id"`
	Name            *string         `json:"name,omitempty" db:"name"`
	PageURL         *string         `json:"page_url,omitempty" db:"page_url"`
	IPAddress       *string         `json:"ip_address,omitempty" db:"ip_address"`
	UserAgent       *string         `json:"user_agent,omitempty" db:"user_agent"`
	BrowserLanguage *string         `json:"browser_language,omitempty" db:"browser_language"`
	DeviceType      *string         `json:"device_type,omitempty" db:"device_type"`
	Platform        string          `json:"platform" db:"platform"`
	ViewTime        time.Time       `json:"view_time" db:"view_time"`
	Type            ViewHistoryType `json:"type" db:"type"`
	CreatedAt       time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time       `json:"updated_at" db:"updated_at"`
}

func (vht *ViewHistory) IsValid() bool {
	return vht.Type == ViewHistoryTypeMovie ||
		vht.Type == ViewHistoryTypeSeries ||
		vht.Type == ViewHistoryTypeAnime
}
func (vht *ViewHistory) IsMovie() bool {
	return vht.Type == ViewHistoryTypeMovie
}
func (vht *ViewHistory) IsSeries() bool {
	return vht.Type == ViewHistoryTypeSeries
}
func (vht *ViewHistory) IsAnime() bool {
	return vht.Type == ViewHistoryTypeAnime
}
