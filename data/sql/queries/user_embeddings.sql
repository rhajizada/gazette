-- name: CreateUserEmbedding :one
INSERT INTO user_embeddings (
  user_id,
  embedding
)
VALUES (
  $1,
  $2
)
RETURNING
  user_id,
  embedding,
  created_at,
  updated_at;

-- name: GetUserEmbedding :one
SELECT
  user_id,
  embedding,
  created_at,
  updated_at
FROM user_embeddings
WHERE user_id = $1;

-- name: UpdateUserEmbedding :one
UPDATE user_embeddings
SET
  embedding  = $2,
  updated_at = now()
WHERE user_id = $1
RETURNING
  user_id,
  embedding,
  created_at,
  updated_at;

-- name: DeleteUserEmbedding :exec
DELETE FROM user_embeddings
WHERE user_id = $1;
