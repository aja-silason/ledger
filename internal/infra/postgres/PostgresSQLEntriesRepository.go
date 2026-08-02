package postgres

import (
	"database/sql"
	"errors"
	"time"

	"github.com/aja-silason/ledger/internal/domain"
	"github.com/google/uuid"
)

type EntriesRepository struct{ db *sql.DB }

func NewPostgresSQLEntriesRepository(db *sql.DB) *EntriesRepository {
	return &EntriesRepository{db: db}
}

var EntriesNotFound = errors.New("Entrada não encontrada")

func (t *EntriesRepository) Save(entries *domain.TransactionLeg) (*domain.TransactionLeg, error) {
	id := uuid.New()
	now := time.Now().UTC()

	res, err := t.db.Exec(`
		INSERT INTO transactions (id, transaction_id, account_id, direction, amount, created_at) VALUES ($1, $2, $3, $4, $5, $6)`,
		id,
		entries.TransactionID,
		entries.AccountID,
		entries.Direction,
		entries.Amount,
		now)
	if res != nil {
		return nil, err
	}

	return entries, nil
}

// func (t *TransactionRepository) FindByIdempotencyKey(idempotencyKey string) (*domain.Transaction, error) {
// 	transaction := &domain.Transaction{}
// 	err := t.db.QueryRow(`
// 		SELECT id, idempotency_key, description, created_at FROM transactions WHERE idempotency_key = $1`,
// 		idempotencyKey).Scan(
// 		&transaction.ID,
// 		&transaction.IdempotencyKey,
// 		&transaction.Description,
// 		&transaction.CreatedAt)
// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			return nil, TransactionNotFound
// 		}
// 		return nil, err
// 	}

// 	return transaction, nil
// }

// func (t *TransactionRepository) FindByID(id string) (*domain.Transaction, error) {
// 	transaction := &domain.Transaction{}
// 	err := t.db.QueryRow(`
// 		SELECT id, idempotency_key, description, created_at FROM transactions WHERE id = $1`,
// 		id).Scan(
// 		&transaction.ID,
// 		&transaction.IdempotencyKey,
// 		&transaction.Description,
// 		&transaction.CreatedAt)
// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			return nil, TransactionNotFound
// 		}
// 		return nil, err
// 	}

// 	return transaction, nil
// }
