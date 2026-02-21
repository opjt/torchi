-- name: CreateNotification :one
INSERT INTO notifications (
    endpoint_id,
    endpoint_name,
    user_id,
    body
)
SELECT 
    e.id, 
    e.name,
    $1,    
    $2    
FROM endpoints e
WHERE e.id = $3 
RETURNING *;

-- name: CreateMuteNotification :one
INSERT INTO notifications (
    endpoint_id,
    endpoint_name,
    user_id,
    body,
    status,
    read_at
)
SELECT 
    e.id, 
    e.name,
    $1,    
    $2,
    $3, 
    now()
FROM endpoints e
WHERE e.id = $4
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
    n.read_at,
    n.created_at,
    n.endpoint_name
FROM notifications n
WHERE n.user_id = $1 
  AND n.is_deleted = false
  AND (sqlc.narg('endpoint_id')::uuid IS NULL OR n.endpoint_id = sqlc.narg('endpoint_id'))
  AND (sqlc.narg('last_id')::uuid IS NULL OR n.id < sqlc.narg('last_id'))
  AND (
    sqlc.narg('query')::text IS NULL 
    OR n.body ILIKE '%' || sqlc.narg('query') || '%'
)
ORDER BY n.id DESC
LIMIT $2;

-- name: MarkNotificationsAsReadBefore :exec
UPDATE notifications
SET read_at = now()
WHERE user_id = $1
  AND read_at is NULL
  AND (sqlc.narg('endpoint_id')::uuid IS NULL OR endpoint_id = sqlc.narg('endpoint_id'))
  AND id >= $2;

-- name: MarkDeleteNotificationByID :exec
UPDATE notifications
SET is_deleted = true
WHERE user_id = $1
  AND id = $2;  