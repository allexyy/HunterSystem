package quest

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/yourname/hunter-system/internal/db"
	"time"
)

type TxRunner interface {
	Transaction(ctx context.Context, fn func(q db.Querier) error) error
}

func updateStreak(tx db.Querier, ctx context.Context, q db.Quest, nowStreak, longestStreak int32) error {
	h, err := tx.GetHabitByID(ctx, q.HabitID.Int64)
	if err != nil {
		return fmt.Errorf("Cant Get Habit: %v", err)
	}

	tx.UpdateHabitStreak(ctx, db.UpdateHabitStreakParams{
		ID:            h.ID,
		CurrentStreak: nowStreak,
		LongestStreak: longestStreak,
		LastCompletedDate: pgtype.Date{
			Time:  time.Now(),
			Valid: true,
		},
	})
	return nil
}
