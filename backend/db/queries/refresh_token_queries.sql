-- name: GetById :one
SELECT * FROM refresh_tokens WHERE id=$1;

-- name: Create :one
INSERT INTO refresh_tokens (
    id, user_id, expires_at, created_at, updated_at
) VALUES (
    $1, $2, $3, $4, $5
)
RETURNING *;

-- name: Revoke :one
UPDATE refresh_tokens SET
    revoked_at=$2,
    updated_at=$3
WHERE id=$1
RETURNING *;