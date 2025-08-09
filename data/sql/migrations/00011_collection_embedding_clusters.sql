-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS collection_embedding_clusters (
  collection_id UUID        NOT NULL REFERENCES collections(id) ON DELETE CASCADE,
  cluster_id    INT         NOT NULL,
  centroid      VECTOR(768) NOT NULL,
  member_count  INT         NOT NULL DEFAULT 0,
  created_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
  PRIMARY KEY (collection_id, cluster_id)
);

CREATE INDEX IF NOT EXISTS idx_collection_embedding_clusters_centroid
  ON collection_embedding_clusters
  USING ivfflat (centroid vector_l2_ops)
  WITH (lists = 50);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP INDEX IF EXISTS idx_collection_embedding_clusters_centroid;
DROP TABLE IF EXISTS collection_embedding_clusters;

-- +goose StatementEnd
