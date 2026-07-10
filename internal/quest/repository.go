package quest

import (
	"context"
	"github.com/yourname/hunter-system/internal/db"
)

type TxRunner interface {
	Transaction(ctx context.Context, fn func(q db.Querier) error) error
}
