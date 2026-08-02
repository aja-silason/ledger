package application

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/aja-silason/ledger/internal/domain"
	"github.com/aja-silason/ledger/internal/infra/postgres"
	"github.com/aja-silason/ledger/internal/security"
	"github.com/google/uuid"
)

type DrawedCardLess struct {
	repo            *postgres.WithdrawRepository
	transactionRepo *postgres.TransactionRepository
	entries         *postgres.EntriesRepository
	balanceRepo     *postgres.BalanceRepository
}

func NewDrawedCardLess(
	repo *postgres.WithdrawRepository,
	transactionRepo *postgres.TransactionRepository,
	entries *postgres.EntriesRepository,
	balanceRepo *postgres.BalanceRepository,
) *DrawedCardLess {
	return &DrawedCardLess{
		repo:            repo,
		transactionRepo: transactionRepo,
		entries:         entries,
		balanceRepo:     balanceRepo,
	}
}

var (
	ThatOperationIsAlreadyRealized = errors.New("Esta operação já foi realizada")
	ThisCodeIsNotPendingErr        = errors.New("Este levantamento não está pendente")
	CodeNotMatchErr                = errors.New("Código secreto incorrecto, tente novamente com outra combinação!")
)

func (c *DrawedCardLess) Drawed(ctx context.Context, key string, input *DrawedCardLessInput) (domain.SuccessMessage, error) {
	existing, err := c.transactionRepo.FindByIdempotencyKey(key)
	if err == nil && existing != nil {
		return nil, ThatOperationIsAlreadyRealized
	}

	withdraw, err := c.repo.FindByCode(input.Reference)
	if err != nil {
		return nil, err
	}

	if err := security.ValidateHash(input.SecretCode, withdraw.CodeHash); err != true {
		return nil, CodeNotMatchErr
	}

	if withdraw != nil && !withdraw.IsPendig() {
		return nil, ThisCodeIsNotPendingErr
	}

	drawed, err := withdraw.Drawed()
	if err != nil {
		return nil, err
	}

	balance, err := c.balanceRepo.FindByAccountId(withdraw.AccountID.String())
	if err != nil {
		return nil, err
	}

	err = c.decreaseAmountDrawed(balance, withdraw.Amount)
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

func (c *DrawedCardLess) decreaseAmountDrawed(balance *domain.Balance, amount int64) error {
	decreaseAmount := balance.CurrentAmount - amount
	decreaseUpdated := time.Now().UTC()
	decrease := &domain.Balance{
		ID:            balance.ID,
		CurrentAmount: decreaseAmount,
		UpdatedAt:     decreaseUpdated,
	}
	_, err := c.balanceRepo.Update(decrease)
	if err != nil {
		log.Printf("[RETIRAR DA CONTA ORIGEM] Falha ao retirar da conta origem")
		return errors.New("Não foi possível retirar os valores da conta origem")
	}
	return nil
}
