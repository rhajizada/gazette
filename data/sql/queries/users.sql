-- name: CreateUser :one
INSERT INTO users (sub, name, email)
VALUES ($1, $2, $3)
RETURNING
  id, sub, name, email, created_at, last_updated_at;

-- name: GetUserByID :one
SELECT
  id, sub, name, email, created_at, last_updated_at
FROM users
WHERE id = $1;

-- name: GetUserBySub :one
SELECT
  id, sub, name, email, created_at, last_updated_at
FROM users
WHERE sub = $1;

-- name: GetUserByEmail :one
SELECT
  id, sub, name, email, created_at, last_updated_at
FROM users
WHERE email = $1;

-- name: ListUsers :many
SELECT
  id, sub, name, email, created_at, last_updated_at
FROM users
ORDER BY created_at DESC
LIMIT  $1
OFFSET $2;

-- name: CountUsers :one
SELECT COUNT(*) AS count
FROM users;

-- name: UpdateUserByID :one
UPDATE users
SET
  sub             = $2,
  name            = $3,
  email           = $4,
  last_updated_at = now()
WHERE id = $1
RETURNING
  id, sub, name, email, created_at, last_updated_at;

-- name: DeleteUserByID :exec
DELETE FROM users
WHERE id = $1;
