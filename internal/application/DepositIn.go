package application

import (
	"context"
	"log"
	"time"

	"github.com/aja-silason/ledger/internal/domain"
	"github.com/aja-silason/ledger/internal/infra/postgres"
)

type DepositIn struct {
	repo        *postgres.BalanceRepository
	accountRepo *postgres.AccountRepository
	// transactionRepo *postgres.TransactionRepository
}

func NewDepositIn(
	repo *postgres.BalanceRepository,
	accountRepo *postgres.AccountRepository,
	// transactionRepo *postgres.TransactionRepository
) *DepositIn {
	return &DepositIn{
		repo:        repo,
		accountRepo: accountRepo,
		// transactionRepo: transactionRepo,
	}
}

var (
	SuccessMessage               = "Depósito realizado com sucesso"
	WasNotPossibleDepositInError = "Não foi possível efectuar o deposito"
)

func (d *DepositIn) Deposit(ctx context.Context, input *domain.DepositInput) (string, error) {

	account, err := d.accountRepo.FindByID(input.AccountId)
	if err != nil {
		return "", err
	}

	balance, err := d.repo.FindByAccountId(input.AccountId)
	if err != nil {
		return "", err
	}

	newAmount := balance.CurrentAmount + input.Amount

	now := time.Now().UTC()
	save := &domain.Balance{
		ID:            balance.ID,
		AccountID:     account.ID,
		CurrencyCode:  balance.CurrencyCode,
		CurrentAmount: newAmount,
		UpdatedAt:     now,
	}

	_, err = d.repo.Update(save)
	if err != nil {
		log.Printf("[ERRO CRIAÇÃO] Falha ao executar o depósito")
		return WasNotPossibleDepositInError, err
	}

	return SuccessMessage, nil
}
