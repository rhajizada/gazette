-- name: CreateUserLike :one
INSERT INTO user_likes (user_id, item_id)
VALUES ($1, $2)
RETURNING user_id, item_id, liked_at;

-- name: GetUserLike :one
SELECT user_id, item_id, liked_at
FROM user_likes
WHERE user_id = $1
  AND item_id = $2;

-- name: ListUserLikesByUser :many
SELECT ul.user_id, ul.item_id, ul.liked_at, i.*
FROM user_likes ul
JOIN items i    ON i.id = ul.item_id
WHERE ul.user_id = $1
ORDER BY ul.liked_at DESC
LIMIT  $2
OFFSET $3;

-- name: ListUserLikesByItem :many
SELECT user_id, item_id, liked_at
FROM user_likes
WHERE item_id = $1
ORDER BY liked_at DESC
LIMIT  $2
OFFSET $3;

-- name: DeleteUserLike :exec
DELETE FROM user_likes
WHERE user_id = $1
  AND item_id = $2;
