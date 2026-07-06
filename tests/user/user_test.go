package user

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/yourname/hunter-system/internal/db"
	"github.com/yourname/hunter-system/internal/user"
	"testing"
)

type mockQuerier struct {
	db.Querier

	getUserFn       func(ctx context.Context, telegramID int64) (db.User, error)
	createUserFn    func(ctx context.Context, arg db.CreateUserParams) (db.User, error)
	initUserStatsFn func(ctx context.Context, userID int64) error
}

type mockTx struct{ q db.Querier }

func (m *mockTx) Transaction(ctx context.Context, fn func(q db.Querier) error) error {
	return fn(m.q)
}

func (m *mockQuerier) InitUserStats(ctx context.Context, userID int64) error {
	return m.initUserStatsFn(ctx, userID)
}

func (m *mockQuerier) GetUserByTelegramID(ctx context.Context, telegramID int64) (db.User, error) {
	return m.getUserFn(ctx, telegramID)
}

func (m *mockQuerier) CreateUser(ctx context.Context, arg db.CreateUserParams) (db.User, error) {
	return m.createUserFn(ctx, arg)
}

func TestRegisterUser_ExistingUser(t *testing.T) {
	existing := db.User{ID: 1, TelegramID: 42}

	m := &mockQuerier{
		getUserFn: func(ctx context.Context, id int64) (db.User, error) {
			return existing, nil
		},
	}

	svc := user.NewService(m, &mockTx{q: m})
	got, err := svc.RegisterUser(context.Background(), 42, "alex")

	if err != nil {
		t.Fatalf("RegisterUser() unexpected error: %v", err)
	}
	if got.ID != existing.ID {
		t.Errorf("RegisterUser() ID = %d, want %d", got.ID, existing.ID)
	}
}

func TestRegisterUser_NewUser(t *testing.T) {
	m := &mockQuerier{
		getUserFn: func(ctx context.Context, id int64) (db.User, error) {
			return db.User{}, pgx.ErrNoRows // юзера нет
		},
		createUserFn: func(ctx context.Context, arg db.CreateUserParams) (db.User, error) {
			return db.User{ID: 7, TelegramID: arg.TelegramID}, nil
		},
		initUserStatsFn: func(ctx context.Context, userID int64) error {
			if userID != 7 {
				t.Errorf("InitUserStats called with userID = %d, want 7", userID)
			}
			return nil
		},
	}
	svc := user.NewService(m, &mockTx{q: m})
	got, err := svc.RegisterUser(context.Background(), 42, "alex")
	if err != nil {
		t.Fatalf("RegisterUser() unexpected error: %v", err)
	}
	if got.ID != 7 {
		t.Errorf("RegisterUser() ID = %d, want %d", got.ID, 1)
	}
}

func TestRegisterUser_DBError(t *testing.T) {
	m := &mockQuerier{
		getUserFn: func(ctx context.Context, id int64) (db.User, error) {
			return db.User{}, errors.New("connection refused") // НЕ ErrNoRows
		},
	}

	svc := user.NewService(m, &mockTx{q: m})
	_, err := svc.RegisterUser(context.Background(), 42, "alex")
	if err == nil {
		t.Fatalf("RegisterUser() expected error but got empty:")
	}
}
