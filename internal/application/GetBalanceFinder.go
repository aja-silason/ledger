package application

import (
	"context"
	"errors"
	"log"

	"github.com/aja-silason/ledger/internal/domain"
	"github.com/aja-silason/ledger/internal/infra/postgres"
)

var (
	BalanceNotFoundError      = errors.New("Balanço da conta não encontrado")
	MustProvideBalanceIdError = errors.New("ID precisa ser fornecido")
)

type GetBalanceFinder struct {
	gateway *postgres.BalanceGateway
}

func NewGetBalanceFinder(gateway *postgres.BalanceGateway) *GetBalanceFinder {
	return &GetBalanceFinder{gateway: gateway}
}

func (b *GetBalanceFinder) FindBalanceByID(ctx context.Context, id string) (*domain.Balance, error) {

	if id == "" {
		return nil, MustProvideBalanceIdError
	}

	balance, err := b.gateway.FindByID(id)
	if err != nil {
		log.Printf("Erro ao obter o balance %v", err)
		return nil, BalanceNotFoundError
	}

	return balance, nil
}

func (b *GetBalanceFinder) FindBalanceByAccountID(ctx context.Context, id string) (*domain.Balance, error) {

	if id == "" {
		return nil, MustProvideBalanceIdError
	}

	balance, err := b.gateway.FindByAccountId(id)
	if err != nil {
		log.Printf("Erro ao obter o balance %v", err)
		return nil, BalanceNotFoundError
	}

	return balance, nil
}
