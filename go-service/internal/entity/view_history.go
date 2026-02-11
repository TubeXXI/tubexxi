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
	ID        uuid.UUID       `json:"id" db:"id"`
	UserID    uuid.UUID       `json:"user_id" db:"user_id"`
	Name      *string         `json:"name,omitempty" db:"name"`
	PageURL   *string         `json:"page_url,omitempty" db:"page_url"`
	ViewTime  time.Time       `json:"view_time" db:"view_time"`
	Type      ViewHistoryType `json:"type" db:"type"`
	CreatedAt time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt time.Time       `json:"updated_at" db:"updated_at"`
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
