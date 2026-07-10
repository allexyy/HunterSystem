-- name: CreateQuest :one
INSERT INTO quests (user_id, habit_id, type, title, description, xp_reward, gold_reward, due_date, deadline_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
RETURNING *;

-- name: GetQuestByID :one
SELECT * FROM quests
WHERE id = $1
LIMIT 1;

-- name: ListPendingQuestsByUser :many
SELECT * FROM quests
WHERE user_id = $1 AND status = 'pending'
ORDER BY created_at;

-- name: ListDailyQuestsByDate :many
SELECT * FROM quests
WHERE user_id = $1 AND type = 'daily' AND due_date = $2
ORDER BY created_at;

-- name: GetDailyQuestForHabit :one
SELECT * FROM quests
WHERE habit_id = $1 AND due_date = $2 AND type = 'daily'
LIMIT 1;

-- name: CompleteQuest :one
UPDATE quests
SET status = 'completed', completed_at = now()
WHERE id = $1
RETURNING *;

-- name: ExpirePendingDailyQuests :many
UPDATE quests
SET status = 'expired'
WHERE user_id = $1 AND type = 'daily' AND due_date < $2 AND status = 'pending'
RETURNING *;

-- name: CreateQuestStatReward :one
INSERT INTO quest_stat_rewards (quest_id, stat_code, amount)
VALUES ($1, $2, $3)
RETURNING *;

-- name: ListQuestStatRewards :many
SELECT * FROM quest_stat_rewards
WHERE quest_id = $1;
