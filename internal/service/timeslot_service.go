// internal/service/timeslot_service.go
package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/npkanaka/meeting-scheduler/internal/errors"
	"github.com/npkanaka/meeting-scheduler/internal/models"
	"github.com/npkanaka/meeting-scheduler/internal/repository"
)

// TimeSlotService handles time slot business logic
type TimeSlotService struct {
	timeslotRepo repository.TimeSlotRepository
	eventRepo    repository.EventRepository
}

// NewTimeSlotService creates a new TimeSlotService
func NewTimeSlotService(
	timeslotRepo repository.TimeSlotRepository,
	eventRepo repository.EventRepository,
) *TimeSlotService {
	return &TimeSlotService{
		timeslotRepo: timeslotRepo,
		eventRepo:    eventRepo,
	}
}

// CreateTimeSlot creates a new time slot for an event
func (s *TimeSlotService) CreateTimeSlot(ctx context.Context, eventID uuid.UUID, req *models.TimeSlotRequest) (*models.TimeSlot, error) {
	// Verify the event exists
	_, err := s.eventRepo.GetByID(ctx, eventID)
	if err != nil {
		return nil, err
	}

	// Parse time strings
	startTime, err := time.Parse(time.RFC3339, req.StartTime)
	if err != nil {
		return nil, err
	}

	endTime, err := time.Parse(time.RFC3339, req.EndTime)
	if err != nil {
		return nil, err
	}

	// Validate time range
	if endTime.Before(startTime) || endTime.Equal(startTime) {
		return nil, errors.ErrInvalidTimeRange
	}

	now := time.Now()
	slot := &models.TimeSlot{
		EventID:   eventID,
		StartTime: startTime,
		EndTime:   endTime,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := s.timeslotRepo.Create(ctx, slot); err != nil {
		return nil, err
	}

	return slot, nil
}

// GetTimeSlot retrieves a time slot by ID
func (s *TimeSlotService) GetTimeSlot(ctx context.Context, id uuid.UUID) (*models.TimeSlot, error) {
	return s.timeslotRepo.GetByID(ctx, id)
}

// UpdateTimeSlot updates an existing time slot
func (s *TimeSlotService) UpdateTimeSlot(ctx context.Context, id uuid.UUID, req *models.TimeSlotRequest) (*models.TimeSlot, error) {
	slot, err := s.timeslotRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Parse time strings
	startTime, err := time.Parse(time.RFC3339, req.StartTime)
	if err != nil {
		return nil, err
	}

	endTime, err := time.Parse(time.RFC3339, req.EndTime)
	if err != nil {
		return nil, err
	}

	// Validate time range
	if endTime.Before(startTime) || endTime.Equal(startTime) {
		return nil, errors.ErrInvalidTimeRange
	}

	slot.StartTime = startTime
	slot.EndTime = endTime
	slot.UpdatedAt = time.Now()

	if err := s.timeslotRepo.Update(ctx, slot); err != nil {
		return nil, err
	}

	return slot, nil
}

// DeleteTimeSlot removes a time slot
func (s *TimeSlotService) DeleteTimeSlot(ctx context.Context, id uuid.UUID) error {
	return s.timeslotRepo.Delete(ctx, id)
}

// GetEventTimeSlots retrieves all time slots for an event
func (s *TimeSlotService) GetEventTimeSlots(ctx context.Context, eventID uuid.UUID) ([]*models.TimeSlot, error) {
	return s.timeslotRepo.GetByEventID(ctx, eventID)
}
