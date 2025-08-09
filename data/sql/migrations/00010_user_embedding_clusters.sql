-- +goose Up
-- +goose StatementBegin

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

CREATE INDEX IF NOT EXISTS idx_user_embedding_clusters_centroid_l2
  ON user_embedding_clusters
  USING ivfflat (centroid vector_l2_ops)
  WITH (lists = 100);

CREATE INDEX IF NOT EXISTS idx_user_embedding_clusters_user
  ON user_embedding_clusters (user_id);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP INDEX IF EXISTS idx_user_embedding_clusters_user;
DROP INDEX IF EXISTS idx_user_embedding_clusters_centroid_l2;
DROP INDEX IF EXISTS idx_user_embedding_clusters_centroid;
DROP TABLE IF EXISTS user_embedding_clusters;

-- +goose StatementEnd
