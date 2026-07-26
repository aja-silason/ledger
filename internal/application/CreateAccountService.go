package application

import (
	"context"
	"errors"
	"time"

	"github.com/aja-silason/ledger/internal/domain"
	"github.com/google/uuid"
)

type AccountRepository interface {
	CreateAccount(ctx context.Context, tx *domain.Account) error
}

type CreateAccountService struct {
	repo AccountRepository
}

func NewCreateAccountService(repo AccountRepository) *CreateAccountService {
	return &CreateAccountService{repo: repo}
}

func (u *CreateAccountService) PostCreateAccount(ctx context.Context, tx *domain.Account) (*domain.Account, error) {

	if err := tx.ValidateAccount(); err == nil {
		return nil, err
	}

	tx.ID = uuid.New()
	tx.CreatedAt = time.Now()
	tx.Name = "Conta de Test"
	tx.Type = domain.Equity

	if err := u.repo.CreateAccount(ctx, tx); err != nil {
		return nil, errors.New("Falha ao criar a conta")
	}

	return tx, nil
}
