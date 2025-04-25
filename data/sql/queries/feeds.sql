-- name: CreateFeed :one
INSERT INTO feeds
  (title, description, link, feed_link, links, updated_parsed, published_parsed,
   authors, language, image, copyright, generator,
   categories, feed_type, feed_version)
VALUES
  ($1, $2, $3, $4, $5, $6, $7,
   $8, $9, $10, $11, $12,
   $13, $14, $15)
RETURNING
  id, title, description, link, feed_link, links, updated_parsed, published_parsed,
  authors, language, image, copyright, generator,
  categories, feed_type, feed_version, created_at, last_updated_at;

-- name: CountFeeds :one
SELECT COUNT(*) AS count
FROM feeds;

-- name: GetFeedByID :one
SELECT
  id, title, description, link, feed_link, links, updated_parsed, published_parsed,
  authors, language, image, copyright, generator,
  categories, feed_type, feed_version, created_at, last_updated_at
FROM feeds
WHERE id = $1;

-- name: GetFeedByFeedLink :one
SELECT
  id, title, description, link, feed_link, links, updated_parsed, published_parsed,
  authors, language, image, copyright, generator,
  categories, feed_type, feed_version, created_at, last_updated_at
FROM feeds
WHERE feed_link = $1;

-- name: ListFeeds :many
SELECT
  id, title, description, link, feed_link, links, updated_parsed, published_parsed,
  authors, language, image, copyright, generator,
  categories, feed_type, feed_version, created_at, last_updated_at
FROM feeds
ORDER BY created_at DESC
LIMIT  $1
OFFSET $2;

-- name: UpdateFeedByID :one
UPDATE feeds
SET
  title           = $2,
  description     = $3,
  link            = $4,
  feed_link       = $5,
  links           = $6,
  updated_parsed  = $7,
  published_parsed= $8,
  authors         = $9,
  language        = $10,
  image           = $11,
  copyright       = $12,
  generator       = $13,
  categories      = $14,
  feed_type       = $15,
  feed_version    = $16,
  last_updated_at = now()
WHERE id = $1
RETURNING
  id, title, description, link, feed_link, links, updated_parsed, published_parsed,
  authors, language, image, copyright, generator,
  categories, feed_type, feed_version, created_at, last_updated_at;

-- name: DeleteFeedByID :exec
DELETE FROM feeds WHERE id = $1;

-- name: CountFeedsByUserID :one
SELECT COUNT(*) AS count
FROM user_feeds
WHERE user_id = $1;

-- name: ListFeedsByUserID :many
SELECT
  f.id, f.title, f.description, f.link, f.feed_link, f.links,
  f.updated_parsed, f.published_parsed,
  f.authors, f.language, f.image, f.copyright, f.generator,
  f.categories, f.feed_type, f.feed_version,
  f.created_at, f.last_updated_at,
  uf.subscribed_at
FROM feeds f
LEFT JOIN user_feeds uf
  ON uf.feed_id = f.id
  AND uf.user_id = $1
WHERE
  -- if subscribed_only = false, return all;
  -- if subscribed_only = true, only those where uf.user_id IS NOT NULL
  (NOT $2) OR (uf.user_id IS NOT NULL)
ORDER BY f.created_at DESC
LIMIT  $3
OFFSET $4;


-- name: GetUserFeedByID :one
SELECT
  f.id, f.title, f.description, f.link, f.feed_link, f.links,
  f.updated_parsed, f.published_parsed,
  f.authors, f.language, f.image, f.copyright, f.generator,
  f.categories, f.feed_type, f.feed_version,
  f.created_at, f.last_updated_at,
  uf.subscribed_at
FROM feeds f
JOIN user_feeds uf ON uf.feed_id = f.id
WHERE uf.user_id = $1
  AND f.id      = $2;

