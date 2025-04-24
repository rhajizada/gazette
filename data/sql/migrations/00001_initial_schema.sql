-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE users (
  id               UUID         PRIMARY KEY DEFAULT uuid_generate_v4(),
  sub              TEXT         UNIQUE NOT NULL,
  name             TEXT         UNIQUE NOT NULL,
  email            TEXT         UNIQUE NOT NULL,
  created_at       TIMESTAMPTZ  NOT NULL DEFAULT now(),
  last_updated_at  TIMESTAMPTZ  NOT NULL DEFAULT now()
);

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

CREATE TABLE user_feeds (
  user_id        UUID         NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  feed_id        UUID         NOT NULL REFERENCES feeds(id) ON DELETE CASCADE,
  subscribed_at  TIMESTAMPTZ  NOT NULL DEFAULT now(),
  PRIMARY KEY (user_id, feed_id)
);

CREATE TABLE user_likes (
  user_id   UUID         NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  item_id   UUID         NOT NULL REFERENCES items(id) ON DELETE CASCADE,
  liked_at  TIMESTAMPTZ  NOT NULL DEFAULT now(),
  PRIMARY KEY (user_id, item_id)
);

CREATE TABLE collections (
  id            UUID         PRIMARY KEY DEFAULT uuid_generate_v4(),
  user_id       UUID         NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  name          TEXT         NOT NULL,
  created_at    TIMESTAMPTZ  NOT NULL DEFAULT now(),
  last_updated  TIMESTAMPTZ  NOT NULL DEFAULT now(),
  UNIQUE (user_id, name)          -- a user cannot have two collections with the same name
);

CREATE TABLE collection_items (
  collection_id UUID         NOT NULL REFERENCES collections(id) ON DELETE CASCADE,
  item_id       UUID         NOT NULL REFERENCES items(id) ON DELETE CASCADE,
  added_at      TIMESTAMPTZ  NOT NULL DEFAULT now(),
  PRIMARY KEY (collection_id, item_id)
);

CREATE OR REPLACE FUNCTION update_feed_last_updated()
  RETURNS TRIGGER AS $$
BEGIN
  IF (TG_OP = 'DELETE') THEN
    UPDATE feeds SET last_updated_at = now() WHERE id = OLD.feed_id;
  ELSE
    UPDATE feeds SET last_updated_at = now() WHERE id = NEW.feed_id;
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

DROP TABLE IF EXISTS collection_items;
DROP TABLE IF EXISTS collections;
DROP TABLE IF EXISTS user_likes;
DROP TABLE IF EXISTS user_feeds;
DROP TABLE IF EXISTS items;
DROP TABLE IF EXISTS feeds;
DROP TABLE IF EXISTS users;
-- +goose StatementEnd
