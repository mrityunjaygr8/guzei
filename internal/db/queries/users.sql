-- name: UserRetrieve :one
SELECT email, "createdAt", "updatedAt", "ID", admin FROM users WHERE email = $1 LIMIT 1;

-- name: UsersList :many
WITH row_data AS (
    SELECT email, "createdAt", "updatedAt", "ID", admin FROM users ORDER BY email LIMIT $1 OFFSET $2
) SELECT
      *,
      (SELECT COUNT(*) FROM users) AS row_data
FROM row_data;

-- name: UserInsert :one
INSERT INTO users (email, password, "ID", admin) VALUES ($1, $2, $3, $4) RETURNING email, "createdAt", "updatedAt", "ID", admin;

-- name: UserUpdatePassword :exec
UPDATE users SET password = $2 WHERE "ID" = $1;

-- name: UserUpdateAdmin :exec
UPDATE users SET admin = $2 WHERE "ID" = $1;

-- name: UserDelete :exec
DELETE FROM users WHERE "ID" = $1;