-- +goose Up
CREATE TABLE feed_follows (
    id UUID PRIMARY KEY, 
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    feed_id UUID references feeds(id) ON DELETE CASCADE,
    UNIQUE (user_id, feed_id)
);

-- +goose Down 
DROP TABLE feed_follows;  