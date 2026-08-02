package application

import (
	"context"

	"github.com/aja-silason/ledger/internal/domain"
	"github.com/aja-silason/ledger/internal/infra/postgres"
)

type CancelWithDrawCardLess struct {
	repo *postgres.WithdrawRepository
}

func NewCancelWithDrawCardLess(
	repo *postgres.WithdrawRepository,
) *CancelWithDrawCardLess {
	return &CancelWithDrawCardLess{repo: repo}
}

func (c *CancelWithDrawCardLess) Cancel(ctx context.Context, input *CancelWithDrawCardLessInput) (domain.SuccessMessage, error) {
	withdraw, err := c.repo.FindById(input.WithDrawID)
	if err != nil {
		return nil, err
	}

	canceled, err := withdraw.Canceled()
	if err != nil {
		return nil, err
	}
	_, err = c.repo.UpdateStatus(input.WithDrawID, string(canceled.Status))
	if err != nil {
		return nil, err
	}

	return domain.NewSuccessMessage("Levantamento sem cartão cancelado"), nil
}
