CREATE INDEX idx_tasks_team_status_updated ON tasks (team_id, status, updated_at);

CREATE INDEX idx_tasks_created_at ON tasks (created_at);
