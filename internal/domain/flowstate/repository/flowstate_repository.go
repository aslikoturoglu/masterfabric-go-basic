package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/masterfabric/masterfabric_go_basic/internal/domain/flowstate/model"
)

// FixedEventRepository defines persistence for fixed events.
type FixedEventRepository interface {
	ListByUserID(ctx context.Context, userID uuid.UUID) ([]*model.FixedEvent, error)
	Create(ctx context.Context, e *model.FixedEvent) error
	Delete(ctx context.Context, id, userID uuid.UUID) error
}

// FlexibleTaskRepository defines persistence for flexible tasks.
type FlexibleTaskRepository interface {
	ListByUserID(ctx context.Context, userID uuid.UUID) ([]*model.FlexibleTask, error)
	Create(ctx context.Context, t *model.FlexibleTask) error
	Delete(ctx context.Context, id, userID uuid.UUID) error
}

// GeneratedScheduleRepository defines persistence for AI-generated schedules.
type GeneratedScheduleRepository interface {
	Upsert(ctx context.Context, s *model.GeneratedSchedule) error
	FindByUserAndWeek(ctx context.Context, userID uuid.UUID, weekID string) (*model.GeneratedSchedule, error)
}
