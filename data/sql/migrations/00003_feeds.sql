-- +goose Up
-- +goose StatementBegin

CREATE TABLE feeds (
  id               UUID         PRIMARY KEY DEFAULT uuid_generate_v4(),
  title            TEXT,
  description      TEXT,
  link             TEXT,
  feed_link        TEXT         UNIQUE NOT NULL,
  links            TEXT[],
  updated_parsed   TIMESTAMPTZ,
  published_parsed TIMESTAMPTZ,
  authors          JSONB,
  language         TEXT,
  image            JSONB,
  copyright        TEXT,
  generator        TEXT,
  categories       TEXT[],
  feed_type        TEXT,
  feed_version     TEXT,
  created_at       TIMESTAMPTZ  NOT NULL DEFAULT now(),
  last_updated_at  TIMESTAMPTZ  NOT NULL DEFAULT now()
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
--
DROP TABLE IF EXISTS feeds;

-- +goose StatementEnd
