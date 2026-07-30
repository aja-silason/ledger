package postgres

import (
	"database/sql"
	"errors"

	"github.com/aja-silason/ledger/internal/domain"
)

type AccountGateway struct{ db *sql.DB }

func NewPostgresSQLAccountGateway(db *sql.DB) *AccountGateway {
	return &AccountGateway{db: db}
}

func (r *AccountGateway) FindByID(id string) (*domain.Account, error) {
	u := &domain.Account{}
	err := r.db.QueryRow(
		`SELECT id, name, type, created_at FROM accounts WHERE id = $1`, id,
	).Scan(&u.ID, &u.Name, &u.Type, &u.CreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrAccountNotFound
	}
	return u, err
}
