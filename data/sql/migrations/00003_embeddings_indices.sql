-- +goose Up
-- +goose StatementBegin
CREATE INDEX IF NOT EXISTS idx_items_published_parsed
  ON items (published_parsed);

CREATE INDEX IF NOT EXISTS idx_item_embeddings_item_id
  ON item_embeddings (item_id);

CREATE INDEX IF NOT EXISTS idx_user_likes_item_user
  ON user_likes (item_id, user_id);

CREATE INDEX IF NOT EXISTS idx_user_embedding_clusters_user
  ON user_embedding_clusters (user_id);

CREATE INDEX IF NOT EXISTS idx_user_embedding_clusters_centroid_l2
  ON user_embedding_clusters
  USING ivfflat (centroid vector_l2_ops)
  WITH (lists = 100);

CREATE INDEX IF NOT EXISTS idx_user_feeds_feed
  ON user_feeds (feed_id);

CREATE INDEX IF NOT EXISTS idx_items_created_at
  ON items (created_at);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_items_created_at;
DROP INDEX IF EXISTS idx_user_feeds_feed;
DROP INDEX IF EXISTS idx_user_embedding_clusters_centroid_l2;
DROP INDEX IF EXISTS idx_user_embedding_clusters_user;
DROP INDEX IF EXISTS idx_user_likes_item_user;
DROP INDEX IF EXISTS idx_item_embeddings_item_id;
DROP INDEX IF EXISTS idx_items_published_parsed;
-- +goose StatementEnd

