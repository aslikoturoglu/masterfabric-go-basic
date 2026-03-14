package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/masterfabric/masterfabric_go_basic/internal/application/flowstate/dto"
	"github.com/masterfabric/masterfabric_go_basic/internal/domain/flowstate/model"
	flowstateRepo "github.com/masterfabric/masterfabric_go_basic/internal/domain/flowstate/repository"
	"github.com/masterfabric/masterfabric_go_basic/internal/shared/middleware"
)

// ListFlexibleTasksUseCase lists flexible tasks for the authenticated user.
type ListFlexibleTasksUseCase struct {
	repo flowstateRepo.FlexibleTaskRepository
}

// NewListFlexibleTasksUseCase creates a new ListFlexibleTasksUseCase.
func NewListFlexibleTasksUseCase(repo flowstateRepo.FlexibleTaskRepository) *ListFlexibleTasksUseCase {
	return &ListFlexibleTasksUseCase{repo: repo}
}

// Execute returns all flexible tasks for the user.
func (uc *ListFlexibleTasksUseCase) Execute(ctx context.Context) ([]*dto.FlexibleTaskResponse, error) {
	userID := middleware.UserIDFromContext(ctx)
	if userID == uuid.Nil {
		return nil, fmt.Errorf("unauthorized")
	}

	tasks, err := uc.repo.ListByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("list flexible tasks: %w", err)
	}

	result := make([]*dto.FlexibleTaskResponse, 0, len(tasks))
	for _, t := range tasks {
		result = append(result, flexibleTaskToDTO(t))
	}
	return result, nil
}

// CreateFlexibleTaskUseCase creates a new flexible task.
type CreateFlexibleTaskUseCase struct {
	repo flowstateRepo.FlexibleTaskRepository
}

// NewCreateFlexibleTaskUseCase creates a new CreateFlexibleTaskUseCase.
func NewCreateFlexibleTaskUseCase(repo flowstateRepo.FlexibleTaskRepository) *CreateFlexibleTaskUseCase {
	return &CreateFlexibleTaskUseCase{repo: repo}
}

// Execute creates a flexible task.
func (uc *CreateFlexibleTaskUseCase) Execute(ctx context.Context, req *dto.FlexibleTaskRequest) (*dto.FlexibleTaskResponse, error) {
	userID := middleware.UserIDFromContext(ctx)
	if userID == uuid.Nil {
		return nil, fmt.Errorf("unauthorized")
	}

	priority := model.TaskPriority(req.Priority)
	if priority == "" {
		priority = model.TaskPriorityMedium
	}
	ctxPref := model.PreferredContext(req.PreferredContext)
	if ctxPref == "" {
		ctxPref = model.PreferredContextMorning
	}

	t := &model.FlexibleTask{
		ID:               uuid.New(),
		UserID:           userID,
		Title:            req.Title,
		DurationMinutes:  req.DurationMinutes,
		FrequencyPerWeek: req.FrequencyPerWeek,
		Priority:         priority,
		PreferredContext: ctxPref,
		Constraints:      req.Constraints,
	}
	if t.Constraints == nil {
		t.Constraints = map[string]any{}
	}

	if err := uc.repo.Create(ctx, t); err != nil {
		return nil, fmt.Errorf("create flexible task: %w", err)
	}
	return flexibleTaskToDTO(t), nil
}

// DeleteFlexibleTaskUseCase deletes a flexible task.
type DeleteFlexibleTaskUseCase struct {
	repo flowstateRepo.FlexibleTaskRepository
}

// NewDeleteFlexibleTaskUseCase creates a new DeleteFlexibleTaskUseCase.
func NewDeleteFlexibleTaskUseCase(repo flowstateRepo.FlexibleTaskRepository) *DeleteFlexibleTaskUseCase {
	return &DeleteFlexibleTaskUseCase{repo: repo}
}

// Execute deletes a flexible task.
func (uc *DeleteFlexibleTaskUseCase) Execute(ctx context.Context, id uuid.UUID) error {
	userID := middleware.UserIDFromContext(ctx)
	if userID == uuid.Nil {
		return fmt.Errorf("unauthorized")
	}
	return uc.repo.Delete(ctx, id, userID)
}

func flexibleTaskToDTO(t *model.FlexibleTask) *dto.FlexibleTaskResponse {
	return &dto.FlexibleTaskResponse{
		ID:               t.ID.String(),
		UserID:           t.UserID.String(),
		Title:            t.Title,
		DurationMinutes:  t.DurationMinutes,
		FrequencyPerWeek: t.FrequencyPerWeek,
		Priority:         string(t.Priority),
		PreferredContext: string(t.PreferredContext),
		Constraints:      t.Constraints,
		CreatedAt:        t.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:        t.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}
