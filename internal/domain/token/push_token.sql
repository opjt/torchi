-- name: UpsertToken :one
INSERT INTO push_tokens (
    user_id,
    p256dh_key,
    auth_key,
    is_active
)
VALUES ($1, $2, $3, true)
ON CONFLICT (user_id, p256dh_key, auth_key)
DO UPDATE SET
    is_active = true
RETURNING user_id;