package application

import (
	"context"
	"errors"
	"log"

	"github.com/aja-silason/ledger/internal/domain"
	"github.com/aja-silason/ledger/internal/infra/postgres"
)

var (
	AccountNotFoundError = errors.New("Conta não encontrada")
	IDMustGiven          = errors.New("ID precisa ser fornecido")
)

type GetAccountService struct {
	gateway *postgres.AccountGateway
}

func NewGetAccountService(gateway *postgres.AccountGateway) *GetAccountService {
	return &GetAccountService{gateway: gateway}
}

func (u *GetAccountService) FindAccountByID(ctx context.Context, id string) (*domain.Account, error) {

	if id == "" {
		return nil, IDMustGiven
	}

	account, err := u.gateway.FindByID(id)
	if err != nil {
		log.Printf("Erro aqui. O erro: %v", err)
		return nil, AccountNotFoundError
	}

	return account, nil

}
