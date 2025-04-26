-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS vector;

CREATE TABLE item_embeddings (
  item_id         UUID PRIMARY KEY REFERENCES items(id) ON DELETE CASCADE,
  embedding       VECTOR(768),
  created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_item_embeddings_vector ON item_embeddings USING ivfflat (embedding vector_l2_ops);

CREATE TABLE user_embeddings (
  user_id         UUID PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
  embedding       VECTOR(768),
  created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_user_embeddings_vector ON user_embeddings USING ivfflat (embedding vector_l2_ops);

CREATE TABLE collection_embeddings (
  collection_id UUID PRIMARY KEY REFERENCES collections(id) ON DELETE CASCADE,
  embedding     VECTOR(768),
  created_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at    TIMESTAMPTZ NOT NULL DEFAULT now() 
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS user_embeddings;
DROP TABLE IF EXISTS item_embeddings;

DROP INDEX IF EXISTS idx_user_embeddings_vector;
DROP INDEX IF EXISTS idx_item_embeddings_vector;

DROP TABLE IF EXISTS collection_embeddings;
DROP INDEX IF EXISTS idx_collection_embeddings_vector;
-- +goose StatementEnd
