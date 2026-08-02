// NOTA - QUIÇA eu faça depois uma ajuste básico
// Ao efectuar um depósito o dinheiro é dado em mão, mas no sistema, ao que parece ele surje por magia, como deveria funcionar isso?
// Hipotese 1: Ter-se-a uma conta principal da empresa detentora (Banco) onde se vai creditar esse valor quando alguém o deposita, visando na transação 2 LEG(entrada e saida)
// Hipotese 2: Para deposito não se prevê esse cenário

//  A principio não vou prever esse cenário, por conta de poder entregar isso, mais tarde eu dou um chega melhor nisso

package application

import (
	"context"
	"log"
	"time"

	"github.com/aja-silason/ledger/internal/domain"
	"github.com/aja-silason/ledger/internal/infra/postgres"
	"github.com/google/uuid"
)

type DepositIn struct {
	repo            *postgres.BalanceRepository
	accountRepo     *postgres.AccountRepository
	transactionRepo *postgres.TransactionRepository
}

func NewDepositIn(
	repo *postgres.BalanceRepository,
	accountRepo *postgres.AccountRepository,
	transactionRepo *postgres.TransactionRepository,
) *DepositIn {
	return &DepositIn{
		repo:            repo,
		accountRepo:     accountRepo,
		transactionRepo: transactionRepo,
	}
}

var (
	SuccessMessage               = "Depósito realizado com sucesso"
	WasNotPossibleDepositInError = "Não foi possível efectuar o deposito"
	IdempontecyError             = "Essa operação já foi realizada"
	TransactionError             = "Erro ao executar a transação"
)

func (d *DepositIn) Deposit(ctx context.Context, input *domain.DepositInput, key string) (string, error) {

	account, err := d.accountRepo.FindByID(input.AccountId)
	if err != nil {
		return "", err
	}

	balance, err := d.repo.FindByAccountId(input.AccountId)
	if err != nil {
		return "", err
	}

	// idempotencyKey := uuid.New()
	existing, err := d.transactionRepo.FindByIdempotencyKey(key)
	if err == nil && existing != nil {
		return IdempontecyError, err
	}

	newAmount := balance.CurrentAmount + input.Amount
	now := time.Now().UTC()
	saveBalance := &domain.Balance{
		ID:            balance.ID,
		AccountID:     account.ID,
		CurrencyCode:  balance.CurrencyCode,
		CurrentAmount: newAmount,
		UpdatedAt:     now,
		CreatedAt:     balance.CreatedAt,
	}

	transactionID := uuid.New()
	transactionCreatedAt := time.Now().UTC()
	saveTransaction := &domain.Transaction{
		ID:             transactionID,
		IdempotencyKey: key,
		Description:    "Depósito Numerário",
		Legs:           []domain.TransactionLeg{},
		CreatedAt:      transactionCreatedAt,
	}

	_, err = d.repo.Update(saveBalance)
	if err != nil {
		log.Printf("[ERRO CRIAÇÃO] Falha ao executar o depósito")
		return WasNotPossibleDepositInError, err
	}

	_, err = d.transactionRepo.Save(saveTransaction)
	if err != nil {
		log.Printf("[ERRO TRANSAÇÃO] Falha ao executar a transação")
		return TransactionError, nil
	}

	return SuccessMessage, nil
}
