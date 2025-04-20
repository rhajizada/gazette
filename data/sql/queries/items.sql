-- name: CreateItem :one
INSERT INTO items
  (feed_id, title, description, content, link, links, updated_parsed, published_parsed,
   authors, guid, image, categories, enclosures)
VALUES
  ($1, $2, $3, $4, $5, $6, $7, $8,
   $9, $10, $11, $12, $13)
RETURNING
  id, feed_id, title, description, content, link, links, updated_parsed, published_parsed,
  authors, guid, image, categories, enclosures, created_at, updated_at;

-- name: GetItemByID :one
SELECT
  id, feed_id, title, description, content, link, links, updated_parsed, published_parsed,
  authors, guid, image, categories, enclosures, created_at, updated_at
FROM items
WHERE id = $1;

-- name: ListItemsByFeedID :many
SELECT
  id, feed_id, title, description, content, link, links, updated_parsed, published_parsed,
  authors, guid, image, categories, enclosures, created_at, updated_at
FROM items
WHERE feed_id = $1
ORDER BY created_at DESC
LIMIT  $2
OFFSET $3;

-- name: UpdateItemByID :one
UPDATE items
SET
  title            = $2,
  description      = $3,
  content          = $4,
  link             = $5,
  links            = $6,
  updated_parsed   = $7,
  published_parsed = $8,
  authors          = $9,
  guid             = $10,
  image            = $11,
  categories       = $12,
  enclosures       = $13,
  updated_at       = now()
WHERE id = $1
RETURNING
  id, feed_id, title, description, content, link, links, updated_parsed, published_parsed,
  authors, guid, image, categories, enclosures, created_at, updated_at;

-- name: DeleteItemByID :exec
DELETE FROM items WHERE id = $1;
