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

func TestCreateTimeSlot(t *testing.T) {
	// Setup mocks
	mockTimeSlotRepo := new(MockTimeSlotRepository)
	mockEventRepo := new(MockEventRepository)
	timeSlotService := service.NewTimeSlotService(mockTimeSlotRepo, mockEventRepo)

	// Prepare test data
	eventID := uuid.New()
	startTime := time.Date(2025, 1, 15, 10, 0, 0, 0, time.UTC)
	endTime := time.Date(2025, 1, 15, 12, 0, 0, 0, time.UTC)

	// Mock event exists
	mockEventRepo.On("GetByID", mock.Anything, eventID).Return(
		&models.Event{ID: eventID},
		nil,
	)

	// Prepare request
	req := &models.TimeSlotRequest{
		StartTime: startTime.Format(time.RFC3339),
		EndTime:   endTime.Format(time.RFC3339),
	}

	// Set expectations for time slot creation
	mockTimeSlotRepo.On("Create", mock.Anything, mock.Anything).Return(nil)

	// Execute the method
	timeSlot, err := timeSlotService.CreateTimeSlot(context.Background(), eventID, req)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, timeSlot)
	assert.Equal(t, eventID, timeSlot.EventID)
	assert.Equal(t, startTime, timeSlot.StartTime)
	assert.Equal(t, endTime, timeSlot.EndTime)

	// Verify mock expectations
	mockEventRepo.AssertExpectations(t)
	mockTimeSlotRepo.AssertExpectations(t)
}

func TestCreateTimeSlotInvalidTimeRange(t *testing.T) {
	// Setup mocks
	mockTimeSlotRepo := new(MockTimeSlotRepository)
	mockEventRepo := new(MockEventRepository)
	timeSlotService := service.NewTimeSlotService(mockTimeSlotRepo, mockEventRepo)

	// Prepare test data
	eventID := uuid.New()
	startTime := time.Date(2025, 1, 15, 12, 0, 0, 0, time.UTC)
	endTime := time.Date(2025, 1, 15, 10, 0, 0, 0, time.UTC)

	// Mock event exists
	mockEventRepo.On("GetByID", mock.Anything, eventID).Return(
		&models.Event{ID: eventID},
		nil,
	)

	// Prepare request with invalid time range
	req := &models.TimeSlotRequest{
		StartTime: startTime.Format(time.RFC3339),
		EndTime:   endTime.Format(time.RFC3339),
	}

	// Execute the method
	timeSlot, err := timeSlotService.CreateTimeSlot(context.Background(), eventID, req)

	// Assertions
	assert.Error(t, err)
	assert.Nil(t, timeSlot)
	assert.Equal(t, errors.ErrInvalidTimeRange, err)

	// Verify mock expectations
	mockEventRepo.AssertExpectations(t)
	mockTimeSlotRepo.AssertExpectations(t)
}

func TestGetTimeSlot(t *testing.T) {
	// Setup mocks
	mockTimeSlotRepo := new(MockTimeSlotRepository)
	mockEventRepo := new(MockEventRepository)
	timeSlotService := service.NewTimeSlotService(mockTimeSlotRepo, mockEventRepo)

	// Prepare test data
	timeSlotID := uuid.New()
	expectedTimeSlot := &models.TimeSlot{
		ID:        timeSlotID,
		EventID:   uuid.New(),
		StartTime: time.Date(2025, 1, 15, 10, 0, 0, 0, time.UTC),
		EndTime:   time.Date(2025, 1, 15, 12, 0, 0, 0, time.UTC),
	}

	// Set expectations
	mockTimeSlotRepo.On("GetByID", mock.Anything, timeSlotID).Return(expectedTimeSlot, nil)

	// Execute the method
	timeSlot, err := timeSlotService.GetTimeSlot(context.Background(), timeSlotID)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, timeSlot)
	assert.Equal(t, expectedTimeSlot, timeSlot)

	// Verify mock expectations
	mockTimeSlotRepo.AssertExpectations(t)
}

func TestUpdateTimeSlot(t *testing.T) {
	// Setup mocks
	mockTimeSlotRepo := new(MockTimeSlotRepository)
	mockEventRepo := new(MockEventRepository)
	timeSlotService := service.NewTimeSlotService(mockTimeSlotRepo, mockEventRepo)

	// Prepare test data
	timeSlotID := uuid.New()
	existingTimeSlot := &models.TimeSlot{
		ID:        timeSlotID,
		EventID:   uuid.New(),
		StartTime: time.Date(2025, 1, 15, 10, 0, 0, 0, time.UTC),
		EndTime:   time.Date(2025, 1, 15, 12, 0, 0, 0, time.UTC),
	}

	// New time range
	newStartTime := time.Date(2025, 1, 15, 14, 0, 0, 0, time.UTC)
	newEndTime := time.Date(2025, 1, 15, 16, 0, 0, 0, time.UTC)

	// Prepare request
	req := &models.TimeSlotRequest{
		StartTime: newStartTime.Format(time.RFC3339),
		EndTime:   newEndTime.Format(time.RFC3339),
	}

	// Set expectations
	mockTimeSlotRepo.On("GetByID", mock.Anything, timeSlotID).Return(existingTimeSlot, nil)
	mockTimeSlotRepo.On("Update", mock.Anything, mock.Anything).Return(nil)

	// Execute the method
	updatedTimeSlot, err := timeSlotService.UpdateTimeSlot(context.Background(), timeSlotID, req)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, updatedTimeSlot)
	assert.Equal(t, newStartTime, updatedTimeSlot.StartTime)
	assert.Equal(t, newEndTime, updatedTimeSlot.EndTime)
	assert.WithinDuration(t, time.Now(), updatedTimeSlot.UpdatedAt, time.Second)

	// Verify mock expectations
	mockTimeSlotRepo.AssertExpectations(t)
}

func TestDeleteTimeSlot(t *testing.T) {
	// Setup mocks
	mockTimeSlotRepo := new(MockTimeSlotRepository)
	mockEventRepo := new(MockEventRepository)
	timeSlotService := service.NewTimeSlotService(mockTimeSlotRepo, mockEventRepo)

	// Prepare test data
	timeSlotID := uuid.New()

	// Set expectations
	mockTimeSlotRepo.On("Delete", mock.Anything, timeSlotID).Return(nil)

	// Execute the method
	err := timeSlotService.DeleteTimeSlot(context.Background(), timeSlotID)

	// Assertions
	assert.NoError(t, err)

	// Verify mock expectations
	mockTimeSlotRepo.AssertExpectations(t)
}

func TestGetEventTimeSlots(t *testing.T) {
	// Setup mocks
	mockTimeSlotRepo := new(MockTimeSlotRepository)
	mockEventRepo := new(MockEventRepository)
	timeSlotService := service.NewTimeSlotService(mockTimeSlotRepo, mockEventRepo)

	// Prepare test data
	eventID := uuid.New()
	expectedTimeSlots := []*models.TimeSlot{
		{
			ID:        uuid.New(),
			EventID:   eventID,
			StartTime: time.Date(2025, 1, 15, 10, 0, 0, 0, time.UTC),
			EndTime:   time.Date(2025, 1, 15, 12, 0, 0, 0, time.UTC),
		},
		{
			ID:        uuid.New(),
			EventID:   eventID,
			StartTime: time.Date(2025, 1, 15, 14, 0, 0, 0, time.UTC),
			EndTime:   time.Date(2025, 1, 15, 16, 0, 0, 0, time.UTC),
		},
	}

	// Set expectations
	mockTimeSlotRepo.On("GetByEventID", mock.Anything, eventID).Return(expectedTimeSlots, nil)

	// Execute the method
	timeSlots, err := timeSlotService.GetEventTimeSlots(context.Background(), eventID)

	// Assertions
	assert.NoError(t, err)
	assert.Len(t, timeSlots, 2)
	assert.Equal(t, expectedTimeSlots, timeSlots)

	// Verify mock expectations
	mockTimeSlotRepo.AssertExpectations(t)
}
