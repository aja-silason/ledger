package postgres

import (
	"database/sql"
	"errors"
	"time"

	"github.com/aja-silason/ledger/internal/domain"
)

type WithdrawRepository struct{ db *sql.DB }

func NewPostgresSQLWithdrawRepository(db *sql.DB) *WithdrawRepository {
	return &WithdrawRepository{db: db}
}

var WithdrawFoundError = errors.New("Levantamento não encontrado")

func (b *WithdrawRepository) Save(wt *domain.Withdraw) (*domain.Withdraw, error) {
	res, err := b.db.Exec(`
		INSERT INTO withdraws (id, account_id, amount, status, code, code_hash, expires_at, creadet_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
		wt.ID, wt.AccountID, wt.Amount, wt.Status, wt.Code, wt.CodeHash, wt.ExpiresAt, wt.CreatedAt, wt.UpdatedAt)
	if res != nil {
		return nil, err
	}
	return wt, nil
}

func (b *WithdrawRepository) FindByCode(code int64) (*domain.Withdraw, error) {
	u := &domain.Withdraw{}
	err := b.db.QueryRow(
		`SELECT id, account_id, amount, status, code_hash, code, expires_at, creadet_at, updated_at FROM withdraws WHERE code = $1`, code,
	).Scan(&u.ID, &u.AccountID, &u.Amount, &u.Status, &u.CodeHash, &u.Code, &u.ExpiresAt, &u.CreatedAt, &u.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, BalanceNotFoundError
	}
	return u, err
}

func (b *WithdrawRepository) FindById(id string) (*domain.Withdraw, error) {
	u := &domain.Withdraw{}
	err := b.db.QueryRow(
		`SELECT id, account_id, amount, status, code_hash, code, expires_at, creadet_at, updated_at FROM withdraws WHERE id = $1`, id,
	).Scan(&u.ID, &u.AccountID, &u.Amount, &u.Status, &u.CodeHash, &u.Code, &u.ExpiresAt, &u.CreatedAt, &u.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, BalanceNotFoundError
	}
	return u, err
}

func (b *WithdrawRepository) FindByAccountId(id string) (*domain.Withdraw, error) {
	u := &domain.Withdraw{}
	err := b.db.QueryRow(
		`SELECT id, account_id, amount, status, code_hash, code, expires_at, creadet_at, updated_at FROM withdraws WHERE id = $1`, id,
	).Scan(&u.ID, &u.AccountID, &u.Amount, &u.Status, &u.CodeHash, &u.Code, &u.ExpiresAt, &u.CreatedAt, &u.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, BalanceNotFoundError
	}
	return u, err
}

func (b *WithdrawRepository) UpdateStatus(id string, status string) (*domain.Withdraw, error) {
	now := time.Now().UTC()
	res, err := b.db.Exec(`
		UPDATE withdraws SET status = $1, updated_at = $2 WHERE id = $3`,
		status,
		now,
		id)
	if res != nil {
		return nil, err
	}
	return nil, nil
}
