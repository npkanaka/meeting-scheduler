// internal/service/recommendation_service.go
package service

import (
	"context"
	"sort"
	"time"

	"github.com/google/uuid"
	"github.com/npkanaka/meeting-scheduler/internal/models"
	"github.com/npkanaka/meeting-scheduler/internal/repository"
)

// RecommendationService handles meeting time recommendations
type RecommendationService struct {
	eventRepo        repository.EventRepository
	timeslotRepo     repository.TimeSlotRepository
	availabilityRepo repository.AvailabilityRepository
	userRepo         repository.UserRepository
}

// NewRecommendationService creates a new RecommendationService
func NewRecommendationService(
	eventRepo repository.EventRepository,
	timeslotRepo repository.TimeSlotRepository,
	availabilityRepo repository.AvailabilityRepository,
	userRepo repository.UserRepository,
) *RecommendationService {
	return &RecommendationService{
		eventRepo:        eventRepo,
		timeslotRepo:     timeslotRepo,
		availabilityRepo: availabilityRepo,
		userRepo:         userRepo,
	}
}

// GetRecommendations generates time slot recommendations for an event
func (s *RecommendationService) GetRecommendations(ctx context.Context, eventID uuid.UUID) (*models.RecommendationResponse, error) {
	// Get the event
	event, err := s.eventRepo.GetByID(ctx, eventID)
	if err != nil {
		return nil, err
	}

	// Get all proposed time slots for the event
	timeSlots, err := s.timeslotRepo.GetByEventID(ctx, eventID)
	if err != nil {
		return nil, err
	}

	if len(timeSlots) == 0 {
		return &models.RecommendationResponse{Recommendations: []models.Recommendation{}}, nil
	}

	// Get all availability data for the event
	availabilities, err := s.availabilityRepo.GetByEventID(ctx, eventID)
	if err != nil {
		return nil, err
	}

	// Group availabilities by user
	userAvailabilities := make(map[uuid.UUID][]*models.Availability)
	userIDs := make(map[uuid.UUID]bool)

	for _, avail := range availabilities {
		userAvailabilities[avail.UserID] = append(userAvailabilities[avail.UserID], avail)
		userIDs[avail.UserID] = true
	}

	// Get all unique users
	var uniqueUserIDs []uuid.UUID
	for userID := range userIDs {
		uniqueUserIDs = append(uniqueUserIDs, userID)
	}

	users, err := s.userRepo.GetByIDs(ctx, uniqueUserIDs)
	if err != nil {
		return nil, err
	}

	// Create a map of user IDs to User objects for quick lookup
	userMap := make(map[uuid.UUID]*models.User)
	for _, user := range users {
		userMap[user.ID] = user
	}

	// Calculate recommendations
	var recommendations []models.Recommendation

	for _, slot := range timeSlots {
		// Calculate meeting end time based on event duration
		meetingEndTime := slot.StartTime.Add(time.Duration(event.Duration) * time.Minute)

		// If the meeting would extend beyond the slot's end time, skip this slot
		if meetingEndTime.After(slot.EndTime) {
			continue
		}

		var attendees []models.UserResponse
		var nonAttendees []models.UserResponse

		// Check each user's availability for this slot
		for userID, userAvails := range userAvailabilities {
			user, exists := userMap[userID]
			if !exists {
				continue
			}

			userResponse := models.UserResponse{
				ID:    user.ID,
				Name:  user.Name,
				Email: user.Email,
			}

			// Check if the user is available for this slot
			available := false
			for _, avail := range userAvails {
				// The user is available if the slot starts after or at the same time as their availability start
				// and the meeting ends before or at the same time as their availability end
				if (slot.StartTime.Equal(avail.StartTime) || slot.StartTime.After(avail.StartTime)) &&
					(meetingEndTime.Equal(avail.EndTime) || meetingEndTime.Before(avail.EndTime)) {
					available = true
					break
				}
			}

			if available {
				attendees = append(attendees, userResponse)
			} else {
				nonAttendees = append(nonAttendees, userResponse)
			}
		}

		// Create a recommendation
		timeSlotResponse := models.TimeSlotResponse{
			ID:        slot.ID,
			StartTime: slot.StartTime,
			EndTime:   meetingEndTime, // Use the calculated meeting end time
		}

		recommendations = append(recommendations, models.Recommendation{
			TimeSlot:     timeSlotResponse,
			Attendees:    attendees,
			NonAttendees: nonAttendees,
			Score:        len(attendees),
		})
	}

	// Sort recommendations by score (number of attendees) in descending order
	sort.Slice(recommendations, func(i, j int) bool {
		return recommendations[i].Score > recommendations[j].Score
	})

	return &models.RecommendationResponse{
		Recommendations: recommendations,
	}, nil
}
