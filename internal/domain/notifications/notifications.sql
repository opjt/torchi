-- name: CreateNotification :one
INSERT INTO notifications (
    service_id,
    body
) VALUES (
    $1, $2
)
RETURNING *;