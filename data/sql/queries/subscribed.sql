-- name: CountSubscribedItemsByUser :one
SELECT
  COUNT(*)::bigint AS total
FROM items AS i
JOIN user_feeds AS fs
  ON fs.feed_id = i.feed_id
WHERE fs.user_id = $1
  AND i.published_parsed >= NOW() - INTERVAL '30 days';

-- name: ListSubscribedItemsByUser :many
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
  ul.item_id IS NOT NULL AS liked,
  ul.liked_at AS liked_at
FROM items AS i
JOIN user_feeds AS fs
  ON fs.feed_id = i.feed_id
LEFT JOIN user_likes AS ul
  ON ul.item_id = i.id AND ul.user_id = $1
WHERE fs.user_id = $1
  AND i.published_parsed >= NOW() - INTERVAL '30 days'
ORDER BY i.published_parsed DESC
LIMIT  $2
OFFSET $3;

