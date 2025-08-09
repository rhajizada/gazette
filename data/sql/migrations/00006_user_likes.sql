-- +goose Up
-- +goose StatementBegin

CREATE TABLE user_likes (
  user_id   UUID         NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  item_id   UUID         NOT NULL REFERENCES items(id) ON DELETE CASCADE,
  liked_at  TIMESTAMPTZ  NOT NULL DEFAULT now(),
  PRIMARY KEY (user_id, item_id)
);

CREATE INDEX IF NOT EXISTS idx_user_likes_item_user
    ON user_likes (item_id, user_id);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP INDEX IF EXISTS idx_user_likes_item_user;
DROP TABLE IF EXISTS user_likes;

-- +goose StatementEnd
