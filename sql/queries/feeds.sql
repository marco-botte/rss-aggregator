-- name: CreateFeed :one
INSERT INTO feeds (id, created_at, updated_at, name, url, user_id)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6
)
RETURNING *;

-- name: DeleteOrphanedFeeds :exec
DELETE FROM feeds WHERE user_id IS NULL;

-- name: GetFeed :one
SELECT * FROM feeds WHERE url = $1;

-- name: GetFeeds :many
SELECT feeds.name, feeds.url, users.name AS username FROM feeds INNER JOIN users ON feeds.user_id=users.id;

-- name: MarkFeedFetched :exec
UPDATE feeds SET last_fetched_at = $2, updated_at = $2 WHERE feeds.id = $1;

-- name: GetNextFeedToFetch :one
SELECT * FROM feeds ORDER BY last_fetched_at NULLS FIRST LIMIT 1;