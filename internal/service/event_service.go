// internal/service/event_service.go
package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/npkanaka/meeting-scheduler/internal/models"
	"github.com/npkanaka/meeting-scheduler/internal/repository"
)

// EventService handles event business logic
type EventService struct {
	eventRepo repository.EventRepository
}

// NewEventService creates a new EventService
func NewEventService(eventRepo repository.EventRepository) *EventService {
	return &EventService{
		eventRepo: eventRepo,
	}
}

// CreateEvent creates a new event
func (s *EventService) CreateEvent(ctx context.Context, req *models.CreateEventRequest, creatorID uuid.UUID) (*models.Event, error) {
	now := time.Now()
	event := &models.Event{
		Title:       req.Title,
		Description: req.Description,
		CreatorID:   creatorID,
		Duration:    req.Duration,
		Status:      models.EventStatusDraft,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := s.eventRepo.Create(ctx, event); err != nil {
		return nil, err
	}

	return event, nil
}

// GetEvent retrieves an event by ID
func (s *EventService) GetEvent(ctx context.Context, id uuid.UUID) (*models.Event, error) {
	return s.eventRepo.GetByID(ctx, id)
}

// UpdateEvent updates an existing event
func (s *EventService) UpdateEvent(ctx context.Context, id uuid.UUID, req *models.CreateEventRequest) (*models.Event, error) {
	event, err := s.eventRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	event.Title = req.Title
	event.Description = req.Description
	event.Duration = req.Duration
	event.UpdatedAt = time.Now()

	if err := s.eventRepo.Update(ctx, event); err != nil {
		return nil, err
	}

	return event, nil
}

// DeleteEvent removes an event
func (s *EventService) DeleteEvent(ctx context.Context, id uuid.UUID) error {
	return s.eventRepo.Delete(ctx, id)
}

// ListEvents returns a paginated list of events
func (s *EventService) ListEvents(ctx context.Context, limit, offset int) ([]*models.Event, error) {
	return s.eventRepo.List(ctx, limit, offset)
}
