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

	getUserFn    func(ctx context.Context, telegramID int64) (db.User, error)
	createUserFn func(ctx context.Context, arg db.CreateUserParams) (db.User, error)
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
			return existing, nil // юзер уже есть
		},
	}

	svc := user.NewService(m)
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
			return db.User{ID: 1, TelegramID: arg.TelegramID}, nil
		},
	}
	svc := user.NewService(m)
	got, err := svc.RegisterUser(context.Background(), 42, "alex")
	if err != nil {
		t.Fatalf("RegisterUser() unexpected error: %v", err)
	}
	if got.ID != 1 {
		t.Errorf("RegisterUser() ID = %d, want %d", got.ID, 1)
	}
}

func TestRegisterUser_DBError(t *testing.T) {
	m := &mockQuerier{
		getUserFn: func(ctx context.Context, id int64) (db.User, error) {
			return db.User{}, errors.New("connection refused") // НЕ ErrNoRows
		},
	}

	svc := user.NewService(m)
	_, err := svc.RegisterUser(context.Background(), 42, "alex")
	if err == nil {
		t.Fatalf("RegisterUser() expected error but got empty:")
	}
}
