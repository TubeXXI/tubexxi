package dto

import "time"

type QueryParamsRequest struct {
	Search         string                 `json:"search,omitempty" validate:"omitempty"`
	SortBy         string                 `json:"sort_by,omitempty" validate:"omitempty"`
	OrderBy        string                 `json:"order_by,omitempty" validate:"omitempty"`
	Page           int                    `json:"page,omitempty" validate:"required,min=1"`
	Limit          int                    `json:"limit,omitempty" validate:"required,min=1,max=100"`
	Status         string                 `json:"status,omitempty" validate:"omitempty"`
	Type           string                 `json:"type,omitempty" validate:"omitempty"`
	IncludeDeleted bool                   `json:"include_deleted,omitempty" validate:"omitempty"`
	IsActive       bool                   `json:"is_active,omitempty" validate:"omitempty"`
	UserID         string                 `json:"user_id,omitempty" validate:"omitempty"`
	DateFrom       time.Time              `json:"date_from,omitempty" validate:"omitempty"`
	DateTo         time.Time              `json:"date_to,omitempty" validate:"omitempty"`
	Extra          map[string]interface{} `json:"extra,omitempty" validate:"dive,key,required,value"`
}

type Pagination struct {
	CurrentPage int   `json:"current_page"`
	Limit       int   `json:"limit"`
	TotalItems  int64 `json:"total_items"`
	TotalPages  int   `json:"total_pages"`
	HasNext     bool  `json:"has_next"`
	HasPrev     bool  `json:"has_prev"`
}
