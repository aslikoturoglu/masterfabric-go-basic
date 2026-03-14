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

// ListFixedEventsUseCase lists fixed events for the authenticated user.
type ListFixedEventsUseCase struct {
	repo flowstateRepo.FixedEventRepository
}

// NewListFixedEventsUseCase creates a new ListFixedEventsUseCase.
func NewListFixedEventsUseCase(repo flowstateRepo.FixedEventRepository) *ListFixedEventsUseCase {
	return &ListFixedEventsUseCase{repo: repo}
}

// Execute returns all fixed events for the user.
func (uc *ListFixedEventsUseCase) Execute(ctx context.Context) ([]*dto.FixedEventResponse, error) {
	userID := middleware.UserIDFromContext(ctx)
	if userID == uuid.Nil {
		return nil, fmt.Errorf("unauthorized")
	}

	events, err := uc.repo.ListByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("list fixed events: %w", err)
	}

	result := make([]*dto.FixedEventResponse, 0, len(events))
	for _, e := range events {
		result = append(result, fixedEventToDTO(e))
	}
	return result, nil
}

// CreateFixedEventUseCase creates a new fixed event.
type CreateFixedEventUseCase struct {
	repo flowstateRepo.FixedEventRepository
}

// NewCreateFixedEventUseCase creates a new CreateFixedEventUseCase.
func NewCreateFixedEventUseCase(repo flowstateRepo.FixedEventRepository) *CreateFixedEventUseCase {
	return &CreateFixedEventUseCase{repo: repo}
}

// Execute creates a fixed event.
func (uc *CreateFixedEventUseCase) Execute(ctx context.Context, req *dto.FixedEventRequest) (*dto.FixedEventResponse, error) {
	userID := middleware.UserIDFromContext(ctx)
	if userID == uuid.Nil {
		return nil, fmt.Errorf("unauthorized")
	}

	e := &model.FixedEvent{
		ID:         uuid.New(),
		UserID:     userID,
		Title:      req.Title,
		StartTime:  req.StartTime,
		EndTime:    req.EndTime,
		DaysOfWeek: req.DaysOfWeek,
		Category:   model.FixedEventCategory(req.Category),
	}
	if err := uc.repo.Create(ctx, e); err != nil {
		return nil, fmt.Errorf("create fixed event: %w", err)
	}
	return fixedEventToDTO(e), nil
}

// DeleteFixedEventUseCase deletes a fixed event.
type DeleteFixedEventUseCase struct {
	repo flowstateRepo.FixedEventRepository
}

// NewDeleteFixedEventUseCase creates a new DeleteFixedEventUseCase.
func NewDeleteFixedEventUseCase(repo flowstateRepo.FixedEventRepository) *DeleteFixedEventUseCase {
	return &DeleteFixedEventUseCase{repo: repo}
}

// Execute deletes a fixed event.
func (uc *DeleteFixedEventUseCase) Execute(ctx context.Context, id uuid.UUID) error {
	userID := middleware.UserIDFromContext(ctx)
	if userID == uuid.Nil {
		return fmt.Errorf("unauthorized")
	}
	return uc.repo.Delete(ctx, id, userID)
}

func fixedEventToDTO(e *model.FixedEvent) *dto.FixedEventResponse {
	return &dto.FixedEventResponse{
		ID:         e.ID.String(),
		UserID:     e.UserID.String(),
		Title:      e.Title,
		StartTime:  e.StartTime,
		EndTime:    e.EndTime,
		DaysOfWeek: e.DaysOfWeek,
		Category:   string(e.Category),
		CreatedAt:  e.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:  e.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}
