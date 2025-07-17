-- name: CreateFeedFollow :many
WITH insert_feed_follow AS (
    INSERT INTO feed_follows (id, created_at , updated_at, user_id, feed_id)
    VALUES (
        $1,
        $2,
        $3,
        $4,
        $5
    ) RETURNING *
) SELECT insert_feed_follow.*, feeds.name, users.name
FROM insert_feed_follow     
INNER JOIN users ON users.id = insert_feed_follow.user_id
INNER JOIN feeds ON feeds.id = insert_feed_follow.feed_id;
