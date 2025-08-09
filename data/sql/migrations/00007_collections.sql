-- +goose Up
-- +goose StatementBegin

CREATE TABLE collections (
  id            UUID         PRIMARY KEY DEFAULT uuid_generate_v4(),
  user_id       UUID         NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  name          TEXT         NOT NULL,
  created_at    TIMESTAMPTZ  NOT NULL DEFAULT now(),
  last_updated  TIMESTAMPTZ  NOT NULL DEFAULT now(),
  UNIQUE (user_id, name)          -- a user cannot have two collections with the same name
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS collections;

-- +goose StatementEnd
