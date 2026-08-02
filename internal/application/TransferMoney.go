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

type TransferMoney struct {
	repo            *postgres.BalanceRepository
	accountRepo     *postgres.AccountRepository
	transactionRepo *postgres.TransactionRepository
	entriesRepo     *postgres.EntriesRepository
}

func NewTransferMoney(
	repo *postgres.BalanceRepository,
	accountRepo *postgres.AccountRepository,
	transactionRepo *postgres.TransactionRepository,
	entriesRepo *postgres.EntriesRepository,
) *TransferMoney {
	return &TransferMoney{
		repo:            repo,
		accountRepo:     accountRepo,
		transactionRepo: transactionRepo,
		entriesRepo:     entriesRepo,
	}
}

var (
	ErrInsufficientBalance          = errors.New("saldo insuficiente para realizar a transferência")
	ErrInvalidAmount                = errors.New("o valor da transferência deve ser maior que zero")
	ErrSameAccount                  = errors.New("a conta de origem e destino não podem ser a mesma")
	ErrOperationHasAlreadyConcluted = errors.New("Esta operação já foi realizada")
)

func (d *TransferMoney) Transfer(ctx context.Context, input *TransferMoneyInput, key string) (domain.SuccessMessage, error) {
	if input.Amount <= 0 {
		return nil, ErrInvalidAmount
	}

	if input.FromAccountID == input.ToAccountID {
		return nil, ErrSameAccount
	}

	fromAccount, err := d.accountRepo.FindByID(input.FromAccountID)
	if err != nil {
		return nil, err
	}

	fromBalance, err := d.repo.FindByAccountId(input.FromAccountID)
	if err != nil {
		return nil, err
	}

	toAccount, err := d.accountRepo.FindByID(input.ToAccountID)
	if err != nil {
		return nil, err
	}

	toBalance, err := d.repo.FindByAccountId(input.ToAccountID)
	if err != nil {
		return nil, err
	}

	existing, err := d.transactionRepo.FindByIdempotencyKey(key)
	if err == nil && existing != nil {
		return nil, ErrOperationHasAlreadyConcluted
	}

	err = d.validateTransfer(fromBalance, input.Amount)
	if err != nil {
		return nil, err
	}

	err = d.decreadeAmount(fromBalance, input.Amount)
	if err != nil {
		return nil, err
	}

	err = d.encreaseAmount(toBalance, input.Amount)
	if err != nil {
		return nil, err
	}

	err = d.transaction(key, fromAccount.ID, toAccount.ID, input.Amount)
	if err != nil {
		return nil, err
	}

	return domain.NewSuccessMessage("Transferência realizada com sucesso"), nil
}

func (d *TransferMoney) validateTransfer(fromBalance *domain.Balance, amount int64) error {
	if fromBalance.CurrentAmount < amount {
		return ErrInsufficientBalance
	}
	return nil
}

func (d *TransferMoney) decreadeAmount(fromBalance *domain.Balance, amount int64) error {

	decreaseAmount := fromBalance.CurrentAmount - amount
	decreaseUpdated := time.Now().UTC()
	decrease := &domain.Balance{
		ID:            fromBalance.ID,
		CurrentAmount: decreaseAmount,
		UpdatedAt:     decreaseUpdated,
	}
	_, err := d.repo.Update(decrease)
	if err != nil {
		log.Printf("[RETIRAR DA CONTA ORIGEM] Falha ao retirar da conta origem")
		return errors.New("Não foi possível retirar os valores da conta origem")
	}

	return nil
}

func (d *TransferMoney) encreaseAmount(toBalance *domain.Balance, amount int64) error {

	encreaseAmount := toBalance.CurrentAmount - amount
	encreaseUpdated := time.Now().UTC()
	encrease := &domain.Balance{
		ID:            toBalance.ID,
		CurrentAmount: encreaseAmount,
		UpdatedAt:     encreaseUpdated,
	}
	_, err := d.repo.Update(encrease)
	if err != nil {
		log.Printf("[DEPOSITAR NA CONTA DESTINO] Falha ao depositar na conta destino")
		return errors.New("Não foi possível depositar valores na conta destino")
	}

	return nil
}

func (d *TransferMoney) transaction(key string, fromAccountId, toAccountId uuid.UUID, amount int64) error {

	transactionId := uuid.New()
	now := time.Now().UTC()

	legDebit := domain.TransactionLeg{
		ID:            uuid.New(),
		TransactionID: transactionId,
		AccountID:     fromAccountId,
		Direction:     domain.Debit,
	}

	legCredit := domain.TransactionLeg{
		ID:            uuid.New(),
		TransactionID: transactionId,
		AccountID:     toAccountId,
		Direction:     domain.Credit,
	}

	legs := []domain.TransactionLeg{legDebit, legCredit}

	savedTransaction := &domain.Transaction{
		ID:             transactionId,
		IdempotencyKey: key,
		Description:    "Transaferência de conta à conta",
		Legs:           legs,
		CreatedAt:      now,
	}

	_, err := d.transactionRepo.Save(savedTransaction)
	if err != nil {
		return err
	}
	return nil
}
