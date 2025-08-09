-- +goose Up
-- +goose StatementBegin

CREATE TABLE items (
  id               UUID         PRIMARY KEY DEFAULT uuid_generate_v4(),
  feed_id          UUID         NOT NULL REFERENCES feeds(id) ON DELETE CASCADE,
  title            TEXT,
  description      TEXT,
  content          TEXT,
  link             TEXT         UNIQUE NOT NULL,
  links            TEXT[],
  updated_parsed   TIMESTAMPTZ,
  published_parsed TIMESTAMPTZ,
  authors          JSONB,
  guid             TEXT,
  image            JSONB,
  categories       TEXT[],
  enclosures       JSONB,
  created_at       TIMESTAMPTZ  NOT NULL DEFAULT now(),
  updated_at       TIMESTAMPTZ  NOT NULL DEFAULT now()
);

CREATE OR REPLACE FUNCTION update_feed_last_updated() RETURNS trigger AS $$
BEGIN
  UPDATE feeds
  SET last_updated_at = now()
  WHERE id = NEW.feed_id;
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_feed_last_updated
AFTER INSERT OR UPDATE ON items
FOR EACH ROW EXECUTE PROCEDURE update_feed_last_updated();

CREATE INDEX IF NOT EXISTS idx_items_published_parsed
    ON items (published_parsed);

CREATE INDEX IF NOT EXISTS idx_items_created_at
    ON items (created_at);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
--
DROP TRIGGER IF EXISTS update_feed_last_updated ON items;
DROP FUNCTION IF EXISTS update_feed_last_updated;
DROP TABLE IF EXISTS items;
DROP INDEX IF EXISTS idx_items_published_parsed;
DROP INDEX IF EXISTS idx_items_created_at;

-- +goose StatementEnd
