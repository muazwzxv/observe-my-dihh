-- name: GetUser :one
SELECT * FROM users WHERE id = ?;

-- name: ListUsers :many
SELECT * FROM users ORDER BY created_at DESC;

-- name: CreateUser :execresult
INSERT INTO users (name, email, created_at, updated_at)
VALUES (?, ?, NOW(), NOW());

-- name: UpdateUser :exec
UPDATE users
SET name = ?, email = ?, updated_at = NOW()
WHERE id = ?;

-- name: DeleteUser :exec
DELETE FROM users WHERE id = ?;
