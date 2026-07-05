package user

import (
	"context"
	"database/sql"
	"errors"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/yourname/hunter-system/internal/db"
)

type Service struct {
	q db.Querier
}

func NewService(q db.Querier) *Service {
	return &Service{q: q}
}

func (s *Service) RegisterUser(ctx context.Context, telegramId int64, username string) (db.User, error) {
	u, err := s.q.GetUserByTelegramID(ctx, telegramId)
	if errors.Is(err, sql.ErrNoRows) {
		return s.createUser(ctx, telegramId, username)
	}
	return u, err
}

func (s *Service) createUser(ctx context.Context, telegramId int64, username string) (db.User, error) {
	u, err := s.q.CreateUser(ctx, db.CreateUserParams{
		telegramId,
		pgtype.Text{username, true},
		//TODO: Got timezone from user
		"+3",
	})
	return u, err
}
