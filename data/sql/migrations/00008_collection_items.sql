-- +goose Up
-- +goose StatementBegin

CREATE TABLE collection_items (
  collection_id UUID         NOT NULL REFERENCES collections(id) ON DELETE CASCADE,
  item_id       UUID         NOT NULL REFERENCES items(id) ON DELETE CASCADE,
  added_at      TIMESTAMPTZ  NOT NULL DEFAULT now(),
  PRIMARY KEY (collection_id, item_id)
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS collection_items;

-- +goose StatementEnd
