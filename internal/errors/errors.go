// internal/errors/errors.go
package errors

import (
	"errors"
)

var (
	// ErrInvalidTimeRange is returned when the time range is invalid
	ErrInvalidTimeRange = errors.New("end time must be after start time")
	// ErrEventNotFound is returned when an event is not found
	ErrEventNotFound = errors.New("event not found")
	// ErrTimeSlotNotFound is returned when a time slot is not found
	ErrTimeSlotNotFound = errors.New("time slot not found")
	// ErrUserNotFound is returned when a user is not found
	ErrUserNotFound = errors.New("user not found")
	// ErrAvailabilityNotFound is returned when an availability record is not found
	ErrAvailabilityNotFound = errors.New("availability not found")
)
