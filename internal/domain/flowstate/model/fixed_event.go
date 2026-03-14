package model

import (
	"time"

	"github.com/google/uuid"
)

// FixedEventCategory represents the type of fixed event.
type FixedEventCategory string

const (
	FixedEventCategoryWork    FixedEventCategory = "Work"
	FixedEventCategorySchool  FixedEventCategory = "School"
	FixedEventCategoryMeeting FixedEventCategory = "Meeting"
)

// FixedEvent is a static recurring event (e.g. work 09:00-18:00, school on Tuesday).
type FixedEvent struct {
	ID         uuid.UUID
	UserID     uuid.UUID
	Title      string
	StartTime  string // "09:00"
	EndTime    string // "18:00"
	DaysOfWeek []int  // 1=Monday .. 7=Sunday
	Category   FixedEventCategory
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
