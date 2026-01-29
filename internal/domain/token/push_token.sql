-- name: UpsertToken :one
INSERT INTO push_tokens (
    user_id,
    p256dh_key,
    auth_key,
    endpoint
)
VALUES ($1, $2, $3, $4)
ON CONFLICT (p256dh_key, auth_key, endpoint)
DO UPDATE SET
    is_active = true,
    user_id = EXCLUDED.user_id
RETURNING user_id;

-- name: DeleteToken :exec
DELETE FROM push_tokens
WHERE endpoint = $1 
    AND p256dh_key = $2 
    AND auth_key = $3;

-- name: FindTokenByUserID :many
SELECT * FROM push_tokens
WHERE user_id = $1 AND is_active = true;

-- name: FindTokenByEndpoint :one
SELECT * FROM push_tokens
WHERE endpoint = $1;

-- name: DeactivatePushToken :exec
UPDATE push_tokens
SET is_active = false
WHERE endpoint = $1;