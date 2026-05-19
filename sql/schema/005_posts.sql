-- +goose Up
CREATE TABLE posts (
    id uuid PRIMARY KEY NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    title TEXT,
    url TEXT NOT NULL,
    description TEXT,
    published_at TIMESTAMP,
    feed_id uuid NOT NULL REFERENCES feeds(id) ON DELETE CASCADE,
    UNIQUE(url)
);
-- +goose Down
DROP TABLE posts;