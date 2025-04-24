-- name: CreateUserFeedSubscription :one
INSERT INTO user_feeds (user_id, feed_id)
VALUES ($1, $2)
RETURNING user_id, feed_id, subscribed_at;

-- name: GetUserFeedSubscription :one
SELECT user_id, feed_id, subscribed_at
FROM user_feeds
WHERE user_id = $1
  AND feed_id = $2;

-- name: ListUserFeedSubscriptions :many
SELECT user_id, feed_id, subscribed_at
FROM user_feeds
WHERE user_id = $1
ORDER BY subscribed_at DESC
LIMIT  $2
OFFSET $3;

-- name: DeleteUserFeedSubscription :exec
DELETE FROM user_feeds
WHERE user_id = $1
  AND feed_id = $2;
