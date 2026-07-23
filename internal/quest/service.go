package quest

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/yourname/hunter-system/internal/db"
	"github.com/yourname/hunter-system/internal/stats"
	"time"
)

type Service struct {
	q  db.Querier
	tx TxRunner
}

func NewService(q db.Querier, tx TxRunner) *Service {
	return &Service{q: q, tx: tx}
}

func (s *Service) FinishNotCompleteTask(ctx context.Context, userId int64) error {
	quests, err := s.q.ListDailyQuestsByDate(ctx, db.ListDailyQuestsByDateParams{
		UserID: userId,
		DueDate: pgtype.Date{
			Time:  time.Now().AddDate(0, 0, -1),
			Valid: true,
		},
	})
	if err != nil {
		return fmt.Errorf("list daily quests: %w", err)
	}
	for _, quest := range quests {
		s.tx.Transaction(ctx, func(q db.Querier) error {
			if quest.DeadlineAt.Time.Before(time.Now()) {
				h, err := q.GetHabitByID(ctx, quest.HabitID.Int64)
				if err != nil {
					return fmt.Errorf("Cant Get Habit: %v", err)
				}
				updateStreak(q, ctx, quest, 0, h.LongestStreak)
				if err != nil {
					return fmt.Errorf("Cant Update streak: %v", err)
				}
			}
			return err
		})
	}
	return nil
}

func (s *Service) GenerateDailyQuests(ctx context.Context, userId int64) ([]db.Quest, error) {
	h, err := s.q.ListActiveHabits(ctx, userId)
	if err != nil {
		return nil, fmt.Errorf("Can't load active habits: %v", err)
	}
	var quests []db.Quest
	err = s.tx.Transaction(ctx, func(q db.Querier) error {
		for _, h := range h {
			q, err := s.CreateQuest(ctx, h.UserID, h.ID, h.Title, h.Description.String, h.XpReward, h.GoldReward)
			if err != nil {
				err = fmt.Errorf("Can't create active quest: %v", err)
			}
			quests = append(quests, *q)
		}
		return err
	})
	return quests, err
}

func (s *Service) CompleteQuest(ctx context.Context, questId int64) error {
	q, err := s.q.GetQuestByID(ctx, questId)
	if err != nil {
		return fmt.Errorf("Quest not found: %v", err)
	}
	s.tx.Transaction(ctx, func(tx db.Querier) error {
		_, err := tx.CompleteQuest(ctx, q.ID)
		if err != nil {
			return fmt.Errorf("Cant complete quest: %v", err)
		}
		u, err := tx.UpdateUserXPGold(ctx, db.UpdateUserXPGoldParams{
			ID:   q.UserID,
			Xp:   int64(q.XpReward),
			Gold: int64(q.GoldReward),
		})
		lvl, xp := stats.LvlEncrease(int(u.Level), int(u.Xp))
		tx.UpdateUserLvl(ctx, db.UpdateUserLvlParams{
			ID:    u.ID,
			Xp:    xp,
			Level: lvl,
		})
		if err != nil {
			return fmt.Errorf("Cant Update xp and gold: %v", err)
		}
		stat, err := tx.ListQuestStatRewards(ctx, q.ID)
		if err != nil {
			return fmt.Errorf("Cant Get stats: %v", err)
		}
		for _, s := range stat {
			tx.UpsertUserStat(ctx, db.UpsertUserStatParams{
				UserID:   u.ID,
				StatCode: s.StatCode,
				Value:    s.Amount,
			})
		}
		h, err := tx.GetHabitByID(ctx, q.HabitID.Int64)
		if err != nil {
			return fmt.Errorf("Cant Get Habit: %v", err)
		}

		nowStreak := h.CurrentStreak + 1
		var longestStreak int32
		if nowStreak > h.LongestStreak {
			longestStreak = nowStreak
		} else {
			longestStreak = h.LongestStreak
		}
		updateStreak(tx, ctx, q, nowStreak, longestStreak)
		if err != nil {
			return fmt.Errorf("Cant Update streak: %v", err)
		}
		return err
	})
	return nil
}

func (s *Service) GetQuestListWithGeneration(ctx context.Context, userId int64) (quests []db.Quest, err error) {
	q, err := s.q.ListDailyQuestsByDate(ctx, db.ListDailyQuestsByDateParams{
		UserID:  userId,
		DueDate: pgtype.Date{Time: time.Now(), Valid: true},
	})
	if len(q) == 0 {
		fmt.Println("No quests found. Start generate")
		q, err = s.GenerateDailyQuests(ctx, userId)
	}
	if err != nil {
		return nil, fmt.Errorf("Quests not found: %v", err)
	}
	return q, err
}

func (s *Service) GetQuestList(ctx context.Context, userId int64) (quests []db.Quest, err error) {
	q, err := s.q.ListDailyQuestsByDate(ctx, db.ListDailyQuestsByDateParams{
		UserID:  userId,
		DueDate: pgtype.Date{Time: time.Now(), Valid: true},
	})
	return q, err
}

func (s *Service) CreateQuest(ctx context.Context, userId int64, habitId int64, title, description string, xpReward, goldReward int32) (*db.Quest, error) {
	tomorrow := time.Now().AddDate(0, 0, 1)
	q, err := s.q.CreateQuest(ctx, db.CreateQuestParams{
		UserID:      userId,
		HabitID:     pgtype.Int8{Int64: habitId, Valid: true},
		Type:        db.QuestTypeDaily,
		Title:       title,
		Description: pgtype.Text{String: description, Valid: true},
		XpReward:    xpReward,
		GoldReward:  goldReward,
		DueDate:     pgtype.Date{Time: tomorrow, Valid: true},
		DeadlineAt:  pgtype.Timestamptz{Time: tomorrow, Valid: true},
	})
	if err != nil {
		return nil, fmt.Errorf("Can't create quest: %v", err)
	}
	return &q, err
}
