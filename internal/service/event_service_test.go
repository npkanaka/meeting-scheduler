package service_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/npkanaka/meeting-scheduler/internal/models"
	"github.com/npkanaka/meeting-scheduler/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateEvent(t *testing.T) {
	// Setup mock repository
	mockEventRepo := new(MockEventRepository)
	eventService := service.NewEventService(mockEventRepo)

	// Prepare test data
	creatorID := uuid.New()
	req := &models.CreateEventRequest{
		Title:       "Team Brainstorming",
		Description: "Quarterly team brainstorming session",
		Duration:    60, // 1-hour meeting
	}

	// Set expectations for event creation
	mockEventRepo.On("Create", mock.Anything, mock.MatchedBy(func(event *models.Event) bool {
		return event.Title == req.Title &&
			event.Description == req.Description &&
			event.Duration == req.Duration &&
			event.CreatorID == creatorID &&
			event.Status == models.EventStatusDraft
	})).Return(nil)

	// Execute the method
	event, err := eventService.CreateEvent(context.Background(), req, creatorID)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, event)
	assert.Equal(t, req.Title, event.Title)
	assert.Equal(t, req.Description, event.Description)
	assert.Equal(t, req.Duration, event.Duration)
	assert.Equal(t, creatorID, event.CreatorID)
	assert.Equal(t, models.EventStatusDraft, event.Status)
	assert.NotZero(t, event.CreatedAt)
	assert.NotZero(t, event.UpdatedAt)

	// Verify mock expectations
	mockEventRepo.AssertExpectations(t)
}

func TestCreateEventRepositoryError(t *testing.T) {
	// Setup mock repository
	mockEventRepo := new(MockEventRepository)
	eventService := service.NewEventService(mockEventRepo)

	// Prepare test data
	creatorID := uuid.New()
	req := &models.CreateEventRequest{
		Title:    "Team Meeting",
		Duration: 45,
	}

	// Set expectations for event creation with error
	mockEventRepo.On("Create", mock.Anything, mock.Anything).Return(assert.AnError)

	// Execute the method
	event, err := eventService.CreateEvent(context.Background(), req, creatorID)

	// Assertions
	assert.Error(t, err)
	assert.Nil(t, event)
	assert.Equal(t, assert.AnError, err)

	// Verify mock expectations
	mockEventRepo.AssertExpectations(t)
}

func TestGetEvent(t *testing.T) {
	// Setup mock repository
	mockEventRepo := new(MockEventRepository)
	eventService := service.NewEventService(mockEventRepo)

	// Prepare test data
	eventID := uuid.New()
	expectedEvent := &models.Event{
		ID:          eventID,
		Title:       "Product Review",
		Description: "Monthly product review meeting",
		CreatorID:   uuid.New(),
		Duration:    90,
		Status:      models.EventStatusActive,
		CreatedAt:   time.Now().Add(-24 * time.Hour),
		UpdatedAt:   time.Now(),
	}

	// Set expectations
	mockEventRepo.On("GetByID", mock.Anything, eventID).Return(expectedEvent, nil)

	// Execute the method
	event, err := eventService.GetEvent(context.Background(), eventID)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, event)
	assert.Equal(t, expectedEvent, event)

	// Verify mock expectations
	mockEventRepo.AssertExpectations(t)
}

func TestGetEventNotFound(t *testing.T) {
	// Setup mock repository
	mockEventRepo := new(MockEventRepository)
	eventService := service.NewEventService(mockEventRepo)

	// Prepare test data
	eventID := uuid.New()

	// Set expectations
	mockEventRepo.On("GetByID", mock.Anything, eventID).Return(nil, assert.AnError)

	// Execute the method
	event, err := eventService.GetEvent(context.Background(), eventID)

	// Assertions
	assert.Error(t, err)
	assert.Nil(t, event)

	// Verify mock expectations
	mockEventRepo.AssertExpectations(t)
}

func TestUpdateEvent(t *testing.T) {
	// Setup mock repository
	mockEventRepo := new(MockEventRepository)
	eventService := service.NewEventService(mockEventRepo)

	// Prepare test data
	eventID := uuid.New()
	existingEvent := &models.Event{
		ID:          eventID,
		Title:       "Old Event Title",
		Description: "Old event description",
		Duration:    60,
	}

	// Update request
	updateReq := &models.CreateEventRequest{
		Title:       "Updated Event Title",
		Description: "Updated event description",
		Duration:    90,
	}

	// Set expectations
	mockEventRepo.On("GetByID", mock.Anything, eventID).Return(existingEvent, nil)
	mockEventRepo.On("Update", mock.Anything, mock.MatchedBy(func(event *models.Event) bool {
		return event.Title == updateReq.Title &&
			event.Description == updateReq.Description &&
			event.Duration == updateReq.Duration
	})).Return(nil)

	// Execute the method
	updatedEvent, err := eventService.UpdateEvent(context.Background(), eventID, updateReq)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, updatedEvent)
	assert.Equal(t, updateReq.Title, updatedEvent.Title)
	assert.Equal(t, updateReq.Description, updatedEvent.Description)
	assert.Equal(t, updateReq.Duration, updatedEvent.Duration)
	assert.WithinDuration(t, time.Now(), updatedEvent.UpdatedAt, time.Second)

	// Verify mock expectations
	mockEventRepo.AssertExpectations(t)
}

func TestUpdateEventNotFound(t *testing.T) {
	// Setup mock repository
	mockEventRepo := new(MockEventRepository)
	eventService := service.NewEventService(mockEventRepo)

	// Prepare test data
	eventID := uuid.New()
	updateReq := &models.CreateEventRequest{
		Title: "Updated Event",
	}

	// Set expectations
	mockEventRepo.On("GetByID", mock.Anything, eventID).Return(nil, assert.AnError)

	// Execute the method
	updatedEvent, err := eventService.UpdateEvent(context.Background(), eventID, updateReq)

	// Assertions
	assert.Error(t, err)
	assert.Nil(t, updatedEvent)

	// Verify mock expectations
	mockEventRepo.AssertExpectations(t)
}

func TestUpdateEventRepositoryError(t *testing.T) {
	// Setup mock repository
	mockEventRepo := new(MockEventRepository)
	eventService := service.NewEventService(mockEventRepo)

	// Prepare test data
	eventID := uuid.New()
	existingEvent := &models.Event{
		ID:    eventID,
		Title: "Existing Event",
	}

	// Update request
	updateReq := &models.CreateEventRequest{
		Title: "Updated Event",
	}

	// Set expectations
	mockEventRepo.On("GetByID", mock.Anything, eventID).Return(existingEvent, nil)
	mockEventRepo.On("Update", mock.Anything, mock.Anything).Return(assert.AnError)

	// Execute the method
	updatedEvent, err := eventService.UpdateEvent(context.Background(), eventID, updateReq)

	// Assertions
	assert.Error(t, err)
	assert.Nil(t, updatedEvent)

	// Verify mock expectations
	mockEventRepo.AssertExpectations(t)
}

func TestDeleteEvent(t *testing.T) {
	// Setup mock repository
	mockEventRepo := new(MockEventRepository)
	eventService := service.NewEventService(mockEventRepo)

	// Prepare test data
	eventID := uuid.New()

	// Set expectations
	mockEventRepo.On("Delete", mock.Anything, eventID).Return(nil)

	// Execute the method
	err := eventService.DeleteEvent(context.Background(), eventID)

	// Assertions
	assert.NoError(t, err)

	// Verify mock expectations
	mockEventRepo.AssertExpectations(t)
}

func TestDeleteEventRepositoryError(t *testing.T) {
	// Setup mock repository
	mockEventRepo := new(MockEventRepository)
	eventService := service.NewEventService(mockEventRepo)

	// Prepare test data
	eventID := uuid.New()

	// Set expectations
	mockEventRepo.On("Delete", mock.Anything, eventID).Return(assert.AnError)

	// Execute the method
	err := eventService.DeleteEvent(context.Background(), eventID)

	// Assertions
	assert.Error(t, err)

	// Verify mock expectations
	mockEventRepo.AssertExpectations(t)
}

func TestListEvents(t *testing.T) {
	// Setup mock repository
	mockEventRepo := new(MockEventRepository)
	eventService := service.NewEventService(mockEventRepo)

	// Prepare test data
	expectedEvents := []*models.Event{
		{
			ID:    uuid.New(),
			Title: "Event 1",
		},
		{
			ID:    uuid.New(),
			Title: "Event 2",
		},
	}

	// Set expectations
	mockEventRepo.On("List", mock.Anything, 10, 0).Return(expectedEvents, nil)

	// Execute the method
	events, err := eventService.ListEvents(context.Background(), 10, 0)

	// Assertions
	assert.NoError(t, err)
	assert.Len(t, events, 2)
	assert.Equal(t, expectedEvents, events)

	// Verify mock expectations
	mockEventRepo.AssertExpectations(t)
}

func TestListEventsRepositoryError(t *testing.T) {
	// Setup mock repository
	mockEventRepo := new(MockEventRepository)
	eventService := service.NewEventService(mockEventRepo)

	// Set expectations
	mockEventRepo.On("List", mock.Anything, 10, 0).Return(nil, assert.AnError)

	// Execute the method
	events, err := eventService.ListEvents(context.Background(), 10, 0)

	// Assertions
	assert.Error(t, err)
	assert.Nil(t, events)

	// Verify mock expectations
	mockEventRepo.AssertExpectations(t)
}
