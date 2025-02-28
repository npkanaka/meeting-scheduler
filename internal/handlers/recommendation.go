package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/npkanaka/meeting-scheduler/internal/service"
)

// RecommendationHandler handles HTTP requests related to time slot recommendations
type RecommendationHandler struct {
	recommendationService *service.RecommendationService
}

// NewRecommendationHandler creates a new RecommendationHandler
func NewRecommendationHandler(recommendationService *service.RecommendationService) *RecommendationHandler {
	return &RecommendationHandler{
		recommendationService: recommendationService,
	}
}

// GetRecommendations returns recommended time slots for an event
func (h *RecommendationHandler) GetRecommendations(c *gin.Context) {
	idStr := c.Param("id") // Changed from eventId to id
	eventID, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid event ID"})
		return
	}

	recommendations, err := h.recommendationService.GetRecommendations(c.Request.Context(), eventID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, recommendations)
}
