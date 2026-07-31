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
	ListHasError         = errors.New("Sistema não consegue exibir a lista")
)

type GetAccountServiceFinder struct {
	gateway *postgres.AccountGateway
}

func NewGetAccountServiceFinder(gateway *postgres.AccountGateway) *GetAccountServiceFinder {
	return &GetAccountServiceFinder{gateway: gateway}
}

func (u *GetAccountServiceFinder) FindAccountByID(ctx context.Context, id string) (*domain.Account, error) {

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

func (u *GetAccountServiceFinder) FindAllAccount(ctx context.Context) (*domain.Account, error) {
	accounts, err := u.gateway.FindAll()
	if err != nil {
		log.Print("Erro aqui. O erro no findAll: %v", err)
		return nil, ListHasError
	}

	return accounts, nil
}
