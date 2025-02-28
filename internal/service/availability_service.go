// internal/service/availability_service.go
package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/npkanaka/meeting-scheduler/internal/errors"
	"github.com/npkanaka/meeting-scheduler/internal/models"
	"github.com/npkanaka/meeting-scheduler/internal/repository"
)

// AvailabilityService handles availability business logic
type AvailabilityService struct {
	availabilityRepo repository.AvailabilityRepository
	eventRepo        repository.EventRepository
	userRepo         repository.UserRepository
}

// NewAvailabilityService creates a new AvailabilityService
func NewAvailabilityService(
	availabilityRepo repository.AvailabilityRepository,
	eventRepo repository.EventRepository,
	userRepo repository.UserRepository,
) *AvailabilityService {
	return &AvailabilityService{
		availabilityRepo: availabilityRepo,
		eventRepo:        eventRepo,
		userRepo:         userRepo,
	}
}

// CreateAvailability creates a new availability record
func (s *AvailabilityService) CreateAvailability(ctx context.Context, eventID uuid.UUID, req *models.AvailabilityRequest) (*models.Availability, error) {
	// Verify the event exists
	_, err := s.eventRepo.GetByID(ctx, eventID)
	if err != nil {
		return nil, err
	}

	// Verify the user exists
	_, err = s.userRepo.GetByID(ctx, req.UserID)
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
	availability := &models.Availability{
		UserID:    req.UserID,
		EventID:   eventID,
		StartTime: startTime,
		EndTime:   endTime,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := s.availabilityRepo.Create(ctx, availability); err != nil {
		return nil, err
	}

	return availability, nil
}

// UpdateAvailability updates an existing availability record
func (s *AvailabilityService) UpdateAvailability(ctx context.Context, id uuid.UUID, req *models.AvailabilityRequest) (*models.Availability, error) {
	// Get existing availabilities for this user and event
	availabilities, err := s.availabilityRepo.GetByUserAndEvent(ctx, req.UserID, id)
	if err != nil || len(availabilities) == 0 {
		return nil, errors.ErrAvailabilityNotFound
	}

	// For simplicity, we'll update the first availability entry
	// In a real app, you might want to handle this differently
	availability := availabilities[0]

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

	availability.StartTime = startTime
	availability.EndTime = endTime
	availability.UpdatedAt = time.Now()

	if err := s.availabilityRepo.Update(ctx, availability); err != nil {
		return nil, err
	}

	return availability, nil
}

// DeleteAvailability removes an availability record
func (s *AvailabilityService) DeleteAvailability(ctx context.Context, id uuid.UUID) error {
	return s.availabilityRepo.Delete(ctx, id)
}

// GetUserEventAvailability retrieves all availability records for a user and event
func (s *AvailabilityService) GetUserEventAvailability(ctx context.Context, userID, eventID uuid.UUID) ([]*models.Availability, error) {
	return s.availabilityRepo.GetByUserAndEvent(ctx, userID, eventID)
}

// GetEventAvailability retrieves all availability records for an event
func (s *AvailabilityService) GetEventAvailability(ctx context.Context, eventID uuid.UUID) ([]*models.Availability, error) {
	return s.availabilityRepo.GetByEventID(ctx, eventID)
}
