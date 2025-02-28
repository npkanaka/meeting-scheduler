package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/npkanaka/meeting-scheduler/internal/models"
	"github.com/npkanaka/meeting-scheduler/internal/service"
)

// TimeSlotHandler handles HTTP requests related to time slots
type TimeSlotHandler struct {
	timeSlotService *service.TimeSlotService
}

// NewTimeSlotHandler creates a new TimeSlotHandler
func NewTimeSlotHandler(timeSlotService *service.TimeSlotService) *TimeSlotHandler {
	return &TimeSlotHandler{
		timeSlotService: timeSlotService,
	}
}

// Create handles the creation of a new time slot
func (h *TimeSlotHandler) Create(c *gin.Context) {
	idStr := c.Param("id") // Changed from eventId to id
	eventID, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid event ID"})
		return
	}

	var req models.TimeSlotRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	timeSlot, err := h.timeSlotService.CreateTimeSlot(c.Request.Context(), eventID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, timeSlot)
}

// Get retrieves a time slot by ID
func (h *TimeSlotHandler) Get(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid time slot ID"})
		return
	}

	timeSlot, err := h.timeSlotService.GetTimeSlot(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, timeSlot)
}

// Update updates an existing time slot
func (h *TimeSlotHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid time slot ID"})
		return
	}

	var req models.TimeSlotRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	timeSlot, err := h.timeSlotService.UpdateTimeSlot(c.Request.Context(), id, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, timeSlot)
}

// Delete removes a time slot
func (h *TimeSlotHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid time slot ID"})
		return
	}

	if err := h.timeSlotService.DeleteTimeSlot(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// List returns all time slots for an event
func (h *TimeSlotHandler) List(c *gin.Context) {
	idStr := c.Param("id") // Changed from eventId to id
	eventID, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid event ID"})
		return
	}

	timeSlots, err := h.timeSlotService.GetEventTimeSlots(c.Request.Context(), eventID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"time_slots": timeSlots})
}
