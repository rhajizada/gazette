-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE feeds (
  id               UUID        PRIMARY KEY DEFAULT uuid_generate_v4(),
  title            TEXT,
  description      TEXT,
  link             TEXT,
  feed_link        TEXT        UNIQUE NOT NULL,
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
  created_at       TIMESTAMPTZ NOT NULL DEFAULT now(),
  last_updated_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE items (
  id               UUID        PRIMARY KEY DEFAULT uuid_generate_v4(),
  feed_id          UUID        NOT NULL REFERENCES feeds(id) ON DELETE CASCADE,
  title            TEXT,
  description      TEXT,
  content          TEXT,
  link             TEXT        UNIQUE NOT NULL,
  links            TEXT[],
  updated_parsed   TIMESTAMPTZ,
  published_parsed TIMESTAMPTZ,
  authors          JSONB,
  guid             TEXT,
  image            JSONB,
  categories       TEXT[],
  enclosures       JSONB,
  created_at       TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at       TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE OR REPLACE FUNCTION update_feed_last_updated()
  RETURNS TRIGGER AS $$
BEGIN
  IF (TG_OP = 'DELETE') THEN
    UPDATE feeds
      SET last_updated_at = now()
    WHERE id = OLD.feed_id;
  ELSE
    UPDATE feeds
      SET last_updated_at = now()
    WHERE id = NEW.feed_id;
  END IF;
  RETURN NULL;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_items_touch_feed
  AFTER INSERT OR UPDATE OR DELETE ON items
  FOR EACH ROW
  EXECUTE FUNCTION update_feed_last_updated();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TRIGGER IF EXISTS trg_items_touch_feed ON items;
DROP FUNCTION IF EXISTS update_feed_last_updated();
DROP TABLE IF EXISTS items;
DROP TABLE IF EXISTS feeds;
-- +goose StatementEnd
