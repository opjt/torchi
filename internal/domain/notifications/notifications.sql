-- name: CreateNotification :one
INSERT INTO notifications (
    endpoint_id,
    user_id,
    body
) VALUES (
    $1, $2, $3
)
RETURNING *;

-- name: UpdateStatusNotification :exec
UPDATE notifications
SET status = $2
WHERE id = $1;

-- name: FindNotificationByUserID :many
SELECT 
    n.*,
    e.name as endpoint_name
FROM notifications n
JOIN endpoints e ON n.endpoint_id = e.id
WHERE n.user_id = $1;

-- name: GetNotificationsWithCursor :many
SELECT 
    n.id,
    n.endpoint_id,
    n.user_id,
    n.body,
    n.status,
    n.is_read,
    n.read_at,
    n.created_at,
    e.name AS endpoint_name
FROM notifications n
JOIN endpoints e ON n.endpoint_id = e.id
WHERE n.user_id = $1 
  AND n.is_deleted = false
  AND (sqlc.narg('last_id')::uuid IS NULL OR n.id < sqlc.narg('last_id'))
ORDER BY n.id DESC
LIMIT $2;

-- name: MarkNotificationsAsReadBefore :exec
UPDATE notifications
SET is_read = true,
    read_at = now()
WHERE user_id = $1 
  AND is_read = false 
  AND id >= $2; 