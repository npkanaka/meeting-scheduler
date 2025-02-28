package service_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/npkanaka/meeting-scheduler/internal/errors"
	"github.com/npkanaka/meeting-scheduler/internal/models"
	"github.com/npkanaka/meeting-scheduler/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateAvailability(t *testing.T) {
	// Setup mocks
	mockAvailabilityRepo := new(MockAvailabilityRepository)
	mockEventRepo := new(MockEventRepository)
	mockUserRepo := new(MockUserRepository)
	availabilityService := service.NewAvailabilityService(
		mockAvailabilityRepo,
		mockEventRepo,
		mockUserRepo,
	)

	// Prepare test data
	eventID := uuid.New()
	userID := uuid.New()
	startTime := time.Date(2025, 1, 15, 10, 0, 0, 0, time.UTC)
	endTime := time.Date(2025, 1, 15, 12, 0, 0, 0, time.UTC)

	// Mock event and user exist
	mockEventRepo.On("GetByID", mock.Anything, eventID).Return(
		&models.Event{ID: eventID},
		nil,
	)
	mockUserRepo.On("GetByID", mock.Anything, userID).Return(
		&models.User{ID: userID},
		nil,
	)

	// Prepare request
	req := &models.AvailabilityRequest{
		UserID:    userID,
		StartTime: startTime.Format(time.RFC3339),
		EndTime:   endTime.Format(time.RFC3339),
	}

	// Set expectations for availability creation
	mockAvailabilityRepo.On("Create", mock.Anything, mock.MatchedBy(func(availability *models.Availability) bool {
		return availability.UserID == userID &&
			availability.EventID == eventID &&
			availability.StartTime.Equal(startTime) &&
			availability.EndTime.Equal(endTime)
	})).Return(nil)

	// Execute the method
	availability, err := availabilityService.CreateAvailability(context.Background(), eventID, req)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, availability)
	assert.Equal(t, eventID, availability.EventID)
	assert.Equal(t, userID, availability.UserID)
	assert.Equal(t, startTime, availability.StartTime)
	assert.Equal(t, endTime, availability.EndTime)
	assert.NotZero(t, availability.CreatedAt)
	assert.NotZero(t, availability.UpdatedAt)

	// Verify mock expectations
	mockEventRepo.AssertExpectations(t)
	mockUserRepo.AssertExpectations(t)
	mockAvailabilityRepo.AssertExpectations(t)
}

func TestCreateAvailabilityEventNotFound(t *testing.T) {
	// Setup mocks
	mockAvailabilityRepo := new(MockAvailabilityRepository)
	mockEventRepo := new(MockEventRepository)
	mockUserRepo := new(MockUserRepository)
	availabilityService := service.NewAvailabilityService(
		mockAvailabilityRepo,
		mockEventRepo,
		mockUserRepo,
	)

	// Prepare test data
	eventID := uuid.New()
	userID := uuid.New()
	startTime := time.Date(2025, 1, 15, 10, 0, 0, 0, time.UTC)
	endTime := time.Date(2025, 1, 15, 12, 0, 0, 0, time.UTC)

	// Mock event not found
	mockEventRepo.On("GetByID", mock.Anything, eventID).Return(nil, assert.AnError)

	// Prepare request
	req := &models.AvailabilityRequest{
		UserID:    userID,
		StartTime: startTime.Format(time.RFC3339),
		EndTime:   endTime.Format(time.RFC3339),
	}

	// Execute the method
	availability, err := availabilityService.CreateAvailability(context.Background(), eventID, req)

	// Assertions
	assert.Error(t, err)
	assert.Nil(t, availability)

	// Verify mock expectations
	mockEventRepo.AssertExpectations(t)
}

func TestCreateAvailabilityUserNotFound(t *testing.T) {
	// Setup mocks
	mockAvailabilityRepo := new(MockAvailabilityRepository)
	mockEventRepo := new(MockEventRepository)
	mockUserRepo := new(MockUserRepository)
	availabilityService := service.NewAvailabilityService(
		mockAvailabilityRepo,
		mockEventRepo,
		mockUserRepo,
	)

	// Prepare test data
	eventID := uuid.New()
	userID := uuid.New()
	startTime := time.Date(2025, 1, 15, 10, 0, 0, 0, time.UTC)
	endTime := time.Date(2025, 1, 15, 12, 0, 0, 0, time.UTC)

	// Mock event exists
	mockEventRepo.On("GetByID", mock.Anything, eventID).Return(
		&models.Event{ID: eventID},
		nil,
	)

	// Mock user not found
	mockUserRepo.On("GetByID", mock.Anything, userID).Return(nil, assert.AnError)

	// Prepare request
	req := &models.AvailabilityRequest{
		UserID:    userID,
		StartTime: startTime.Format(time.RFC3339),
		EndTime:   endTime.Format(time.RFC3339),
	}

	// Execute the method
	availability, err := availabilityService.CreateAvailability(context.Background(), eventID, req)

	// Assertions
	assert.Error(t, err)
	assert.Nil(t, availability)

	// Verify mock expectations
	mockEventRepo.AssertExpectations(t)
	mockUserRepo.AssertExpectations(t)
}

func TestCreateAvailabilityInvalidTimeRange(t *testing.T) {
	// Setup mocks
	mockAvailabilityRepo := new(MockAvailabilityRepository)
	mockEventRepo := new(MockEventRepository)
	mockUserRepo := new(MockUserRepository)
	availabilityService := service.NewAvailabilityService(
		mockAvailabilityRepo,
		mockEventRepo,
		mockUserRepo,
	)

	// Prepare test data
	eventID := uuid.New()
	userID := uuid.New()
	startTime := time.Date(2025, 1, 15, 12, 0, 0, 0, time.UTC)
	endTime := time.Date(2025, 1, 15, 10, 0, 0, 0, time.UTC)

	// Mock event and user exist
	mockEventRepo.On("GetByID", mock.Anything, eventID).Return(
		&models.Event{ID: eventID},
		nil,
	)
	mockUserRepo.On("GetByID", mock.Anything, userID).Return(
		&models.User{ID: userID},
		nil,
	)

	// Prepare request with invalid time range
	req := &models.AvailabilityRequest{
		UserID:    userID,
		StartTime: startTime.Format(time.RFC3339),
		EndTime:   endTime.Format(time.RFC3339),
	}

	// Execute the method
	availability, err := availabilityService.CreateAvailability(context.Background(), eventID, req)

	// Assertions
	assert.Error(t, err)
	assert.Nil(t, availability)
	assert.Equal(t, errors.ErrInvalidTimeRange, err)

	// Verify mock expectations
	mockEventRepo.AssertExpectations(t)
	mockUserRepo.AssertExpectations(t)
}

func TestUpdateAvailability(t *testing.T) {
	// Setup mocks
	mockAvailabilityRepo := new(MockAvailabilityRepository)
	mockEventRepo := new(MockEventRepository)
	mockUserRepo := new(MockUserRepository)
	availabilityService := service.NewAvailabilityService(
		mockAvailabilityRepo,
		mockEventRepo,
		mockUserRepo,
	)

	// Prepare test data
	eventID := uuid.New()
	userID := uuid.New()

	// Existing availabilities
	existingAvailabilities := []*models.Availability{
		{
			ID:        uuid.New(),
			UserID:    userID,
			EventID:   eventID,
			StartTime: time.Date(2025, 1, 15, 10, 0, 0, 0, time.UTC),
			EndTime:   time.Date(2025, 1, 15, 12, 0, 0, 0, time.UTC),
		},
	}

	// New time range
	newStartTime := time.Date(2025, 1, 15, 14, 0, 0, 0, time.UTC)
	newEndTime := time.Date(2025, 1, 15, 16, 0, 0, 0, time.UTC)

	// Prepare request
	req := &models.AvailabilityRequest{
		UserID:    userID,
		StartTime: newStartTime.Format(time.RFC3339),
		EndTime:   newEndTime.Format(time.RFC3339),
	}

	// Set expectations
	mockAvailabilityRepo.On("GetByUserAndEvent", mock.Anything, userID, eventID).Return(existingAvailabilities, nil)
	mockAvailabilityRepo.On("Update", mock.Anything, mock.MatchedBy(func(availability *models.Availability) bool {
		return availability.StartTime.Equal(newStartTime) &&
			availability.EndTime.Equal(newEndTime)
	})).Return(nil)

	// Execute the method
	updatedAvailability, err := availabilityService.UpdateAvailability(context.Background(), eventID, req)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, updatedAvailability)
	assert.Equal(t, newStartTime, updatedAvailability.StartTime)
	assert.Equal(t, newEndTime, updatedAvailability.EndTime)
	assert.WithinDuration(t, time.Now(), updatedAvailability.UpdatedAt, time.Second)

	// Verify mock expectations
	mockAvailabilityRepo.AssertExpectations(t)
}

func TestUpdateAvailabilityNoExistingAvailability(t *testing.T) {
	// Setup mocks
	mockAvailabilityRepo := new(MockAvailabilityRepository)
	mockEventRepo := new(MockEventRepository)
	mockUserRepo := new(MockUserRepository)
	availabilityService := service.NewAvailabilityService(
		mockAvailabilityRepo,
		mockEventRepo,
		mockUserRepo,
	)

	// Prepare test data
	eventID := uuid.New()
	userID := uuid.New()

	// No existing availabilities
	mockAvailabilityRepo.On("GetByUserAndEvent", mock.Anything, userID, eventID).Return([]*models.Availability{}, nil)

	// Prepare request
	req := &models.AvailabilityRequest{
		UserID:    userID,
		StartTime: time.Date(2025, 1, 15, 14, 0, 0, 0, time.UTC).Format(time.RFC3339),
		EndTime:   time.Date(2025, 1, 15, 16, 0, 0, 0, time.UTC).Format(time.RFC3339),
	}

	// Execute the method
	updatedAvailability, err := availabilityService.UpdateAvailability(context.Background(), eventID, req)

	// Assertions
	assert.Error(t, err)
	assert.Nil(t, updatedAvailability)
	assert.Equal(t, errors.ErrAvailabilityNotFound, err)

	// Verify mock expectations
	mockAvailabilityRepo.AssertExpectations(t)
}

func TestUpdateAvailabilityInvalidTimeRange(t *testing.T) {
	// Setup mocks
	mockAvailabilityRepo := new(MockAvailabilityRepository)
	mockEventRepo := new(MockEventRepository)
	mockUserRepo := new(MockUserRepository)
	availabilityService := service.NewAvailabilityService(
		mockAvailabilityRepo,
		mockEventRepo,
		mockUserRepo,
	)

	// Prepare test data
	eventID := uuid.New()
	userID := uuid.New()

	// Existing availabilities
	existingAvailabilities := []*models.Availability{
		{
			ID:        uuid.New(),
			UserID:    userID,
			EventID:   eventID,
			StartTime: time.Date(2025, 1, 15, 10, 0, 0, 0, time.UTC),
			EndTime:   time.Date(2025, 1, 15, 12, 0, 0, 0, time.UTC),
		},
	}

	// Invalid time range
	newStartTime := time.Date(2025, 1, 15, 14, 0, 0, 0, time.UTC)
	newEndTime := time.Date(2025, 1, 15, 13, 0, 0, 0, time.UTC)

	// Prepare request
	req := &models.AvailabilityRequest{
		UserID:    userID,
		StartTime: newStartTime.Format(time.RFC3339),
		EndTime:   newEndTime.Format(time.RFC3339),
	}

	// Set expectations
	mockAvailabilityRepo.On("GetByUserAndEvent", mock.Anything, userID, eventID).Return(existingAvailabilities, nil)

	// Execute the method
	updatedAvailability, err := availabilityService.UpdateAvailability(context.Background(), eventID, req)

	// Assertions
	assert.Error(t, err)
	assert.Nil(t, updatedAvailability)
	assert.Equal(t, errors.ErrInvalidTimeRange, err)

	// Verify mock expectations
	mockAvailabilityRepo.AssertExpectations(t)
}

func TestDeleteAvailability(t *testing.T) {
	// Setup mocks
	mockAvailabilityRepo := new(MockAvailabilityRepository)
	mockEventRepo := new(MockEventRepository)
	mockUserRepo := new(MockUserRepository)
	availabilityService := service.NewAvailabilityService(
		mockAvailabilityRepo,
		mockEventRepo,
		mockUserRepo,
	)

	// Prepare test data
	availabilityID := uuid.New()

	// Set expectations
	mockAvailabilityRepo.On("Delete", mock.Anything, availabilityID).Return(nil)

	// Execute the method
	err := availabilityService.DeleteAvailability(context.Background(), availabilityID)

	// Assertions
	assert.NoError(t, err)

	// Verify mock expectations
	mockAvailabilityRepo.AssertExpectations(t)
}

func TestDeleteAvailabilityError(t *testing.T) {
	// Setup mocks
	mockAvailabilityRepo := new(MockAvailabilityRepository)
	mockEventRepo := new(MockEventRepository)
	mockUserRepo := new(MockUserRepository)
	availabilityService := service.NewAvailabilityService(
		mockAvailabilityRepo,
		mockEventRepo,
		mockUserRepo,
	)

	// Prepare test data
	availabilityID := uuid.New()

	// Set expectations
	mockAvailabilityRepo.On("Delete", mock.Anything, availabilityID).Return(assert.AnError)

	// Execute the method
	err := availabilityService.DeleteAvailability(context.Background(), availabilityID)

	// Assertions
	assert.Error(t, err)

	// Verify mock expectations
	mockAvailabilityRepo.AssertExpectations(t)
}

func TestGetUserEventAvailability(t *testing.T) {
	// Setup mocks
	mockAvailabilityRepo := new(MockAvailabilityRepository)
	mockEventRepo := new(MockEventRepository)
	mockUserRepo := new(MockUserRepository)
	availabilityService := service.NewAvailabilityService(
		mockAvailabilityRepo,
		mockEventRepo,
		mockUserRepo,
	)

	// Prepare test data
	eventID := uuid.New()
	userID := uuid.New()
	expectedAvailabilities := []*models.Availability{
		{
			ID:        uuid.New(),
			UserID:    userID,
			EventID:   eventID,
			StartTime: time.Date(2025, 1, 15, 10, 0, 0, 0, time.UTC),
			EndTime:   time.Date(2025, 1, 15, 12, 0, 0, 0, time.UTC),
		},
		{
			ID:        uuid.New(),
			UserID:    userID,
			EventID:   eventID,
			StartTime: time.Date(2025, 1, 15, 14, 0, 0, 0, time.UTC),
			EndTime:   time.Date(2025, 1, 15, 16, 0, 0, 0, time.UTC),
		},
	}

	// Set expectations
	mockAvailabilityRepo.On("GetByUserAndEvent", mock.Anything, userID, eventID).Return(expectedAvailabilities, nil)

	// Execute the method
	availabilities, err := availabilityService.GetUserEventAvailability(context.Background(), userID, eventID)

	// Assertions
	assert.NoError(t, err)
	assert.Len(t, availabilities, 2)
	assert.Equal(t, expectedAvailabilities, availabilities)

	// Verify mock expectations
	mockAvailabilityRepo.AssertExpectations(t)
}

func TestGetUserEventAvailabilityError(t *testing.T) {
	// Setup mocks
	mockAvailabilityRepo := new(MockAvailabilityRepository)
	mockEventRepo := new(MockEventRepository)
	mockUserRepo := new(MockUserRepository)
	availabilityService := service.NewAvailabilityService(
		mockAvailabilityRepo,
		mockEventRepo,
		mockUserRepo,
	)

	// Prepare test data
	eventID := uuid.New()
	userID := uuid.New()

	// Set expectations
	mockAvailabilityRepo.On("GetByUserAndEvent", mock.Anything, userID, eventID).Return(nil, assert.AnError)

	// Execute the method
	availabilities, err := availabilityService.GetUserEventAvailability(context.Background(), userID, eventID)

	// Assertions
	assert.Error(t, err)
	assert.Nil(t, availabilities)

	// Verify mock expectations
	mockAvailabilityRepo.AssertExpectations(t)
}

func TestGetEventAvailability(t *testing.T) {
	// Setup mocks
	mockAvailabilityRepo := new(MockAvailabilityRepository)
	mockEventRepo := new(MockEventRepository)
	mockUserRepo := new(MockUserRepository)
	availabilityService := service.NewAvailabilityService(
		mockAvailabilityRepo,
		mockEventRepo,
		mockUserRepo,
	)

	// Prepare test data
	eventID := uuid.New()
	expectedAvailabilities := []*models.Availability{
		{
			ID:        uuid.New(),
			UserID:    uuid.New(),
			EventID:   eventID,
			StartTime: time.Date(2025, 1, 15, 10, 0, 0, 0, time.UTC),
			EndTime:   time.Date(2025, 1, 15, 12, 0, 0, 0, time.UTC),
		},
		{
			ID:        uuid.New(),
			UserID:    uuid.New(),
			EventID:   eventID,
			StartTime: time.Date(2025, 1, 15, 14, 0, 0, 0, time.UTC),
			EndTime:   time.Date(2025, 1, 15, 16, 0, 0, 0, time.UTC),
		},
	}

	// Set expectations
	mockAvailabilityRepo.On("GetByEventID", mock.Anything, eventID).Return(expectedAvailabilities, nil)

	// Execute the method
	availabilities, err := availabilityService.GetEventAvailability(context.Background(), eventID)

	// Assertions
	assert.NoError(t, err)
	assert.Len(t, availabilities, 2)
	assert.Equal(t, expectedAvailabilities, availabilities)

	// Verify mock expectations
	mockAvailabilityRepo.AssertExpectations(t)
}

func TestGetEventAvailabilityError(t *testing.T) {
	// Setup mocks
	mockAvailabilityRepo := new(MockAvailabilityRepository)
	mockEventRepo := new(MockEventRepository)
	mockUserRepo := new(MockUserRepository)
	availabilityService := service.NewAvailabilityService(
		mockAvailabilityRepo,
		mockEventRepo,
		mockUserRepo,
	)

	// Prepare test data
	eventID := uuid.New()

	// Set expectations
	mockAvailabilityRepo.On("GetByEventID", mock.Anything, eventID).Return(nil, assert.AnError)

	// Execute the method
	availabilities, err := availabilityService.GetEventAvailability(context.Background(), eventID)

	// Assertions
	assert.Error(t, err)
	assert.Nil(t, availabilities)

	// Verify mock expectations
	mockAvailabilityRepo.AssertExpectations(t)
}
