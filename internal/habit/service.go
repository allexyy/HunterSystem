package habit

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/yourname/hunter-system/internal/db"
	"github.com/yourname/hunter-system/internal/reward"
	"strings"
)

type Service struct {
	q  db.Querier
	tx TxRunner
}

func NewService(q db.Querier, tx TxRunner) *Service {
	return &Service{q: q, tx: tx}
}

func (s *Service) CreateHabit(ctx context.Context, userId int64, data HabitData) (db.Habit, error) {
	var habit db.Habit
	rewards := reward.GetRewardByDifficult(data.Difficult)

	err := s.tx.Transaction(ctx, func(q db.Querier) error {
		var err error
		habit, err = q.CreateHabit(ctx, db.CreateHabitParams{
			UserID:      userId,
			Title:       strings.TrimSpace(strings.TrimPrefix(data.Title, "/newhabit")),
			Description: pgtype.Text{String: data.Description, Valid: true},
			Difficulty:  data.Difficult,
			XpReward:    rewards.XP,
			GoldReward:  rewards.Gold,
		})
		if err != nil {
			return fmt.Errorf("create habit: %w", err)
		}
		for _, stat := range data.StatCodes {
			if _, err := q.CreateHabitStatReward(ctx, db.CreateHabitStatRewardParams{
				HabitID:  habit.ID,
				StatCode: strings.ToUpper(strings.TrimSpace(stat)),
				Amount:   1,
			}); err != nil {
				return fmt.Errorf("add stat reward %s: %w", stat, err)
			}
		}
		if err != nil {
			return err
		}
		return nil
	})

	return habit, err
}

func (s *Service) GetHabits(ctx context.Context, userId int64) ([]db.Habit, error) {
	return s.q.ListActiveHabits(ctx, userId)
}

func (s *Service) Deactivate(ctx context.Context, habitId int64) error {
	_, err := s.q.DeactivateHabit(ctx, habitId)
	if errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("habit %d not found", habitId)
	}
	return err
}
