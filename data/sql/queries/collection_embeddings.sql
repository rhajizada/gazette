-- name: CreateCollectionEmbedding :one
INSERT INTO collection_embeddings (
  collection_id,
  embedding
)
VALUES (
  $1,
  $2
)
RETURNING
  collection_id,
  embedding,
  created_at,
  updated_at;

-- name: GetCollectionEmbeddingByID :one
SELECT
  collection_id,
  embedding,
  created_at,
  updated_at
FROM collection_embeddings
WHERE collection_id = $1;

-- name: UpdateCollectionEmbeddingByID :one
UPDATE collection_embeddings
SET
  embedding  = $2,
  updated_at = now()
WHERE collection_id = $1
RETURNING
  collection_id,
  embedding,
  created_at,
  updated_at;

-- name: DeleteCollectionEmbeddingByID :exec
DELETE FROM collection_embeddings
WHERE collection_id = $1;
