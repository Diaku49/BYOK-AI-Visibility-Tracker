-- +goose Up
CREATE INDEX IF NOT EXISTS idx_projects_user_id_created_at
ON projects(user_id, created_at DESC);

CREATE INDEX IF NOT EXISTS idx_provider_keys_user_engine
ON provider_keys(user_id, engine_id);

CREATE INDEX IF NOT EXISTS idx_prompts_project_active
ON prompts(project_id, active);

CREATE INDEX IF NOT EXISTS idx_scans_project_id_created_at
ON scans(project_id, created_at DESC);

CREATE INDEX IF NOT EXISTS idx_scan_runs_scan_status
ON scan_runs(scan_id, status);


-- +goose Down
DROP INDEX IF EXISTS idx_scan_runs_scan_status;
DROP INDEX IF EXISTS idx_scans_project_id_created_at;
DROP INDEX IF EXISTS idx_prompts_project_active;
DROP INDEX IF EXISTS idx_provider_keys_user_engine;
DROP INDEX IF EXISTS idx_projects_user_id_created_at;