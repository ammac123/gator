-- name: CreateFeed :one
INSERT INTO feeds (id, created_at, updated_at, name, url
) VALUES (
    $1,
    $2,
    $3,
    $4,
    $5
)
RETURNING *;

-- name: GetAllFeeds :many
SELECT
    feeds.name,
    feeds.url,
    users.name as user
FROM feeds
INNER JOIN users
    ON feeds.user_id = users.id
;

-- name: GetFeedByURL :one
SELECT *
FROM feeds
WHERE url = $1;


-- name: CreateFeedFollow :one
WITH new_feed_follow AS (
    INSERT INTO feed_follows (id, created_at, updated_at, user_id, feed_id)
    VALUES (
        $1,
        $2,
        $3,
        $4,
        $5
    )
    RETURNING *
)
SELECT new_feed_follow.*, users.name as user, feeds.name as feed
FROM new_feed_follow
INNER JOIN users ON new_feed_follow.user_id = users.id
INNER JOIN feeds ON new_feed_follow.feed_id = feeds.id;

-- name: GetFeedFollowsForUser :many
SELECT feed_follows.*, users.name as user, feeds.name as feed
FROM feed_follows
INNER JOIN users ON feed_follows.user_id = users.id
INNER JOIN feeds ON feed_follows.feed_id = feeds.id
WHERE feed_follows.user_id = $1;

-- name: ResetAllFeeds :exec
DELETE FROM feeds
WHERE 1=1;