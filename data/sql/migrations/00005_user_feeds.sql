-- +goose Up
-- +goose StatementBegin

CREATE TABLE user_feeds (
  user_id        UUID         NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  feed_id        UUID         NOT NULL REFERENCES feeds(id) ON DELETE CASCADE,
  subscribed_at  TIMESTAMPTZ  NOT NULL DEFAULT now(),
  PRIMARY KEY (user_id, feed_id)
);

CREATE INDEX IF NOT EXISTS idx_user_feeds_feed
    ON user_feeds (feed_id);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP INDEX IF EXISTS idx_user_feeds_feed;
DROP TABLE IF EXISTS user_feeds;

-- +goose StatementEnd
