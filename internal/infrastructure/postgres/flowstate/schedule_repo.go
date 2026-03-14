package flowstate

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/masterfabric/masterfabric_go_basic/internal/domain/flowstate/model"
)

// ScheduleRepo implements GeneratedScheduleRepository.
type ScheduleRepo struct {
	db *pgxpool.Pool
}

// NewScheduleRepo creates a new ScheduleRepo.
func NewScheduleRepo(db *pgxpool.Pool) *ScheduleRepo {
	return &ScheduleRepo{db: db}
}

func (r *ScheduleRepo) Upsert(ctx context.Context, s *model.GeneratedSchedule) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO flowstate_generated_schedules (id, user_id, week_identifier, schedule_data)
		VALUES ($1, $2, $3, $4::jsonb)
		ON CONFLICT (user_id, week_identifier) DO UPDATE SET schedule_data = EXCLUDED.schedule_data
	`, s.ID, s.UserID, s.WeekIdentifier, toJSONB(s.ScheduleData))
	return err
}

func (r *ScheduleRepo) FindByUserAndWeek(ctx context.Context, userID uuid.UUID, weekID string) (*model.GeneratedSchedule, error) {
	var s model.GeneratedSchedule
	var data []byte
	err := r.db.QueryRow(ctx, `
		SELECT id, user_id, week_identifier, schedule_data, created_at
		FROM flowstate_generated_schedules WHERE user_id = $1 AND week_identifier = $2
	`, userID, weekID).Scan(&s.ID, &s.UserID, &s.WeekIdentifier, &data, &s.CreatedAt)
	if err != nil {
		return nil, err
	}
	s.ScheduleData = parseJSONB(data)
	return &s, nil
}
