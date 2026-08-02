package postgres

import (
	"database/sql"
	"errors"

	"github.com/aja-silason/ledger/internal/domain"
)

type BalanceGateway struct{ db *sql.DB }

func NewPostgresSQLBalanceGateway(db *sql.DB) *BalanceGateway {
	return &BalanceGateway{db: db}
}

var GetBalanceNotFoundError = errors.New("Saldo não encontrado")

func (b *BalanceGateway) FindByID(id string) (*domain.Balance, error) {
	u := &domain.Balance{}
	err := b.db.QueryRow(
		`SELECT id, account_id, currency_code, current_amount, updated_at, created_at FROM account_balances WHERE id = $1`, id,
	).Scan(&u.ID, &u.AccountID, &u.CurrencyCode, &u.CurrentAmount, &u.UpdatedAt, &u.CreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, BalanceNotFoundError
	}
	return u, err
}

func (b *BalanceGateway) FindByAccountId(id string) (*domain.Balance, error) {
	u := &domain.Balance{}
	err := b.db.QueryRow(
		`SELECT id, account_id, currency_code, current_amount, updated_at, created_at FROM account_balances WHERE account_id = $1`, id,
	).Scan(&u.ID, &u.AccountID, &u.CurrencyCode, &u.CurrentAmount, &u.UpdatedAt, &u.CreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, GetBalanceNotFoundError
	}
	return u, err
}
