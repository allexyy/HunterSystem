package user

import (
	"context"
	"database/sql"
	"errors"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/yourname/hunter-system/internal/db"
)

type Service struct {
	q  db.Querier
	tx TxRunner
}

func NewService(q db.Querier, tx TxRunner) *Service {
	return &Service{q: q, tx: tx}
}

func (s *Service) RegisterUser(ctx context.Context, telegramId int64, username string) (db.User, error) {
	u, err := s.q.GetUserByTelegramID(ctx, telegramId)
	if errors.Is(err, sql.ErrNoRows) {
		return s.createUser(ctx, telegramId, username)
	}
	return u, err
}

func (s *Service) GetUser(ctx context.Context, telegramId int64) (db.User, error) {
	u, err := s.q.GetUserByTelegramID(ctx, telegramId)
	return u, err
}

func (s *Service) createUser(ctx context.Context, telegramId int64, username string) (db.User, error) {
	var created db.User
	err := s.tx.Transaction(ctx, func(q db.Querier) error {
		var err error
		created, err = q.CreateUser(ctx, db.CreateUserParams{
			telegramId,
			pgtype.Text{username, true},
			//TODO: Got timezone from user
			"+3",
		})
		if err != nil {
			return err
		}
		return q.InitUserStats(ctx, created.ID)
	})
	if err != nil {
		return db.User{}, err
	}
	return created, nil
}
