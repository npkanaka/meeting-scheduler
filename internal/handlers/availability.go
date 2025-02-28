package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/npkanaka/meeting-scheduler/internal/models"
	"github.com/npkanaka/meeting-scheduler/internal/service"
)

// AvailabilityHandler handles HTTP requests related to user availability
type AvailabilityHandler struct {
	availabilityService *service.AvailabilityService
}

// NewAvailabilityHandler creates a new AvailabilityHandler
func NewAvailabilityHandler(availabilityService *service.AvailabilityService) *AvailabilityHandler {
	return &AvailabilityHandler{
		availabilityService: availabilityService,
	}
}

// Create handles the creation of a new availability record
func (h *AvailabilityHandler) Create(c *gin.Context) {
	idStr := c.Param("id") // Changed from eventId to id
	eventID, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid event ID"})
		return
	}

	var req models.AvailabilityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	availability, err := h.availabilityService.CreateAvailability(c.Request.Context(), eventID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, availability)
}

// Update updates an existing availability record
func (h *AvailabilityHandler) Update(c *gin.Context) {
	idStr := c.Param("id") // Changed from eventId to id
	eventID, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid event ID"})
		return
	}

	userIDStr := c.Param("userId")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	var req models.AvailabilityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	req.UserID = userID
	availability, err := h.availabilityService.UpdateAvailability(c.Request.Context(), eventID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, availability)
}

// Delete removes an availability record
func (h *AvailabilityHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid availability ID"})
		return
	}

	if err := h.availabilityService.DeleteAvailability(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// GetUserAvailability retrieves all availability records for a user and event
func (h *AvailabilityHandler) GetUserAvailability(c *gin.Context) {
	idStr := c.Param("id") // Changed from eventId to id
	eventID, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid event ID"})
		return
	}

	userIDStr := c.Param("userId")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	availabilities, err := h.availabilityService.GetUserEventAvailability(c.Request.Context(), userID, eventID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"availabilities": availabilities})
}

// GetEventAvailability retrieves all availability records for an event
func (h *AvailabilityHandler) GetEventAvailability(c *gin.Context) {
	idStr := c.Param("id") // Changed from eventId to id
	eventID, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid event ID"})
		return
	}

	availabilities, err := h.availabilityService.GetEventAvailability(c.Request.Context(), eventID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"availabilities": availabilities})
}
