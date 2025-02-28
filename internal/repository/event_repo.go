package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/npkanaka/meeting-scheduler/internal/models"
	"gorm.io/gorm"
)

// EventRepository defines the interface for event data access
type EventRepository interface {
	Create(ctx context.Context, event *models.Event) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Event, error)
	Update(ctx context.Context, event *models.Event) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, limit, offset int) ([]*models.Event, error)
}

// GormEventRepository implements EventRepository using GORM
type GormEventRepository struct {
	db *gorm.DB
}

// NewGormEventRepository creates a new GormEventRepository
func NewGormEventRepository(db *gorm.DB) *GormEventRepository {
	return &GormEventRepository{db: db}
}

// Create saves a new event to the database
func (r *GormEventRepository) Create(ctx context.Context, event *models.Event) error {
	if event.ID == uuid.Nil {
		event.ID = uuid.New()
	}
	return r.db.WithContext(ctx).Create(event).Error
}

// GetByID retrieves an event by its ID
func (r *GormEventRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Event, error) {
	var event models.Event
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&event).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("event not found")
		}
		return nil, err
	}
	return &event, nil
}

// Update updates an existing event
func (r *GormEventRepository) Update(ctx context.Context, event *models.Event) error {
	return r.db.WithContext(ctx).Model(&models.Event{}).Where("id = ?", event.ID).Updates(event).Error
}

// Delete removes an event by its ID
func (r *GormEventRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&models.Event{}, id).Error
}

// List retrieves a paginated list of events
func (r *GormEventRepository) List(ctx context.Context, limit, offset int) ([]*models.Event, error) {
	var events []*models.Event
	err := r.db.WithContext(ctx).Limit(limit).Offset(offset).Find(&events).Error
	return events, err
}
