package entity

import (
	"time"

	"github.com/google/uuid"
)

const (
	TicketStatusOpen     = "open"
	TicketStatusResolved = "resolved"
	TicketStatusClosed   = "closed"
)

type TicketStatus string

type Ticket struct {
	ID          uuid.UUID    `json:"id" gorm:"primaryKey"`
	UserID      uuid.UUID    `json:"user_id" gorm:"not null"`
	Title       string       `json:"title" gorm:"not null"`
	Description string       `json:"description" gorm:"not null"`
	Status      TicketStatus `json:"status" gorm:"not null;default:'open'"`
	CreatedAt   time.Time    `json:"created_at" gorm:"not null;default:now()"`
	UpdatedAt   time.Time    `json:"updated_at" gorm:"not null;default:now()"`
}

func (t *Ticket) TableName() string {
	return "tickets"
}
func (t *TicketStatus) String() string {
	return string(*t)
}
func (t *TicketStatus) IsOpen() bool {
	return *t == TicketStatusOpen
}
func (t *TicketStatus) IsResolved() bool {
	return *t == TicketStatusResolved
}
func (t *TicketStatus) IsClosed() bool {
	return *t == TicketStatusClosed
}
