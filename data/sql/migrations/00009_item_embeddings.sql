-- +goose Up
-- +goose StatementBegin

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

CREATE INDEX IF NOT EXISTS idx_item_embeddings_item_id
  ON item_embeddings (item_id);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP INDEX IF EXISTS idx_item_embeddings_item_id;
DROP INDEX IF EXISTS idx_item_embeddings_vector;
DROP TABLE IF EXISTS item_embeddings;

-- +goose StatementEnd
