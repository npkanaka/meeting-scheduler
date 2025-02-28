// internal/service/recommendation_service_test.go
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

// MockEventRepository is a mock for the EventRepository
type MockEventRepository struct {
	mock.Mock
}

func (m *MockEventRepository) Create(ctx context.Context, event *models.Event) error {
	args := m.Called(ctx, event)
	return args.Error(0)
}

func (m *MockEventRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Event, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Event), args.Error(1)
}

func (m *MockEventRepository) Update(ctx context.Context, event *models.Event) error {
	args := m.Called(ctx, event)
	return args.Error(0)
}

func (m *MockEventRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockEventRepository) List(ctx context.Context, limit, offset int) ([]*models.Event, error) {
	args := m.Called(ctx, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Event), args.Error(1)
}

// MockTimeSlotRepository is a mock for the TimeSlotRepository
type MockTimeSlotRepository struct {
	mock.Mock
}

func (m *MockTimeSlotRepository) Create(ctx context.Context, slot *models.TimeSlot) error {
	args := m.Called(ctx, slot)
	return args.Error(0)
}

func (m *MockTimeSlotRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.TimeSlot, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.TimeSlot), args.Error(1)
}

func (m *MockTimeSlotRepository) Update(ctx context.Context, slot *models.TimeSlot) error {
	args := m.Called(ctx, slot)
	return args.Error(0)
}

func (m *MockTimeSlotRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockTimeSlotRepository) GetByEventID(ctx context.Context, eventID uuid.UUID) ([]*models.TimeSlot, error) {
	args := m.Called(ctx, eventID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.TimeSlot), args.Error(1)
}

// MockAvailabilityRepository is a mock for the AvailabilityRepository
type MockAvailabilityRepository struct {
	mock.Mock
}

func (m *MockAvailabilityRepository) Create(ctx context.Context, availability *models.Availability) error {
	args := m.Called(ctx, availability)
	return args.Error(0)
}

func (m *MockAvailabilityRepository) Update(ctx context.Context, availability *models.Availability) error {
	args := m.Called(ctx, availability)
	return args.Error(0)
}

func (m *MockAvailabilityRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockAvailabilityRepository) GetByUserAndEvent(ctx context.Context, userID, eventID uuid.UUID) ([]*models.Availability, error) {
	args := m.Called(ctx, userID, eventID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Availability), args.Error(1)
}

func (m *MockAvailabilityRepository) GetByEventID(ctx context.Context, eventID uuid.UUID) ([]*models.Availability, error) {
	args := m.Called(ctx, eventID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Availability), args.Error(1)
}

// MockUserRepository is a mock for the UserRepository
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(ctx context.Context, user *models.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) GetByIDs(ctx context.Context, ids []uuid.UUID) ([]*models.User, error) {
	args := m.Called(ctx, ids)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.User), args.Error(1)
}

func (m *MockUserRepository) Update(ctx context.Context, user *models.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestGetRecommendations(t *testing.T) {
	// Setup mocks
	mockEventRepo := new(MockEventRepository)
	mockTimeSlotRepo := new(MockTimeSlotRepository)
	mockAvailabilityRepo := new(MockAvailabilityRepository)
	mockUserRepo := new(MockUserRepository)

	// Create the service with mocks
	recommendationService := service.NewRecommendationService(
		mockEventRepo,
		mockTimeSlotRepo,
		mockAvailabilityRepo,
		mockUserRepo,
	)

	ctx := context.Background()

	// Setup test data
	eventID := uuid.New()
	user1ID := uuid.New()
	user2ID := uuid.New()

	// Create a test event with 1-hour duration
	testEvent := &models.Event{
		ID:          eventID,
		Title:       "Test Meeting",
		Description: "Test Description",
		CreatorID:   user1ID,
		Duration:    60, // 1 hour meeting
		Status:      models.EventStatusActive,
	}

	// Create test time slots
	slot1Start := time.Date(2025, 1, 15, 10, 0, 0, 0, time.UTC)
	slot1End := time.Date(2025, 1, 15, 12, 0, 0, 0, time.UTC)

	slot2Start := time.Date(2025, 1, 15, 14, 0, 0, 0, time.UTC)
	slot2End := time.Date(2025, 1, 15, 16, 0, 0, 0, time.UTC)

	testTimeSlots := []*models.TimeSlot{
		{
			ID:        uuid.New(),
			EventID:   eventID,
			StartTime: slot1Start,
			EndTime:   slot1End,
		},
		{
			ID:        uuid.New(),
			EventID:   eventID,
			StartTime: slot2Start,
			EndTime:   slot2End,
		},
	}

	// Create test availability
	// User1 is available for both slots
	user1Availability := []*models.Availability{
		{
			ID:        uuid.New(),
			UserID:    user1ID,
			EventID:   eventID,
			StartTime: slot1Start,
			EndTime:   slot1End,
		},
		{
			ID:        uuid.New(),
			UserID:    user1ID,
			EventID:   eventID,
			StartTime: slot2Start,
			EndTime:   slot2End,
		},
	}

	// User2 is only available for the second slot
	user2Availability := []*models.Availability{
		{
			ID:        uuid.New(),
			UserID:    user2ID,
			EventID:   eventID,
			StartTime: slot2Start,
			EndTime:   slot2End,
		},
	}

	allAvailability := append(user1Availability, user2Availability...)

	// Create test users
	user1 := &models.User{
		ID:    user1ID,
		Name:  "User 1",
		Email: "user1@example.com",
	}

	user2 := &models.User{
		ID:    user2ID,
		Name:  "User 2",
		Email: "user2@example.com",
	}

	testUsers := []*models.User{user1, user2}

	// Setup expectations
	mockEventRepo.On("GetByID", ctx, eventID).Return(testEvent, nil)
	mockTimeSlotRepo.On("GetByEventID", ctx, eventID).Return(testTimeSlots, nil)
	mockAvailabilityRepo.On("GetByEventID", ctx, eventID).Return(allAvailability, nil)
	mockUserRepo.On("GetByIDs", ctx, mock.MatchedBy(func(ids []uuid.UUID) bool {
		// Check if the slice contains both user IDs, regardless of order
		if len(ids) != 2 {
			return false
		}
		hasUser1 := false
		hasUser2 := false
		for _, id := range ids {
			if id == user1ID {
				hasUser1 = true
			}
			if id == user2ID {
				hasUser2 = true
			}
		}
		return hasUser1 && hasUser2
	})).Return(testUsers, nil)

	// Execute the function being tested
	recommendations, err := recommendationService.GetRecommendations(ctx, eventID)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, recommendations)
	assert.Len(t, recommendations.Recommendations, 2)

	// The second slot should have a higher score (2 attendees vs 1)
	assert.Equal(t, 2, recommendations.Recommendations[0].Score)
	assert.Equal(t, testTimeSlots[1].ID, recommendations.Recommendations[0].TimeSlot.ID)
	assert.Len(t, recommendations.Recommendations[0].Attendees, 2)
	assert.Len(t, recommendations.Recommendations[0].NonAttendees, 0)

	assert.Equal(t, 1, recommendations.Recommendations[1].Score)
	assert.Equal(t, testTimeSlots[0].ID, recommendations.Recommendations[1].TimeSlot.ID)
	assert.Len(t, recommendations.Recommendations[1].Attendees, 1)
	assert.Len(t, recommendations.Recommendations[1].NonAttendees, 1)

	// Verify all the mocks were called as expected
	mockEventRepo.AssertExpectations(t)
	mockTimeSlotRepo.AssertExpectations(t)
	mockAvailabilityRepo.AssertExpectations(t)
	mockUserRepo.AssertExpectations(t)
}
