package flowstate

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/masterfabric/masterfabric_go_basic/internal/domain/flowstate/model"
)

// FlexibleTaskRepo implements FlexibleTaskRepository.
type FlexibleTaskRepo struct {
	db *pgxpool.Pool
}

// NewFlexibleTaskRepo creates a new FlexibleTaskRepo.
func NewFlexibleTaskRepo(db *pgxpool.Pool) *FlexibleTaskRepo {
	return &FlexibleTaskRepo{db: db}
}

func (r *FlexibleTaskRepo) ListByUserID(ctx context.Context, userID uuid.UUID) ([]*model.FlexibleTask, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, user_id, title, duration_minutes, frequency_per_week, priority::text, preferred_context::text, constraints, created_at, updated_at
		FROM flowstate_flexible_tasks WHERE user_id = $1 ORDER BY priority DESC, created_at
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []*model.FlexibleTask
	for rows.Next() {
		var t model.FlexibleTask
		var constraints []byte
		err := rows.Scan(&t.ID, &t.UserID, &t.Title, &t.DurationMinutes, &t.FrequencyPerWeek, &t.Priority, &t.PreferredContext, &constraints, &t.CreatedAt, &t.UpdatedAt)
		if err != nil {
			return nil, err
		}
		t.Constraints = parseJSONB(constraints)
		tasks = append(tasks, &t)
	}
	return tasks, rows.Err()
}

func (r *FlexibleTaskRepo) Create(ctx context.Context, t *model.FlexibleTask) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO flowstate_flexible_tasks (id, user_id, title, duration_minutes, frequency_per_week, priority, preferred_context, constraints)
		VALUES ($1, $2, $3, $4, $5, $6::flowstate_task_priority, $7::flowstate_preferred_context, COALESCE($8::jsonb, '{}'))
	`, t.ID, t.UserID, t.Title, t.DurationMinutes, t.FrequencyPerWeek, t.Priority, t.PreferredContext, toJSONB(t.Constraints))
	return err
}

func (r *FlexibleTaskRepo) Delete(ctx context.Context, id, userID uuid.UUID) error {
	result, err := r.db.Exec(ctx, `DELETE FROM flowstate_flexible_tasks WHERE id = $1 AND user_id = $2`, id, userID)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return errors.New("not found")
	}
	return nil
}
