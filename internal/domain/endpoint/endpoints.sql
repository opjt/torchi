-- name: CreateEndpoint :one
INSERT INTO endpoints (
    user_id,
    name,
    endpoint
) VALUES (
    $1, 
    $2, 
    $3
)
RETURNING *;

-- name: FindByUserID :many
SELECT * FROM endpoints
WHERE user_id = $1;