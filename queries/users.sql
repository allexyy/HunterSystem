-- name: CreateUser :one
INSERT INTO users (telegram_id, username, timezone)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetUserByTelegramID :one
SELECT * FROM users
WHERE telegram_id = $1
LIMIT 1;

-- name: GetUserByID :one
SELECT * FROM users
WHERE id = $1
LIMIT 1;

-- name: UpdateUserXPGold :one
UPDATE users
SET xp = xp + $2, gold = gold + $3, updated_at = now()
WHERE id = $1
RETURNING *;

-- name: UpdateUserRank :one
UPDATE users
SET rank = $2, updated_at = now()
WHERE id = $1
RETURNING *;

-- name: UpdateUserTimezone :one
UPDATE users
SET timezone = $2, updated_at = now()
WHERE id = $1
RETURNING *;
