// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: collection_embeddings.sql

package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/pgvector/pgvector-go"
)

const createCollectionEmbedding = `-- name: CreateCollectionEmbedding :one
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
  updated_at
`

type CreateCollectionEmbeddingParams struct {
	CollectionID uuid.UUID        `json:"collectionId"`
	Embedding    *pgvector.Vector `json:"embedding"`
}

func (q *Queries) CreateCollectionEmbedding(ctx context.Context, arg CreateCollectionEmbeddingParams) (CollectionEmbedding, error) {
	row := q.db.QueryRow(ctx, createCollectionEmbedding, arg.CollectionID, arg.Embedding)
	var i CollectionEmbedding
	err := row.Scan(
		&i.CollectionID,
		&i.Embedding,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const deleteCollectionEmbeddingByID = `-- name: DeleteCollectionEmbeddingByID :exec
DELETE FROM collection_embeddings
WHERE collection_id = $1
`

func (q *Queries) DeleteCollectionEmbeddingByID(ctx context.Context, collectionID uuid.UUID) error {
	_, err := q.db.Exec(ctx, deleteCollectionEmbeddingByID, collectionID)
	return err
}

const getCollectionEmbeddingByID = `-- name: GetCollectionEmbeddingByID :one
SELECT
  collection_id,
  embedding,
  created_at,
  updated_at
FROM collection_embeddings
WHERE collection_id = $1
`

func (q *Queries) GetCollectionEmbeddingByID(ctx context.Context, collectionID uuid.UUID) (CollectionEmbedding, error) {
	row := q.db.QueryRow(ctx, getCollectionEmbeddingByID, collectionID)
	var i CollectionEmbedding
	err := row.Scan(
		&i.CollectionID,
		&i.Embedding,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const updateCollectionEmbeddingByID = `-- name: UpdateCollectionEmbeddingByID :one
UPDATE collection_embeddings
SET
  embedding  = $2,
  updated_at = now()
WHERE collection_id = $1
RETURNING
  collection_id,
  embedding,
  created_at,
  updated_at
`

type UpdateCollectionEmbeddingByIDParams struct {
	CollectionID uuid.UUID        `json:"collectionId"`
	Embedding    *pgvector.Vector `json:"embedding"`
}

func (q *Queries) UpdateCollectionEmbeddingByID(ctx context.Context, arg UpdateCollectionEmbeddingByIDParams) (CollectionEmbedding, error) {
	row := q.db.QueryRow(ctx, updateCollectionEmbeddingByID, arg.CollectionID, arg.Embedding)
	var i CollectionEmbedding
	err := row.Scan(
		&i.CollectionID,
		&i.Embedding,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}
