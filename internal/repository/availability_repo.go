package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/npkanaka/meeting-scheduler/internal/models"
	"gorm.io/gorm"
)

// AvailabilityRepository defines the interface for availability data access
type AvailabilityRepository interface {
	Create(ctx context.Context, availability *models.Availability) error
	Update(ctx context.Context, availability *models.Availability) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetByUserAndEvent(ctx context.Context, userID, eventID uuid.UUID) ([]*models.Availability, error)
	GetByEventID(ctx context.Context, eventID uuid.UUID) ([]*models.Availability, error)
}

// GormAvailabilityRepository implements AvailabilityRepository using GORM
type GormAvailabilityRepository struct {
	db *gorm.DB
}

// NewGormAvailabilityRepository creates a new GormAvailabilityRepository
func NewGormAvailabilityRepository(db *gorm.DB) *GormAvailabilityRepository {
	return &GormAvailabilityRepository{db: db}
}

// Create saves a new availability to the database
func (r *GormAvailabilityRepository) Create(ctx context.Context, availability *models.Availability) error {
	if availability.ID == uuid.Nil {
		availability.ID = uuid.New()
	}
	return r.db.WithContext(ctx).Create(availability).Error
}

// Update updates an existing availability
func (r *GormAvailabilityRepository) Update(ctx context.Context, availability *models.Availability) error {
	return r.db.WithContext(ctx).Model(&models.Availability{}).Where("id = ?", availability.ID).Updates(availability).Error
}

// Delete removes an availability by its ID
func (r *GormAvailabilityRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&models.Availability{}, id).Error
}

// GetByUserAndEvent retrieves all availability entries for a user and event
func (r *GormAvailabilityRepository) GetByUserAndEvent(ctx context.Context, userID, eventID uuid.UUID) ([]*models.Availability, error) {
	var availabilities []*models.Availability
	err := r.db.WithContext(ctx).Where("user_id = ? AND event_id = ?", userID, eventID).Find(&availabilities).Error
	return availabilities, err
}

// GetByEventID retrieves all availability entries for an event
func (r *GormAvailabilityRepository) GetByEventID(ctx context.Context, eventID uuid.UUID) ([]*models.Availability, error) {
	var availabilities []*models.Availability
	err := r.db.WithContext(ctx).Where("event_id = ?", eventID).Find(&availabilities).Error
	return availabilities, err
}
