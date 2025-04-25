-- name: CreateCollection :one
INSERT INTO collections (user_id, name)
VALUES ($1, $2)
RETURNING id, user_id, name, created_at, last_updated;

-- name: CountCollectionsByUserID :one
SELECT COUNT(*) AS count
FROM collections
WHERE user_id = $1;

-- name: GetCollectionByID :one
SELECT id, user_id, name, created_at, last_updated
FROM collections
WHERE id = $1;

-- name: ListCollectionsByUser :many
SELECT id, user_id, name, created_at, last_updated
FROM collections
WHERE user_id = $1
ORDER BY created_at DESC
LIMIT  $2
OFFSET $3;

-- name: UpdateCollectionByID :one
UPDATE collections
SET name         = $2,
    last_updated = now()
WHERE id = $1
RETURNING id, user_id, name, created_at, last_updated;

-- name: DeleteCollectionByID :exec
DELETE FROM collections
WHERE id = $1;
