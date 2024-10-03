-- name: CreateFeed :one
INSERT INTO feeds(id, created_at, updated_at, name, url, user_id)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6
)
RETURNING *;

-- name: GetFeeds :many
SELECT * FROM feeds;

-- name: GetFeedFromUrl :one
SELECT * FROM feeds WHERE feeds.url = $1;

-- name: MarkFeedFetched :one
-- Requires migration 004: Add last_fetched_at column to feeds table
UPDATE feeds
SET last_fetched_at = NOW(), updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: GetNextFeedToFetch :one
-- Requires migration 004: Add last_fetched_at column to feeds table
SELECT *
FROM feeds
ORDER BY last_fetched_at NULLS FIRST
LIMIT 1;