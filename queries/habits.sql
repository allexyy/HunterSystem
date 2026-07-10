-- name: CreateHabit :one
INSERT INTO habits (user_id, title, description, difficulty, xp_reward, gold_reward)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetHabitByID :one
SELECT * FROM habits
WHERE id = $1
LIMIT 1;

-- name: ListActiveHabits :many
SELECT * FROM habits
WHERE user_id = $1 AND is_active = true
ORDER BY created_at;

-- name: UpdateHabit :one
UPDATE habits
SET title = $2, description = $3, difficulty = $4, xp_reward = $5, gold_reward = $6
WHERE id = $1
RETURNING *;

-- name: DeactivateHabit :one
UPDATE habits
SET is_active = false
WHERE id = $1
RETURNING *;

-- name: UpdateHabitStreak :one
UPDATE habits
SET
    current_streak      = $2,
    longest_streak      = $3,
    last_completed_date = $4
WHERE id = $1
RETURNING *;

-- name: CreateHabitStatReward :one
INSERT INTO habit_stat_rewards (habit_id, stat_code, amount)
VALUES ($1, $2, $3)
ON CONFLICT (habit_id, stat_code) DO UPDATE SET amount = EXCLUDED.amount
RETURNING *;

-- name: ListHabitStatRewards :many
SELECT * FROM habit_stat_rewards
WHERE habit_id = $1;
