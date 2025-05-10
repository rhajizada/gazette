-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS vector;

CREATE TABLE IF NOT EXISTS item_embeddings (
  id            UUID        PRIMARY KEY DEFAULT uuid_generate_v4(),
  item_id       UUID        NOT NULL REFERENCES items(id) ON DELETE CASCADE,
  chunk_index   INT         NOT NULL DEFAULT 0,
  embedding     VECTOR(768) NOT NULL,
  created_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at    TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS idx_item_embeddings_vector
  ON item_embeddings
  USING ivfflat (embedding vector_l2_ops)
  WITH (lists = 100);

CREATE TABLE IF NOT EXISTS user_embedding_clusters (
  user_id      UUID        NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  cluster_id   INT         NOT NULL,
  centroid     VECTOR(768) NOT NULL,
  member_count INT         NOT NULL DEFAULT 0,
  created_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
  PRIMARY KEY (user_id, cluster_id)
);
CREATE INDEX IF NOT EXISTS idx_user_embedding_clusters_centroid
  ON user_embedding_clusters
  USING ivfflat (centroid vector_l2_ops)
  WITH (lists = 50);

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

DROP INDEX IF EXISTS idx_user_embedding_clusters_centroid;
DROP TABLE IF EXISTS user_embedding_clusters;
DROP INDEX IF EXISTS idx_item_embeddings_vector;
DROP TABLE IF EXISTS item_embeddings;
-- +goose StatementEnd
