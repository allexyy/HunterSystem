-- ============ ENUMS ============
CREATE TYPE quest_type AS ENUM ('daily', 'main', 'boss');
CREATE TYPE quest_status AS ENUM ('pending', 'completed', 'failed', 'expired');

-- ============ USERS ============
CREATE TABLE users (
    id           BIGSERIAL PRIMARY KEY,
    telegram_id  BIGINT UNIQUE NOT NULL,
    username     TEXT,
    timezone     TEXT NOT NULL DEFAULT 'UTC',
    rank         TEXT NOT NULL DEFAULT 'E',
    level        INT NOT NULL DEFAULT 1,
    xp           BIGINT NOT NULL DEFAULT 0,
    gold         BIGINT NOT NULL DEFAULT 0,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- ============ STATS (справочник) ============
CREATE TABLE stats (
    code  TEXT PRIMARY KEY,   -- 'INT', 'STR', 'END', 'WIS'
    name  TEXT NOT NULL,
    icon  TEXT
);

CREATE TABLE user_stats (
    user_id   BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    stat_code TEXT NOT NULL REFERENCES stats(code) ON DELETE CASCADE,
    value     INT NOT NULL DEFAULT 0,
    PRIMARY KEY (user_id, stat_code)
);

-- ============ HABITS ============
CREATE TABLE habits (
    id                  BIGSERIAL PRIMARY KEY,
    user_id             BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title               TEXT NOT NULL,
    description         TEXT,
    difficulty          TEXT NOT NULL DEFAULT 'normal', -- easy/normal/hard/epic
    xp_reward           INT NOT NULL,
    gold_reward         INT NOT NULL,
    is_active           BOOLEAN NOT NULL DEFAULT true,
    current_streak      INT NOT NULL DEFAULT 0,
    longest_streak      INT NOT NULL DEFAULT 0,
    last_completed_date DATE,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE habit_stat_rewards (
    habit_id  BIGINT NOT NULL REFERENCES habits(id) ON DELETE CASCADE,
    stat_code TEXT NOT NULL REFERENCES stats(code) ON DELETE CASCADE,
    amount    INT NOT NULL,
    PRIMARY KEY (habit_id, stat_code)
);

-- ============ QUESTS ============
CREATE TABLE quests (
    id           BIGSERIAL PRIMARY KEY,
    user_id      BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    habit_id     BIGINT REFERENCES habits(id) ON DELETE SET NULL,
    type         quest_type NOT NULL,
    title        TEXT NOT NULL,
    description  TEXT,
    xp_reward    INT NOT NULL,
    gold_reward  INT NOT NULL,
    status       quest_status NOT NULL DEFAULT 'pending',
    due_date     DATE,
    deadline_at  TIMESTAMPTZ,
    completed_at TIMESTAMPTZ,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE quest_stat_rewards (
    quest_id  BIGINT NOT NULL REFERENCES quests(id) ON DELETE CASCADE,
    stat_code TEXT NOT NULL REFERENCES stats(code) ON DELETE CASCADE,
    amount    INT NOT NULL,
    PRIMARY KEY (quest_id, stat_code)
);

-- ============ TRANSACTIONS LOG ============
CREATE TABLE xp_gold_transactions (
    id         BIGSERIAL PRIMARY KEY,
    user_id    BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    quest_id   BIGINT REFERENCES quests(id) ON DELETE SET NULL,
    xp_delta   INT NOT NULL DEFAULT 0,
    gold_delta INT NOT NULL DEFAULT 0,
    reason     TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- ============ INDEXES ============
CREATE INDEX idx_quests_user_status ON quests(user_id, status);
CREATE INDEX idx_quests_user_due_date ON quests(user_id, due_date) WHERE type = 'daily';
CREATE INDEX idx_habits_user_active ON habits(user_id) WHERE is_active = true;
CREATE INDEX idx_xp_gold_tx_user ON xp_gold_transactions(user_id);

-- защита от дублей daily-квеста на один день для одной привычки
CREATE UNIQUE INDEX uq_daily_quest_per_habit_per_day
    ON quests(habit_id, due_date)
    WHERE type = 'daily' AND habit_id IS NOT NULL;

-- ============ SEED: базовые статы ============
INSERT INTO stats (code, name, icon) VALUES
    ('INT', 'Intelligence', '🧠'),
    ('STR', 'Strength', '💪'),
    ('END', 'Endurance', '🛡️'),
    ('WIS', 'Wisdom', '📖');
