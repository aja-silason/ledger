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
