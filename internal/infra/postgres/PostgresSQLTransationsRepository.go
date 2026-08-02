package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/aja-silason/ledger/internal/domain"
)

type TransactionRepository struct{ db *sql.DB }

func NewPostgresSQLTransactionRepository(db *sql.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

var TransactionNotFound = errors.New("Transação não encontrada")

func (t *TransactionRepository) Save(ctx context.Context, transaction *domain.Transaction) (*domain.Transaction, error) {

	tx, err := t.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	defer tx.Rollback()

	// _, err := t.db.Exec(`
	// 	INSERT INTO transactions (id, idempotency_key, description, created_at) VALUES ($1, $2, $3, $4)`,
	// 	transaction.ID,
	// 	transaction.IdempotencyKey,
	// 	transaction.Description,
	// 	now)
	// if err != nil {
	// 	return nil, err
	// }

	_, err = tx.ExecContext(ctx, `
        INSERT INTO transactions (id, idempotency_key, description, created_at) 
        VALUES ($1, $2, $3, $4)`,
		transaction.ID,
		transaction.IdempotencyKey,
		transaction.Description,
		transaction.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	if len(transaction.Legs) > 0 {
		queryLeg := `
            INSERT INTO entries (id, transaction_id, account_id, amount, direction)
            VALUES ($1, $2, $3, $4, $5)`

		for _, leg := range transaction.Legs {
			_, err = tx.ExecContext(ctx, queryLeg,
				leg.ID,
				leg.TransactionID,
				leg.AccountID,
				leg.Amount,
				leg.Direction,
			)
			if err != nil {
				return nil, err
			}
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return transaction, nil
}

func (t *TransactionRepository) FindByIdempotencyKey(idempotencyKey string) (*domain.Transaction, error) {
	transaction := &domain.Transaction{}
	err := t.db.QueryRow(`
		SELECT id, idempotency_key, description, created_at FROM transactions WHERE idempotency_key = $1`,
		idempotencyKey).Scan(
		&transaction.ID,
		&transaction.IdempotencyKey,
		&transaction.Description,
		&transaction.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, TransactionNotFound
		}
		return nil, err
	}

	return transaction, nil
}

func (t *TransactionRepository) FindByID(id string) (*domain.Transaction, error) {
	transaction := &domain.Transaction{}
	err := t.db.QueryRow(`
		SELECT id, idempotency_key, description, created_at FROM transactions WHERE id = $1`,
		id).Scan(
		&transaction.ID,
		&transaction.IdempotencyKey,
		&transaction.Description,
		&transaction.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, TransactionNotFound
		}
		return nil, err
	}

	return transaction, nil
}
