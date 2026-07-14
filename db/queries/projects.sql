-- name: CreateProject :one
INSERT INTO projects (
    id,
    user_id,
    brand_name,
    domain,
    language,
    region
)
VALUES (
    $1, $2, $3, $4, $5, $6
)
RETURNING *;

-- name: GetProjectByID :one
SELECT *
FROM projects
WHERE id = $1;

-- name: GetProjectByIDForUser :one
SELECT *
FROM projects
WHERE id = $1
  AND user_id = $2;

-- name: ListProjectsByUserID :many
SELECT *
FROM projects
WHERE user_id = $1
ORDER BY created_at DESC;

-- name: UpdateProject :one
UPDATE projects
SET
    brand_name = $2,
    domain = $3,
    language = $4,
    region = $5,
    updated_at = now()
WHERE id = $1
RETURNING *;

-- name: UpdateProjectForUser :one
UPDATE projects
SET
    brand_name = $3,
    domain = $4,
    language = $5,
    region = $6,
    updated_at = now()
WHERE id = $1
  AND user_id = $2
RETURNING *;