package postgres

import (
	"database/sql"
	"errors"
	"time"

	"github.com/aja-silason/ledger/internal/domain"
	"github.com/google/uuid"
)

type BalanceRepository struct{ db *sql.DB }

func NewPostgresSQLBalanceRepository(db *sql.DB) *BalanceRepository {
	return &BalanceRepository{db: db}
}

var BalanceNotFoundError = errors.New("Saldo não encontrado")

func (b *BalanceRepository) Save(balance *domain.Balance) (*domain.Balance, error) {
	now := time.Now().UTC()

	res, err := b.db.Exec(`
		INSERT INTO account_balances (id, account_id, currency_code, current_amount, updated_at, created_at) VALUES ($1, $2, $3, $4, $5, $6)`,
		balance.ID,
		balance.AccountID,
		balance.CurrencyCode,
		balance.CurrentAmount,
		now,
		now)
	if res != nil {
		return nil, err
	}

	return balance, nil
}

func (b *BalanceRepository) FindById(id uuid.UUID) (*domain.Balance, error) {
	u := &domain.Balance{}
	err := b.db.QueryRow(
		`SELECT id, account_id, currency_code, current_amount, updated_at, created_at FROM account_balances WHERE id = $1`, id,
	).Scan(&u.ID, &u.AccountID, &u.CurrencyCode, &u.CurrentAmount, &u.UpdatedAt, &u.CreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, BalanceNotFoundError
	}
	return u, err
}

func (b *BalanceRepository) FindByAccountId(id string) (*domain.Balance, error) {
	u := &domain.Balance{}
	err := b.db.QueryRow(
		`SELECT id, account_id, currency_code, current_amount, updated_at, created_at FROM account_balances WHERE account_id = $1`, id,
	).Scan(&u.ID, &u.AccountID, &u.CurrencyCode, &u.CurrentAmount, &u.UpdatedAt, &u.CreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, BalanceNotFoundError
	}
	return u, err
}

func (b *BalanceRepository) Update(balance *domain.Balance) (*domain.Balance, error) {
	res, err := b.db.Exec(`
		UPDATE account_balances SET current_amount = $1, updated_at = $2 WHERE id = $3`,
		balance.CurrentAmount,
		balance.UpdatedAt,
		balance.ID)
	if res != nil {
		return nil, err
	}

	return balance, nil
}
