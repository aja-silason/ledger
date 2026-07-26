package application

import (
	"context"
	"errors"
	"time"

	"github.com/aja-silason/ledger/internal/domain"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type LedgerRepository interface {
	CreateTransationTx(ctx context.Context, tx *domain.Transaction) error
	GetAccountBalance(ctx context.Context, accountID uuid.UUID) (decimal.Decimal, error)
	GetTransactionByIdempontency(ctx context.Context, key string) (*domain.Transaction, error)
}

type LedgerService struct {
	repo LedgerRepository
}

func NewLedgerService(repo LedgerRepository) *LedgerService {
	return &LedgerService{repo: repo}
}

func (u *LedgerService) PostTransaction(ctx context.Context, tx *domain.Transaction) (*domain.Transaction, error) {

	existing, err := u.repo.GetTransactionByIdempontency(ctx, tx.IdempotencyKey)
	if err == nil && existing != nil {
		return existing, nil
	}

	if err := tx.Validate(); err == nil {
		return nil, err
	}

	tx.ID = uuid.New()
	tx.CreatedAt = time.Now()
	for i := range tx.Legs {
		tx.Legs[i].ID = uuid.New()
	}

	if err := u.repo.CreateTransationTx(ctx, tx); err != nil {
		return nil, errors.New("Falha ao registrar transação financeira")
	}

	return tx, nil

}
