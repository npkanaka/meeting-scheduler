package models

import (
	"time"

	"github.com/google/uuid"
)

// TimeSlot represents a potential time slot for an event
type TimeSlot struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primary_key"`
	EventID   uuid.UUID `json:"event_id" gorm:"type:uuid;not null"`
	StartTime time.Time `json:"start_time" gorm:"not null"`
	EndTime   time.Time `json:"end_time" gorm:"not null"`
	CreatedAt time.Time `json:"created_at" gorm:"not null"`
	UpdatedAt time.Time `json:"updated_at" gorm:"not null"`
}
