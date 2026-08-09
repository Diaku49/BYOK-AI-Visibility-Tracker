-- name: CreateProviderKey :one
INSERT INTO provider_keys (
    id,
    user_id,
    engine_id,
    name,
    encrypted_key,
    key_nonce,
    active,
    monthly_run_limit
)
VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8
)
RETURNING *;

-- name: GetProviderKeyByID :one
SELECT *
FROM provider_keys
WHERE id = $1;

-- name: GetProviderKeyByIDForUser :one
SELECT *
FROM provider_keys
WHERE id = $1
  AND user_id = $2;

-- name: GetProviderKeyByIDForUserAndEngine :one
SELECT *
FROM provider_keys
WHERE id = $1
  AND user_id = $2
  AND engine_id = $3;

-- name: ListProviderKeysByUserID :many
SELECT *
FROM provider_keys
WHERE user_id = $1
ORDER BY created_at DESC;

-- name: ListProviderKeysByUserIDAndEngine :many
SELECT *
FROM provider_keys
WHERE user_id = $1
  AND engine_id = $2
ORDER BY created_at DESC;

-- name: UpdateProviderKeyMetadataForUser :one
UPDATE provider_keys
SET
    name = $3,
    active = $4,
    monthly_run_limit = $5,
    updated_at = now()
WHERE id = $1
  AND user_id = $2
RETURNING *;

-- name: DeleteProviderKeyForUser :exec
DELETE FROM provider_keys
WHERE id = $1
  AND user_id = $2;
