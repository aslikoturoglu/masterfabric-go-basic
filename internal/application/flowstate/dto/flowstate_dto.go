package dto

// FixedEventRequest is input for creating a fixed event.
type FixedEventRequest struct {
	Title      string
	StartTime  string
	EndTime    string
	DaysOfWeek []int
	Category   string
}

// FixedEventResponse is returned to the client.
type FixedEventResponse struct {
	ID         string
	UserID     string
	Title      string
	StartTime  string
	EndTime    string
	DaysOfWeek []int
	Category   string
	CreatedAt  string
	UpdatedAt  string
}

// FlexibleTaskRequest is input for creating a flexible task.
type FlexibleTaskRequest struct {
	Title             string
	DurationMinutes   int
	FrequencyPerWeek  int
	Priority          string
	PreferredContext  string
	Constraints       map[string]any
}

// FlexibleTaskResponse is returned to the client.
type FlexibleTaskResponse struct {
	ID               string
	UserID           string
	Title            string
	DurationMinutes  int
	FrequencyPerWeek int
	Priority         string
	PreferredContext string
	Constraints      map[string]any
	CreatedAt        string
	UpdatedAt        string
}

// GenerateScheduleRequest is input for generating a schedule.
type GenerateScheduleRequest struct {
	WeekIdentifier string
}

// ScheduleResponse is the AI-generated schedule.
type ScheduleResponse struct {
	WeekIdentifier string
	ScheduleData   map[string]any
}
