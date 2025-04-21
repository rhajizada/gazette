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
