-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, name)
VALUES (
    $1,
    $2, 
    $3,
    $4
)
RETURNING *;

-- name: GetUser :one
SELECT * FROM users WHERE name = $1;


-- name: DropUsers :exec 
DROP TABLE users;

-- name: CreateUsers :exec
CREATE TABLE users (
    id UUID,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    name TEXT UNIQUE
);

-- name: GetUsers :many
SELECT name FROM users;