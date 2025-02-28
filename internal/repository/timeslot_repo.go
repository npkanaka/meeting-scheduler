package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/npkanaka/meeting-scheduler/internal/models"
	"gorm.io/gorm"
)

// TimeSlotRepository defines the interface for time slot data access
type TimeSlotRepository interface {
	Create(ctx context.Context, slot *models.TimeSlot) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.TimeSlot, error)
	Update(ctx context.Context, slot *models.TimeSlot) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetByEventID(ctx context.Context, eventID uuid.UUID) ([]*models.TimeSlot, error)
}

// GormTimeSlotRepository implements TimeSlotRepository using GORM
type GormTimeSlotRepository struct {
	db *gorm.DB
}

// NewGormTimeSlotRepository creates a new GormTimeSlotRepository
func NewGormTimeSlotRepository(db *gorm.DB) *GormTimeSlotRepository {
	return &GormTimeSlotRepository{db: db}
}

// Create saves a new time slot to the database
func (r *GormTimeSlotRepository) Create(ctx context.Context, slot *models.TimeSlot) error {
	if slot.ID == uuid.Nil {
		slot.ID = uuid.New()
	}
	return r.db.WithContext(ctx).Create(slot).Error
}

// GetByID retrieves a time slot by its ID
func (r *GormTimeSlotRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.TimeSlot, error) {
	var slot models.TimeSlot
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&slot).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("time slot not found")
		}
		return nil, err
	}
	return &slot, nil
}

// Update updates an existing time slot
func (r *GormTimeSlotRepository) Update(ctx context.Context, slot *models.TimeSlot) error {
	return r.db.WithContext(ctx).Model(&models.TimeSlot{}).Where("id = ?", slot.ID).Updates(slot).Error
}

// Delete removes a time slot by its ID
func (r *GormTimeSlotRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&models.TimeSlot{}, id).Error
}

// GetByEventID retrieves all time slots for an event
func (r *GormTimeSlotRepository) GetByEventID(ctx context.Context, eventID uuid.UUID) ([]*models.TimeSlot, error) {
	var slots []*models.TimeSlot
	err := r.db.WithContext(ctx).Where("event_id = ?", eventID).Find(&slots).Error
	return slots, err
}
