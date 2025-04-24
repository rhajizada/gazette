-- name: AddItemToCollection :one
INSERT INTO collection_items (collection_id, item_id)
VALUES ($1, $2)
RETURNING collection_id, item_id, added_at;

-- name: GetCollectionItem :one
SELECT collection_id, item_id, added_at
FROM collection_items
WHERE collection_id = $1
  AND item_id       = $2;

-- name: ListItemsInCollection :many
SELECT ci.collection_id, ci.item_id, ci.added_at, i.*
FROM collection_items ci
JOIN items i ON i.id = ci.item_id
WHERE ci.collection_id = $1
ORDER BY ci.added_at DESC
LIMIT  $2
OFFSET $3;

-- name: CountItemsInCollection :one
SELECT COUNT(*) AS count
FROM collection_items
WHERE collection_id = $1;

-- name: RemoveItemFromCollection :exec
DELETE FROM collection_items
WHERE collection_id = $1
  AND item_id       = $2;
