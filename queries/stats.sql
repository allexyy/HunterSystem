-- name: ListStats :many
SELECT * FROM stats
ORDER BY code;

-- name: GetUserStat :one
SELECT * FROM user_stats
WHERE user_id = $1 AND stat_code = $2;

-- name: ListUserStats :many
SELECT * FROM user_stats
WHERE user_id = $1
ORDER BY stat_code;

-- name: InitUserStats :exec
INSERT INTO user_stats (user_id, stat_code, value)
SELECT $1, code, 0 FROM stats
ON CONFLICT (user_id, stat_code) DO NOTHING;

-- name: UpsertUserStat :one
INSERT INTO user_stats (user_id, stat_code, value)
VALUES ($1, $2, $3)
ON CONFLICT (user_id, stat_code) DO UPDATE
    SET value = user_stats.value + EXCLUDED.value
RETURNING *;

-- name: CreateXPGoldTransaction :one
INSERT INTO xp_gold_transactions (user_id, quest_id, xp_delta, gold_delta, reason)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: ListUserTransactions :many
SELECT * FROM xp_gold_transactions
WHERE user_id = $1
ORDER BY created_at DESC
LIMIT $2;
