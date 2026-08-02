package application

import (
	"context"
	"errors"
	"time"

	"github.com/aja-silason/ledger/internal/domain"
	"github.com/aja-silason/ledger/internal/infra/postgres"
	"github.com/google/uuid"
)

type DrawedCardLess struct {
	repo            *postgres.WithdrawRepository
	transactionRepo *postgres.TransactionRepository
	entries         *postgres.EntriesRepository
}

func NewDrawedCardLess(
	repo *postgres.WithdrawRepository,
	transactionRepo *postgres.TransactionRepository,
	entries *postgres.EntriesRepository,
) *DrawedCardLess {
	return &DrawedCardLess{
		repo:            repo,
		transactionRepo: transactionRepo,
		entries:         entries,
	}
}

var ThatOperationIsAlreadyRealized = errors.New("Esta operação já foi realizada")

func (c *DrawedCardLess) Drawed(ctx context.Context, key string, input *DrawedCardLessInput) (domain.SuccessMessage, error) {
	existing, err := c.transactionRepo.FindByIdempotencyKey(key)
	if err == nil && existing != nil {
		return nil, ThatOperationIsAlreadyRealized
	}

	withdraw, err := c.repo.FindByCode(input.Reference)
	if err != nil {
		return nil, err
	}

	drawed, err := withdraw.Drawed()
	if err != nil {
		return nil, err
	}

	err = c.transactionByDrawed(ctx, key, withdraw.AccountID.String(), withdraw.Amount)
	if err != nil {
		return nil, err
	}

	_, err = c.repo.UpdateStatus(withdraw.ID.String(), string(drawed.Status))
	if err != nil {
		return nil, err
	}

	return domain.NewSuccessMessage("Levantamento sem cartão realizado com sucesso."), nil
}

func (c *DrawedCardLess) transactionByDrawed(ctx context.Context, key, accountId string, amount int64) error {

	transactionID := uuid.New()
	now := time.Now().UTC()

	legDebit := domain.TransactionLeg{
		ID:            uuid.New(),
		TransactionID: transactionID,
		AccountID:     uuid.MustParse(accountId),
		Direction:     domain.Debit,
		Amount:        amount,
	}

	legs := []domain.TransactionLeg{legDebit}

	transaction := &domain.Transaction{
		ID:             transactionID,
		IdempotencyKey: key,
		Description:    "Levantamento sem cartão - Saque directo da conta",
		CreatedAt:      now,
		Legs:           legs,
	}

	_, err := c.transactionRepo.Save(ctx, transaction)
	if err != nil {
		return err
	}
	return nil
}
