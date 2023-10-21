-- name: GetUser :one
SELECT * FROM users WHERE email = $1 LIMIT 1;

-- name: ListUsers :one
SELECT * FROM users ORDER BY email LIMIT $1 OFFSET $2;