-- name: CountDistinctCategories :one
SELECT
  COUNT(DISTINCT category) AS unique_category_count
FROM items
CROSS JOIN LATERAL unnest(categories) AS category;

-- name: ListDistinctCategories :many
SELECT DISTINCT category
FROM items
CROSS JOIN LATERAL unnest(categories) AS category
ORDER BY category ASC
LIMIT  $1  -- max number of categories to return
OFFSET $2; -- number of categories to skip

-- name: CountItemsByCategoryForUser :one
SELECT
  COUNT(*) AS count
FROM items i
WHERE
  i.categories @> $1;

-- name: ListItemsByCategoryForUser :many
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
  (ul.user_id IS NOT NULL) AS liked,
  ul.liked_at
FROM items i
LEFT JOIN user_likes ul
  ON ul.item_id = i.id
  AND ul.user_id  = $2
WHERE
  i.categories @> $1
ORDER BY
  i.published_parsed DESC
LIMIT  $3
OFFSET $4;
