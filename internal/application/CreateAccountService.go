package application

import (
	"context"
	"errors"
	"log"

	"github.com/aja-silason/ledger/internal/domain"
	"github.com/aja-silason/ledger/internal/infra/postgres"
)

var (
	createdAccountErro = errors.New("Falha ao criar a conta")
)

type CreateAccountService struct {
	repo *postgres.AccountRepository
}

func NewCreateAccountService(repo *postgres.AccountRepository) *CreateAccountService {
	return &CreateAccountService{repo: repo}
}

func (u *CreateAccountService) PostCreateAccount(ctx context.Context, tx *domain.Account) (*domain.Account, error) {

	if err := tx.ValidateAccount(); err != nil {
		return nil, err
	}

	if err := tx.ValidateAccountType(string(tx.Type)); err != nil {
		return nil, err
	}

	_, err := u.repo.Save(tx)
	if err != nil {
		log.Printf("[ERRO BANCO DE DADOS] Falha no Save: %v", err)
		return nil, createdAccountErro
	}

	return tx, nil
}
