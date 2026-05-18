-- +goose Up
CREATE TABLE feed_follows (
    id uuid NOT NULL PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    user_id uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    feed_id uuid NOT NULL REFERENCES feeds(id) ON DELETE CASCADE,
    UNIQUE(user_id, feed_id)
);

INSERT INTO feed_follows (id, created_at, updated_at, user_id, feed_id)
SELECT gen_random_uuid(), created_at, updated_at, user_id, id
FROM feeds;

ALTER TABLE feeds
DROP COLUMN user_id;

-- +goose Down
ALTER TABLE feeds
ADD COLUMN user_id uuid REFERENCES users(id) ON DELETE CASCADE;

WITH first_follows AS (
    SELECT DISTINCT ON (feed_id) feed_id, user_id
    FROM feed_follows
    ORDER BY feed_id, created_at ASC
)
UPDATE feeds
SET user_id = first_follows.user_id
FROM first_follows
WHERE first_follows.feed_id = feeds.id;

ALTER TABLE feeds
ALTER COLUMN user_id SET NOT NULL;

DROP TABLE feed_follows;