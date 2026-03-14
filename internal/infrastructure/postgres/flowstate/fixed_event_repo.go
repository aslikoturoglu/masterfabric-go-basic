package flowstate

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/masterfabric/masterfabric_go_basic/internal/domain/flowstate/model"
)

// FixedEventRepo implements FixedEventRepository.
type FixedEventRepo struct {
	db *pgxpool.Pool
}

// NewFixedEventRepo creates a new FixedEventRepo.
func NewFixedEventRepo(db *pgxpool.Pool) *FixedEventRepo {
	return &FixedEventRepo{db: db}
}

func (r *FixedEventRepo) ListByUserID(ctx context.Context, userID uuid.UUID) ([]*model.FixedEvent, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, user_id, title, start_time::text, end_time::text, days_of_week, category::text, created_at, updated_at
		FROM flowstate_fixed_events WHERE user_id = $1 ORDER BY start_time
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []*model.FixedEvent
	for rows.Next() {
		var e model.FixedEvent
		err := rows.Scan(&e.ID, &e.UserID, &e.Title, &e.StartTime, &e.EndTime, &e.DaysOfWeek, &e.Category, &e.CreatedAt, &e.UpdatedAt)
		if err != nil {
			return nil, err
		}
		events = append(events, &e)
	}
	return events, rows.Err()
}

func (r *FixedEventRepo) Create(ctx context.Context, e *model.FixedEvent) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO flowstate_fixed_events (id, user_id, title, start_time, end_time, days_of_week, category)
		VALUES ($1, $2, $3, $4::time, $5::time, $6, $7::flowstate_event_category)
	`, e.ID, e.UserID, e.Title, e.StartTime, e.EndTime, e.DaysOfWeek, e.Category)
	return err
}

func (r *FixedEventRepo) Delete(ctx context.Context, id, userID uuid.UUID) error {
	result, err := r.db.Exec(ctx, `DELETE FROM flowstate_fixed_events WHERE id = $1 AND user_id = $2`, id, userID)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return errors.New("not found")
	}
	return nil
}

func toJSONB(m map[string]any) []byte {
	if m == nil {
		return []byte("{}")
	}
	b, _ := json.Marshal(m)
	return b
}

func parseJSONB(b []byte) map[string]any {
	if len(b) == 0 {
		return map[string]any{}
	}
	var m map[string]any
	_ = json.Unmarshal(b, &m)
	return m
}
