package habit

import (
	"context"
	"github.com/yourname/hunter-system/internal/db"
	"testing"
)

type mockQuerier struct {
	db.Querier

	createStatRewFn func(ctx context.Context, arg db.CreateHabitStatRewardParams) (db.HabitStatReward, error)
	createHabitFn   func(ctx context.Context, arg db.CreateHabitParams) (db.Habit, error)
	listHabitsFn    func(ctx context.Context, userID int64) ([]db.Habit, error)
}

type mockTx struct{ q db.Querier }

func (m *mockTx) Transaction(ctx context.Context, fn func(q db.Querier) error) error {
	return fn(m.q)
}

func (m *mockQuerier) CreateHabit(ctx context.Context, arg db.CreateHabitParams) (db.Habit, error) {
	return m.createHabitFn(ctx, arg)
}
func (m *mockQuerier) CreateHabitStatReward(ctx context.Context, arg db.CreateHabitStatRewardParams) (db.HabitStatReward, error) {
	return m.createStatRewFn(ctx, arg)
}

func (m *mockQuerier) ListActiveHabits(ctx context.Context, userID int64) ([]db.Habit, error) {
	return m.listHabitsFn(ctx, userID)
}

func TestCreateHabit(t *testing.T) {

	m := &mockQuerier{
		createHabitFn: func(ctx context.Context, arg db.CreateHabitParams) (db.Habit, error) {
			return db.Habit{}, nil
		},
		createStatRewFn: func(ctx context.Context, arg db.CreateHabitStatRewardParams) (db.HabitStatReward, error) {
			return db.HabitStatReward{}, nil
		},
	}

	svc := NewService(m, &mockTx{q: m})
	d := NewHabitData("Create |Create Test |Hard |Int ")
	_, err := svc.CreateHabit(context.Background(), 42, d)

	if err != nil {
		t.Fatalf("CreateHabit() unexpected error: %v", err)
	}
}

func TestGetHabit(t *testing.T) {
	m := &mockQuerier{
		createHabitFn: func(ctx context.Context, arg db.CreateHabitParams) (db.Habit, error) {
			return db.Habit{}, nil
		},
		listHabitsFn: func(ctx context.Context, userID int64) ([]db.Habit, error) {
			return []db.Habit{{Title: "Test"}}, nil
		},
	}
	svc := NewService(m, &mockTx{q: m})
	h, err := svc.GetHabits(context.Background(), 42)

	if err != nil {
		t.Fatalf("GetHabits() unexpected error: %v", err)
	}
	if len(h) == 0 {
		t.Fatalf("GetHabits() return empty list")
	}

}
