package model

import (
	"time"

	"github.com/google/uuid"
)

// TaskPriority represents task priority.
type TaskPriority string

const (
	TaskPriorityHigh   TaskPriority = "High"
	TaskPriorityMedium TaskPriority = "Medium"
	TaskPriorityLow    TaskPriority = "Low"
)

// PreferredContext represents when the user prefers to do the task.
type PreferredContext string

const (
	PreferredContextMorning   PreferredContext = "Morning"
	PreferredContextEvening   PreferredContext = "Evening"
	PreferredContextWorkBreak PreferredContext = "WorkBreak"
	PreferredContextWeekend   PreferredContext = "Weekend"
)

// FlexibleTask is a dynamic task with frequency and constraints.
type FlexibleTask struct {
	ID               uuid.UUID
	UserID           uuid.UUID
	Title            string
	DurationMinutes  int
	FrequencyPerWeek int
	Priority         TaskPriority
	PreferredContext PreferredContext
	Constraints      map[string]any // e.g. {"prevent_back_to_back": "leg_day", "requires_home": true}
	CreatedAt        time.Time
	UpdatedAt        time.Time
}
