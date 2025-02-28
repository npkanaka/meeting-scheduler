package timeutil

import (
	"time"
)

// TimeRange represents a time range with a start and end time
type TimeRange struct {
	Start time.Time
	End   time.Time
}

// Overlaps checks if this time range overlaps with another
func (tr TimeRange) Overlaps(other TimeRange) bool {
	return (tr.Start.Before(other.End) || tr.Start.Equal(other.End)) &&
		(tr.End.After(other.Start) || tr.End.Equal(other.Start))
}

// Contains checks if this time range contains another
func (tr TimeRange) Contains(other TimeRange) bool {
	return (tr.Start.Before(other.Start) || tr.Start.Equal(other.Start)) &&
		(tr.End.After(other.End) || tr.End.Equal(other.End))
}

// Duration returns the duration of the time range
func (tr TimeRange) Duration() time.Duration {
	return tr.End.Sub(tr.Start)
}

// FindOverlap returns the overlap between two time ranges
func FindOverlap(a, b TimeRange) (TimeRange, bool) {
	if !a.Overlaps(b) {
		return TimeRange{}, false
	}

	start := a.Start
	if b.Start.After(a.Start) {
		start = b.Start
	}

	end := a.End
	if b.End.Before(a.End) {
		end = b.End
	}

	return TimeRange{
		Start: start,
		End:   end,
	}, true
}

// FindCommonAvailability finds common availability time ranges between multiple availability ranges
func FindCommonAvailability(availabilities []TimeRange) []TimeRange {
	if len(availabilities) == 0 {
		return []TimeRange{}
	}

	if len(availabilities) == 1 {
		return availabilities
	}

	// Start with the first availability
	commonRanges := []TimeRange{availabilities[0]}

	// Intersect with each subsequent availability
	for i := 1; i < len(availabilities); i++ {
		var newCommonRanges []TimeRange

		for _, common := range commonRanges {
			if overlap, exists := FindOverlap(common, availabilities[i]); exists {
				newCommonRanges = append(newCommonRanges, overlap)
			}
		}

		commonRanges = newCommonRanges
	}

	return commonRanges
}

// FormatTime formats a time in RFC3339 format
func FormatTime(t time.Time) string {
	return t.Format(time.RFC3339)
}

// ParseTime parses a time string in RFC3339 format
func ParseTime(s string) (time.Time, error) {
	return time.Parse(time.RFC3339, s)
}

// ConvertTimeZone converts a time from one time zone to another
func ConvertTimeZone(t time.Time, fromTZ, toTZ string) (time.Time, error) {
	loc, err := time.LoadLocation(fromTZ)
	if err != nil {
		return time.Time{}, err
	}

	tInFromTZ := t.In(loc)

	toLocation, err := time.LoadLocation(toTZ)
	if err != nil {
		return time.Time{}, err
	}

	return tInFromTZ.In(toLocation), nil
}

// FormatTimeWithTZ formats a time with time zone information
func FormatTimeWithTZ(t time.Time, timezone string) (string, error) {
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		return "", err
	}

	return t.In(loc).Format(time.RFC3339), nil
}

// ParseTimeWithTZ parses a time string with specified time zone
func ParseTimeWithTZ(s string, timezone string) (time.Time, error) {
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return time.Time{}, err
	}

	loc, err := time.LoadLocation(timezone)
	if err != nil {
		return time.Time{}, err
	}

	return t.In(loc), nil
}

// GetCommonTimeSlots finds time slots that overlap between different time zones
func GetCommonTimeSlots(slots []TimeRange, duration time.Duration) []TimeRange {
	if len(slots) == 0 {
		return []TimeRange{}
	}

	// Find overlapping time ranges
	commonRanges := FindCommonAvailability(slots)

	// Filter out ranges that are shorter than the required duration
	var validRanges []TimeRange
	for _, r := range commonRanges {
		if r.Duration() >= duration {
			validRanges = append(validRanges, r)
		}
	}

	return validRanges
}

// GetTimeZoneAbbreviation returns common abbreviation for a time zone
func GetTimeZoneAbbreviation(timezone string, t time.Time) (string, error) {
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		return "", err
	}

	name, _ := t.In(loc).Zone()
	return name, nil
}
