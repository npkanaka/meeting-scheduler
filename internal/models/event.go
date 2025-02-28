package models

import (
	"time"

	"github.com/google/uuid"
)

// EventStatus represents the status of an event
type EventStatus string

const (
	EventStatusDraft    EventStatus = "draft"
	EventStatusActive   EventStatus = "active"
	EventStatusCanceled EventStatus = "canceled"
)

// Event represents a meeting or event
type Event struct {
	ID          uuid.UUID   `json:"id" gorm:"type:uuid;primary_key"`
	Title       string      `json:"title" gorm:"not null"`
	Description string      `json:"description"`
	CreatorID   uuid.UUID   `json:"creator_id" gorm:"type:uuid;not null"`
	Duration    int         `json:"duration" gorm:"not null"` // Duration in minutes
	Status      EventStatus `json:"status" gorm:"not null"`
	CreatedAt   time.Time   `json:"created_at" gorm:"not null"`
	UpdatedAt   time.Time   `json:"updated_at" gorm:"not null"`
}
