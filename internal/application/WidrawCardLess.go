package application

import (
	"context"

	"github.com/aja-silason/ledger/internal/domain"
	"github.com/aja-silason/ledger/internal/infra/postgres"
)

type WidrawCardLess struct {
	repo            *postgres.BalanceRepository
	accountRepo     *postgres.AccountRepository
	transictionRepo *postgres.TransactionRepository
}

func NewWidrawCardLess(
	repo *postgres.BalanceRepository,
	accountRepo *postgres.AccountRepository,
	transictionRepo *postgres.TransactionRepository,
) *WidrawCardLess {
	return &WidrawCardLess{
		repo:            repo,
		accountRepo:     accountRepo,
		transictionRepo: transictionRepo,
	}
}

func (w *WidrawCardLess) DemandWidraw(ctx context.Context) (domain.SuccessMessage, error) {
	return nil, nil
}
