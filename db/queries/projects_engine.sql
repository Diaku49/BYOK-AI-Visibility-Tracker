-- name: CreateProjectEngineForUser :one
INSERT INTO project_engines (
    project_id,
    engine_id,
    provider_key_id
)
SELECT
    sqlc.arg(project_id),
    sqlc.arg(engine_id),
    sqlc.arg(provider_key_id)
FROM projects p
JOIN provider_keys pk
    ON pk.id = sqlc.arg(provider_key_id)
WHERE p.id = sqlc.arg(project_id)
  AND p.user_id = sqlc.arg(user_id)
  AND pk.user_id = sqlc.arg(user_id)
  AND pk.engine_id = sqlc.arg(engine_id)
RETURNING *;

-- name: ListActiveProjectEnginesForScan :many
SELECT pe.engine_id, pe.provider_key_id
FROM project_engines pe
JOIN projects p
  ON p.id = pe.project_id
JOIN provider_keys pk
  ON pk.id = pe.provider_key_id
WHERE pe.project_id = sqlc.arg(project_id)
  AND p.user_id = sqlc.arg(user_id)
  AND pk.active = true
ORDER BY pe.engine_id ASC;

-- name: UpdateProjectEngineForUser :one
UPDATE project_engines pe
SET
    provider_key_id = sqlc.arg(provider_key_id),
    updated_at = now()
FROM projects p
JOIN provider_keys pk
    ON pk.id = sqlc.arg(provider_key_id)
WHERE pe.project_id = sqlc.arg(project_id)
  AND pe.engine_id = sqlc.arg(engine_id)
  AND p.id = pe.project_id
  AND p.user_id = sqlc.arg(user_id)
  AND pk.user_id = sqlc.arg(user_id)
  AND pk.engine_id = pe.engine_id
RETURNING pe.*;

-- name: DeleteProjectEngineForUser :exec
DELETE FROM project_engines pe
USING projects p
WHERE pe.project_id = sqlc.arg(project_id)
  AND pe.engine_id = sqlc.arg(engine_id)
  AND p.id = pe.project_id
  AND p.user_id = sqlc.arg(user_id);
