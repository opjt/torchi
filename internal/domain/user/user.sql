-- name: UpsertUserByEmail :one
INSERT INTO users (email, updated_at)
VALUES ($1, now())
ON CONFLICT (email) 
DO UPDATE SET 
    updated_at = EXCLUDED.updated_at
RETURNING *;


-- name: FindUserById :one
SELECT id, email, created_at, updated_at, terms_agreed
FROM users
WHERE id = $1;

-- name: UpdateUserTermsAgreed :exec
UPDATE users
SET terms_agreed = true
WHERE id = $1;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;
