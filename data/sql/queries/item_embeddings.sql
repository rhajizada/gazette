-- name: CreateItemEmbedding :one
INSERT INTO item_embeddings (
  item_id,
  embedding
)
VALUES (
  $1,
  $2
)
RETURNING
  item_id,
  embedding,
  created_at,
  updated_at;

-- name: GetItemEmbeddingByID :one
SELECT
  item_id,
  embedding,
  created_at,
  updated_at
FROM item_embeddings
WHERE item_id = $1;

-- name: UpdateItemEmbeddingByID :one
UPDATE item_embeddings
SET
  embedding  = $2,
  updated_at = now()
WHERE item_id = $1
RETURNING
  item_id,
  embedding,
  created_at,
  updated_at;

-- name: DeleteItemEmbeddingByID :exec
DELETE FROM item_embeddings
WHERE item_id = $1;
