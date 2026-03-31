-- name: CreateUser :one
INSERT INTO users (name, email, password, role)
VALUES ($1, $2, $3, $4)
RETURNING id, name, email, role, created_at;

-- name: GetUserByEmail :one
SELECT id, name, email, password, role, created_at
FROM users
WHERE email = $1;

-- name: GetUserByID :one
SELECT id, name, email, role, created_at
FROM users
WHERE id = $1;

-- name: UpdateUserRole :execrows
UPDATE users
SET role = $2, updated_at = NOW()
WHERE id = $1;

-- REFRESH TOKENS

-- name: CreateRefreshToken :exec
INSERT INTO refresh_tokens (user_id, token_hash, expires_at)
VALUES ($1, $2, $3);

-- name: GetRefreshTokenByHash :one
SELECT id, user_id, expires_at
FROM refresh_tokens
WHERE token_hash = $1;

-- name: DeleteRefreshTokensByUserID :exec
DELETE FROM refresh_tokens
WHERE user_id = $1;
