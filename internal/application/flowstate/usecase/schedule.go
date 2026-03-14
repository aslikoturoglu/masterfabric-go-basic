package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/masterfabric/masterfabric_go_basic/internal/application/flowstate/dto"
	"github.com/masterfabric/masterfabric_go_basic/internal/domain/flowstate/model"
	flowstateRepo "github.com/masterfabric/masterfabric_go_basic/internal/domain/flowstate/repository"
	"github.com/masterfabric/masterfabric_go_basic/internal/domain/iam/repository"
	"github.com/masterfabric/masterfabric_go_basic/internal/shared/middleware"
)

// ScheduleGenerator is the AI schedule generator interface.
type ScheduleGenerator interface {
	Generate(ctx context.Context, events []*model.FixedEvent, tasks []*model.FlexibleTask, preferences map[string]any) (map[string]any, error)
}

// GenerateScheduleUseCase generates a weekly schedule via AI.
type GenerateScheduleUseCase struct {
	eventRepo    flowstateRepo.FixedEventRepository
	taskRepo     flowstateRepo.FlexibleTaskRepository
	scheduleRepo flowstateRepo.GeneratedScheduleRepository
	userRepo     repository.UserRepository
	generator    ScheduleGenerator
}

// NewGenerateScheduleUseCase creates a new GenerateScheduleUseCase.
func NewGenerateScheduleUseCase(
	eventRepo flowstateRepo.FixedEventRepository,
	taskRepo flowstateRepo.FlexibleTaskRepository,
	scheduleRepo flowstateRepo.GeneratedScheduleRepository,
	userRepo repository.UserRepository,
	generator ScheduleGenerator,
) *GenerateScheduleUseCase {
	return &GenerateScheduleUseCase{
		eventRepo:    eventRepo,
		taskRepo:     taskRepo,
		scheduleRepo: scheduleRepo,
		userRepo:     userRepo,
		generator:    generator,
	}
}

// Execute generates and stores a weekly schedule.
func (uc *GenerateScheduleUseCase) Execute(ctx context.Context, req *dto.GenerateScheduleRequest) (*dto.ScheduleResponse, error) {
	userID := middleware.UserIDFromContext(ctx)
	if userID == uuid.Nil {
		return nil, fmt.Errorf("unauthorized")
	}

	events, err := uc.eventRepo.ListByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("list events: %w", err)
	}

	tasks, err := uc.taskRepo.ListByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("list tasks: %w", err)
	}

	user, err := uc.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("find user: %w", err)
	}
	preferences := map[string]any{}
	// User preferences could be stored in user.Bio or a JSON field - for now empty
	_ = user

	scheduleData, err := uc.generator.Generate(ctx, events, tasks, preferences)
	if err != nil {
		return nil, fmt.Errorf("generate schedule: %w", err)
	}

	weekID := req.WeekIdentifier
	if weekID == "" {
		year, week := time.Now().ISOWeek()
		weekID = fmt.Sprintf("%d-W%02d", year, week)
	}
	if w, ok := scheduleData["week_identifier"].(string); ok && w != "" {
		weekID = w
	}

	s := &model.GeneratedSchedule{
		ID:             uuid.New(),
		UserID:         userID,
		WeekIdentifier: weekID,
		ScheduleData:   scheduleData,
	}
	if err := uc.scheduleRepo.Upsert(ctx, s); err != nil {
		return nil, fmt.Errorf("upsert schedule: %w", err)
	}

	return &dto.ScheduleResponse{
		WeekIdentifier: weekID,
		ScheduleData:   scheduleData,
	}, nil
}

// GetScheduleUseCase retrieves a generated schedule.
type GetScheduleUseCase struct {
	repo flowstateRepo.GeneratedScheduleRepository
}

// NewGetScheduleUseCase creates a new GetScheduleUseCase.
func NewGetScheduleUseCase(repo flowstateRepo.GeneratedScheduleRepository) *GetScheduleUseCase {
	return &GetScheduleUseCase{repo: repo}
}

// Execute returns the schedule for the given week.
func (uc *GetScheduleUseCase) Execute(ctx context.Context, weekID string) (*dto.ScheduleResponse, error) {
	userID := middleware.UserIDFromContext(ctx)
	if userID == uuid.Nil {
		return nil, fmt.Errorf("unauthorized")
	}

	if weekID == "" {
		year, week := time.Now().ISOWeek()
		weekID = fmt.Sprintf("%d-W%02d", year, week)
	}

	s, err := uc.repo.FindByUserAndWeek(ctx, userID, weekID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil // not found — return null for nullable query
		}
		return nil, fmt.Errorf("find schedule: %w", err)
	}

	return &dto.ScheduleResponse{
		WeekIdentifier: s.WeekIdentifier,
		ScheduleData:   s.ScheduleData,
	}, nil
}
