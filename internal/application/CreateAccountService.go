package application

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/aja-silason/ledger/internal/domain"
	"github.com/aja-silason/ledger/internal/infra/postgres"
	"github.com/google/uuid"
)

var (
	createdAccountError = errors.New("Falha ao criar a conta")
)

type CreateAccountService struct {
	repo        *postgres.AccountRepository
	balanceRepo *postgres.BalanceRepository
}

func NewCreateAccountService(repo *postgres.AccountRepository, balanceRepo *postgres.BalanceRepository) *CreateAccountService {
	return &CreateAccountService{repo: repo, balanceRepo: balanceRepo}
}

func (u *CreateAccountService) PostCreateAccount(ctx context.Context, tx *domain.Account) (*domain.Account, error) {

	id := uuid.New()

	tx.ID = id

	if err := tx.ValidateAccount(); err != nil {
		return nil, err
	}

	if err := tx.ValidateAccountType(); err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	accountBalance := &domain.Balance{
		ID:            uuid.New(),
		AccountID:     id,
		CurrencyCode:  "AOA",
		CurrentAmount: 0,
		UpdatedAt:     now,
		CreatedAt:     now,
	}

	if err := accountBalance.ValidateCurrency(); err != nil {
		log.Printf("[ERRO VALIDAÇÃO] Falha na validação da moeda: %v", err)
		return nil, err
	}

	_, err := u.repo.Save(tx)
	if err != nil {
		log.Printf("[ERRO BANCO DE DADOS] Falha no Save: %v", err)
		return nil, createdAccountError
	}

	_, err = u.balanceRepo.Save(accountBalance)
	if err != nil {
		log.Printf("[ERRO BANCO DE DADOS] Falha no Save: %v", err)
		return nil, createdAccountError
	}

	return tx, nil
}
