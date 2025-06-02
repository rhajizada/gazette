-- name: CountSuggestedItemsByUser :one
SELECT COUNT(*)::bigint AS total
FROM (
  SELECT i.id
  FROM items AS i
  JOIN item_embeddings AS e ON e.item_id = i.id
  JOIN user_embedding_clusters AS c ON c.user_id = $1
  WHERE i.published_parsed >= NOW() - INTERVAL '90 days'
    AND NOT EXISTS (
      SELECT 1
      FROM user_likes ul
      WHERE ul.user_id = $1 AND ul.item_id = i.id
    )
  GROUP BY i.id
) AS suggestions;

-- name: ListSuggestedItemsByUser :many
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
  (1.0 - MIN(e.embedding <=> c.centroid))::double precision AS similarity,
  EXP(
    - (LN(2)::double precision / 7.0)
      * (EXTRACT(EPOCH FROM NOW() - i.published_parsed) / 86400.0)
  )::double precision AS freshness,
  (
    (1.0 - MIN(e.embedding <=> c.centroid))
      * EXP(
          - (LN(2)::double precision / 7.0)
            * (EXTRACT(EPOCH FROM NOW() - i.published_parsed) / 86400.0)
        )
  )::double precision AS score
FROM items AS i
JOIN item_embeddings AS e ON e.item_id = i.id
JOIN user_embedding_clusters AS c ON c.user_id = $1
WHERE i.published_parsed >= NOW() - INTERVAL '90 days'
  AND NOT EXISTS (
    SELECT 1
    FROM user_likes ul
    WHERE ul.user_id = $1 AND ul.item_id = i.id
  )
GROUP BY
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
  i.updated_at
ORDER BY score DESC
LIMIT  $2
OFFSET $3;

-- name: ListSuggestedItemsForUserCache :many
SELECT
  i.id,
  (1.0 - MIN(e.embedding <=> c.centroid))::double precision AS similarity,
  EXP(
    - (LN(2)::double precision / 7.0)
      * (EXTRACT(EPOCH FROM NOW() - i.published_parsed) / 86400.0)
  )::double precision AS freshness,
  (
    (1.0 - MIN(e.embedding <=> c.centroid))
      * EXP(
          - (LN(2)::double precision / 7.0)
            * (EXTRACT(EPOCH FROM NOW() - i.published_parsed) / 86400.0)
        )
  )::double precision AS score
FROM items AS i
JOIN item_embeddings AS e ON e.item_id = i.id
JOIN user_embedding_clusters AS c ON c.user_id = $1
WHERE i.published_parsed >= NOW() - INTERVAL '90 days'
  AND NOT EXISTS (
    SELECT 1
    FROM user_likes ul
    WHERE ul.user_id = $1 AND ul.item_id = i.id
  )
GROUP BY
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
  i.updated_at
ORDER BY score DESC;
