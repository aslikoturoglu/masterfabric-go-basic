package model

import (
	"time"

	"github.com/google/uuid"
)

// GeneratedSchedule holds the AI-generated weekly calendar.
type GeneratedSchedule struct {
	ID             uuid.UUID
	UserID         uuid.UUID
	WeekIdentifier string // e.g. "2024-W42"
	ScheduleData   map[string]any
	CreatedAt      time.Time
}
