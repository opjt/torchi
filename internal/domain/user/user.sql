-- name: UpsertUserByProvider :one
INSERT INTO users (provider, provider_id, email, updated_at)
VALUES ($1, $2, $3, now())
ON CONFLICT (provider, provider_id)
DO UPDATE SET
    email = COALESCE(EXCLUDED.email, users.email),
    updated_at = EXCLUDED.updated_at
RETURNING *;

-- name: FindUserById :one
SELECT id, email, provider, provider_id, guest, terms_agreed, created_at, updated_at
FROM users
WHERE id = $1;

-- name: UpdateUserTermsAgreed :exec
UPDATE users
SET terms_agreed = true
WHERE id = $1;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;

-- name: InsertGuestUser :one
INSERT INTO users (guest, terms_agreed)
VALUES (true, true)
RETURNING *;

