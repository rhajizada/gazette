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


-- name: GetLastItem :one
SELECT
  id, feed_id, title, description, content, link, links,
  updated_parsed, published_parsed, authors, guid, image,
  categories, enclosures, created_at, updated_at
FROM items
WHERE feed_id = $1
ORDER BY published_parsed DESC
LIMIT 1;

-- name: GetItemByID :one
SELECT
  id, feed_id, title, description, content, link, links, updated_parsed, published_parsed,
  authors, guid, image, categories, enclosures, created_at, updated_at
FROM items
WHERE id = $1;

-- name: CountItemsByFeedID :one
SELECT COUNT(*) AS count
FROM items
WHERE feed_id = $1;

-- name: ListItemsByFeedID :many
SELECT
  id, feed_id, title, description, content, link, links, updated_parsed, published_parsed,
  authors, guid, image, categories, enclosures, created_at, updated_at
FROM items
WHERE feed_id = $1
ORDER BY published_parsed DESC
LIMIT  $2
OFFSET $3;


-- name: ListItemsByFeedIDForUser :many
SELECT
  i.id,
  i.feed_id,
  i.title,
  i.description,
  i.content,
  i.link,
  i.links,
  i.updated_parsed,
  i.published_parsed,
  i.authors,
  i.guid,
  i.image,
  i.categories,
  i.enclosures,
  i.created_at,
  i.updated_at,
  (ul.user_id IS NOT NULL)        AS liked,
  ul.liked_at
FROM items i
LEFT JOIN user_likes ul
  ON ul.item_id = i.id
  AND ul.user_id = $2
WHERE i.feed_id = $1
ORDER BY i.published_parsed DESC
LIMIT  $3
OFFSET $4;

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
