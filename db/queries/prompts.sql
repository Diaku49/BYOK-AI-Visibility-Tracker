-- name: CreatePromptForUser :one
INSERT INTO prompts (
    id,
    project_id,
    text,
    active
)
SELECT
    sqlc.arg(id),
    sqlc.arg(project_id),
    sqlc.arg(text),
    sqlc.arg(active)
FROM projects p
WHERE p.id = sqlc.arg(project_id)
  AND p.user_id = sqlc.arg(user_id)
RETURNING *;

-- name: CreatePromptsForUser :execrows
INSERT INTO prompts (
    id,
    project_id,
    text,
    active
)
SELECT
    unnest(sqlc.arg(ids)::uuid[]),
    sqlc.arg(project_id),
    unnest(sqlc.arg(texts)::text[]),
    unnest(sqlc.arg(actives)::boolean[])
WHERE EXISTS (
    SELECT 1
    FROM projects p
    WHERE p.id = sqlc.arg(project_id)
      AND p.user_id = sqlc.arg(user_id)
);

-- name: GetPromptByIDForUser :one
SELECT pr.*
FROM prompts pr
JOIN projects p
  ON p.id = pr.project_id
WHERE pr.id = sqlc.arg(id)
  AND p.user_id = sqlc.arg(user_id);

-- name: ListPromptsByProjectForUser :many
SELECT pr.*
FROM prompts pr
JOIN projects p
  ON p.id = pr.project_id
WHERE pr.project_id = sqlc.arg(project_id)
  AND p.user_id = sqlc.arg(user_id)
ORDER BY pr.created_at DESC;

-- name: UpdatePromptForUser :one
UPDATE prompts pr
SET
    text = sqlc.arg(text),
    active = sqlc.arg(active),
    updated_at = now()
FROM projects p
WHERE pr.id = sqlc.arg(id)
  AND pr.project_id = p.id
  AND p.user_id = sqlc.arg(user_id)
RETURNING pr.*;
