-- name: CreateTableFeeds :exec
CREATE TABLE feeds (
    id UUID PRIMARY KEY, 
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    name TEXT,
    url TEXT UNIQUE,
    user_id UUID REFERENCES users(id) ON DELETE CASCADE 
);

-- name: DropTableFeeds :exec
DROP TABLE feeds;

-- name: CreateFeed :one
INSERT INTO feeds (id, created_at, updated_at, name, url ,user_id)
VALUES (
    $1,
    $2, 
    $3,
    $4,
    $5,
    $6
)
RETURNING *;

-- name: SelectAllFeeds :many 
SELECT * FROM feeds;

-- name: SelectFeedByURL :one
SELECT * FROM feeds WHERE url = $1;
