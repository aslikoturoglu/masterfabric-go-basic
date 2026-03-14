package flowstate

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/masterfabric/masterfabric_go_basic/internal/domain/flowstate/model"
)

// ScheduleGenerator produces optimized weekly schedules.
// Replace with OpenAI/Cursor API for production.
type ScheduleGenerator interface {
	Generate(ctx context.Context, events []*model.FixedEvent, tasks []*model.FlexibleTask, preferences map[string]any) (map[string]any, error)
}

// MockGenerator returns a placeholder schedule.
type MockGenerator struct{}

// NewMockGenerator creates a MockGenerator.
func NewMockGenerator() *MockGenerator {
	return &MockGenerator{}
}

func (m *MockGenerator) Generate(ctx context.Context, events []*model.FixedEvent, tasks []*model.FlexibleTask, preferences map[string]any) (map[string]any, error) {
	_ = ctx
	_ = preferences
	prompt := buildPrompt(events, tasks, preferences)
	_ = prompt

	days := map[string]any{
		"monday": []any{}, "tuesday": []any{}, "wednesday": []any{}, "thursday": []any{},
		"friday": []any{}, "saturday": []any{}, "sunday": []any{},
	}
	dayNames := []string{"monday", "tuesday", "wednesday", "thursday", "friday", "saturday", "sunday"}

	for _, e := range events {
		for _, d := range e.DaysOfWeek {
			if d >= 1 && d <= 7 {
				dayKey := dayNames[d-1]
				slots, _ := days[dayKey].([]any)
				slots = append(slots, map[string]any{
					"time": e.StartTime, "task": e.Title, "duration": 0, "type": "fixed", "category": string(e.Category),
				})
				days[dayKey] = slots
			}
		}
	}

	for _, t := range tasks {
		for i := 0; i < t.FrequencyPerWeek && i < 7; i++ {
			dayKey := dayNames[i]
			slots, _ := days[dayKey].([]any)
			slots = append(slots, map[string]any{
				"time": "18:00", "task": t.Title, "duration": t.DurationMinutes, "type": "flexible", "priority": string(t.Priority),
			})
			days[dayKey] = slots
		}
	}

	year, week := time.Now().ISOWeek()
	weekID := fmt.Sprintf("%d-W%02d", year, week)

	return map[string]any{
		"week_identifier": weekID,
		"generated_at":    time.Now().Format(time.RFC3339),
		"days":            days,
		"meta":            map[string]any{"events_count": len(events), "tasks_count": len(tasks)},
	}, nil
}

func buildPrompt(events []*model.FixedEvent, tasks []*model.FlexibleTask, preferences map[string]any) string {
	eventsJSON, _ := json.MarshalIndent(events, "", "  ")
	tasksJSON, _ := json.MarshalIndent(tasks, "", "  ")
	prefsJSON, _ := json.MarshalIndent(preferences, "", "  ")
	return fmt.Sprintf("Events: %s\nTasks: %s\nPreferences: %s", string(eventsJSON), string(tasksJSON), string(prefsJSON))
}
