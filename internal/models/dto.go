package models

import (
	"time"

	"github.com/google/uuid"
)

// CreateEventRequest represents the request to create a new event
type CreateEventRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	Duration    int    `json:"duration" binding:"required,min=1"`
}

// CreateEventResponse represents the response after creating an event
type CreateEventResponse struct {
	ID          uuid.UUID   `json:"id"`
	Title       string      `json:"title"`
	Description string      `json:"description"`
	CreatorID   uuid.UUID   `json:"creator_id"`
	Duration    int         `json:"duration"`
	Status      EventStatus `json:"status"`
	CreatedAt   time.Time   `json:"created_at"`
}

// TimeSlotRequest represents a request to create or update a time slot
type TimeSlotRequest struct {
	StartTime string `json:"start_time" binding:"required"` // ISO 8601 format
	EndTime   string `json:"end_time" binding:"required"`   // ISO 8601 format
}

// AvailabilityRequest represents a request to create or update availability
type AvailabilityRequest struct {
	UserID    uuid.UUID `json:"user_id" binding:"required"`
	StartTime string    `json:"start_time" binding:"required"` // ISO 8601 format
	EndTime   string    `json:"end_time" binding:"required"`   // ISO 8601 format
}

// RecommendationResponse represents the recommendation API response
type RecommendationResponse struct {
	Recommendations []Recommendation `json:"recommendations"`
}

// Recommendation represents a single time slot recommendation
type Recommendation struct {
	TimeSlot     TimeSlotResponse `json:"time_slot"`
	Attendees    []UserResponse   `json:"attendees"`
	NonAttendees []UserResponse   `json:"non_attendees"`
	Score        int              `json:"score"` // Number of attendees
}

// TimeSlotResponse represents a time slot in API responses
type TimeSlotResponse struct {
	ID        uuid.UUID `json:"id"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
}

// UserResponse represents a user in API responses
type UserResponse struct {
	ID    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	Email string    `json:"email"`
}
