-- name: CreateUser :one
INSERT INTO users (
    id,
    email,
    password,
    name
)
VALUES (
    $1, $2, $3, $4
)
RETURNING *;


-- name: GetUserByID :one
SELECT *
FROM users
WHERE id = $1;


-- name: GetUserByEmail :one
SELECT *
FROM users
WHERE email = $1;