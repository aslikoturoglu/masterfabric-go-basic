-- 007_create_flowstate_tables.down.sql

DROP TABLE IF EXISTS flowstate_generated_schedules;
DROP TABLE IF EXISTS flowstate_flexible_tasks;
DROP TYPE IF EXISTS flowstate_preferred_context;
DROP TYPE IF EXISTS flowstate_task_priority;
DROP TABLE IF EXISTS flowstate_fixed_events;
DROP TYPE IF EXISTS flowstate_event_category;
