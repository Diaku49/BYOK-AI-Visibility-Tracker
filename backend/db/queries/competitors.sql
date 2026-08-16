-- name: CreateCompetitorForUser :one
INSERT INTO competitors (
    id,
    project_id,
    name,
    domain
)
SELECT
    sqlc.arg(id),
    sqlc.arg(project_id),
    sqlc.arg(name),
    sqlc.narg(domain)
FROM projects p
WHERE p.id = sqlc.arg(project_id)
  AND p.user_id = sqlc.arg(user_id)
RETURNING *;

-- name: ListCompetitorsByProjectForUser :many
SELECT c.*
FROM competitors c
JOIN projects p
  ON p.id = c.project_id
WHERE c.project_id = sqlc.arg(project_id)
  AND p.user_id = sqlc.arg(user_id)
ORDER BY c.created_at DESC;

-- name: DeleteCompetitorForUser :execrows
DELETE FROM competitors c
USING projects p
WHERE c.id = sqlc.arg(id)
  AND c.project_id = p.id
  AND p.user_id = sqlc.arg(user_id);

-- name: ListCompetitorsByProject :many
SELECT c.*
FROM competitors c
WHERE c.project_id = sqlc.arg(project_id)
ORDER BY c.created_at DESC;
