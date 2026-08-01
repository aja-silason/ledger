package postgres

import (
	"database/sql"
	"errors"
	"time"

	"github.com/aja-silason/ledger/internal/domain"
	"github.com/google/uuid"
)

type TransactionRepository struct{ db *sql.DB }

func NewPostgresSQLTransactionRepository(db *sql.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

var TransactionNotFound = errors.New("Transação não encontrada")

func (t *TransactionRepository) Save(transaction *domain.Transaction) (*domain.Transaction, error) {
	id := uuid.New()
	now := time.Now().UTC()

	res, err := t.db.Exec(`
		INSERT INTO transactions (id, idempotency_key, description, created_at) VALUES ($1, $2, $3, $4)`,
		id,
		transaction.IdempotencyKey,
		transaction.Description,
		now)
	if res != nil {
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
