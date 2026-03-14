-- 007_create_flowstate_tables.up.sql
-- FlowState AI: Fixed events, flexible tasks, generated schedules.
-- user_id references existing users(id).

-- Fixed_Events (Statik Planlar)
CREATE TYPE flowstate_event_category AS ENUM ('Work', 'School', 'Meeting');

CREATE TABLE IF NOT EXISTS flowstate_fixed_events (
    id           UUID                     PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id      UUID                     NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title        TEXT                     NOT NULL,
    start_time   TIME                     NOT NULL,
    end_time     TIME                     NOT NULL,
    days_of_week INT[]                    NOT NULL DEFAULT '{}',
    category     flowstate_event_category NOT NULL,
    created_at   TIMESTAMPTZ              NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ              NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_flowstate_fixed_events_user_id ON flowstate_fixed_events(user_id);

-- Flexible_Tasks (Dinamik Görevler)
CREATE TYPE flowstate_task_priority AS ENUM ('High', 'Medium', 'Low');
CREATE TYPE flowstate_preferred_context AS ENUM ('Morning', 'Evening', 'WorkBreak', 'Weekend');

CREATE TABLE IF NOT EXISTS flowstate_flexible_tasks (
    id                 UUID                         PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id            UUID                         NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title              TEXT                         NOT NULL,
    duration_minutes   INT                          NOT NULL CHECK (duration_minutes > 0),
    frequency_per_week INT                          NOT NULL DEFAULT 1 CHECK (frequency_per_week >= 1 AND frequency_per_week <= 7),
    priority           flowstate_task_priority      NOT NULL DEFAULT 'Medium',
    preferred_context  flowstate_preferred_context  NOT NULL DEFAULT 'Morning',
    constraints        JSONB                        NOT NULL DEFAULT '{}',
    created_at         TIMESTAMPTZ                  NOT NULL DEFAULT NOW(),
    updated_at         TIMESTAMPTZ                  NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_flowstate_flexible_tasks_user_id ON flowstate_flexible_tasks(user_id);

-- Generated_Schedules (AI Çıktısı)
CREATE TABLE IF NOT EXISTS flowstate_generated_schedules (
    id              UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id         UUID        NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    week_identifier TEXT        NOT NULL,
    schedule_data   JSONB       NOT NULL,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(user_id, week_identifier)
);

CREATE INDEX IF NOT EXISTS idx_flowstate_generated_schedules_user_id ON flowstate_generated_schedules(user_id);
CREATE INDEX IF NOT EXISTS idx_flowstate_generated_schedules_week ON flowstate_generated_schedules(week_identifier);
