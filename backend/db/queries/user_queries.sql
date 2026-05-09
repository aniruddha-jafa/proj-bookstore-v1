-- name: GetUser :one
SELECT * FROM users WHERE id=$1;

-- name: EmailIsTaken :one
SELECT EXISTS (SELECT 1 FROM users WHERE email=$1) as user_email;

-- name: GetByEmail :one
SELECT * FROM users WHERE email=$1;

-- name: GetAllUsers :many
SELECT * FROM users;

-- name: CreateUser :one
INSERT INTO USERS (
    id, email, password, created_at, updated_at
) VALUES (
    $1, $2, $3, $4, $5
)
RETURNING *;

-- name: UpdateUser :one
UPDATE users SET
    email=$2,
    updated_at=$3
WHERE id=$1
RETURNING *;

-- name: DeleteUser :exec
DELETE FROM users WHERE id=$1;