-- +goose Up
-- +goose StatementBegin

CREATE TABLE users (
  id               UUID         PRIMARY KEY DEFAULT uuid_generate_v4(),
  sub              TEXT         UNIQUE NOT NULL,
  name             TEXT         UNIQUE NOT NULL,
  email            TEXT         UNIQUE NOT NULL,
  created_at       TIMESTAMPTZ  NOT NULL DEFAULT now(),
  last_updated_at  TIMESTAMPTZ  NOT NULL DEFAULT now()
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS users;

-- +goose StatementEnd
