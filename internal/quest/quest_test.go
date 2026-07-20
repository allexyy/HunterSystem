package quest

import (
	"context"
	"github.com/yourname/hunter-system/internal/db"
	"testing"
)

type mockQuerier struct {
	db.Querier

	createQuestFn func(ctx context.Context, arg db.CreateQuestParams) (db.Quest, error)
}

type mockTx struct{ q db.Querier }

func (m *mockTx) Transaction(ctx context.Context, fn func(q db.Querier) error) error {
	return fn(m.q)
}
func (m *mockQuerier) CreateQuest(ctx context.Context, arg db.CreateQuestParams) (db.Quest, error) {
	return m.createQuestFn(ctx, arg)
}

func TestCreateQuest(t *testing.T) {

	m := &mockQuerier{
		createQuestFn: func(ctx context.Context, arg db.CreateQuestParams) (db.Quest, error) {
			return db.Quest{}, nil
		},
	}

	svc := NewService(m, &mockTx{q: m})
	_, err := svc.CreateQuest(context.Background(), 42, 1, "", "", 1, 1)

	if err != nil {
		t.Fatalf("CreateHabit() unexpected error: %v", err)
	}
}
