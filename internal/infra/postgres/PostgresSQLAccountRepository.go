package postgres

import (
	"database/sql"
	"errors"
	"time"

	"github.com/aja-silason/ledger/internal/domain"
	"github.com/google/uuid"
)

type AccountRepository struct{ db *sql.DB }

func NewPostgresSQLAccountRepository(db *sql.DB) *AccountRepository {
	return &AccountRepository{db: db}
}

var ErrAccountNotFound = errors.New("Conta não encontrada")

func (r *AccountRepository) Save(account *domain.Account) (*domain.Account, error) {
	id := uuid.New()
	now := time.Now().UTC()
	res, err := r.db.Exec(
		`INSERT INTO accounts (id, name, type, created_at) VALUES ($1, $2, $3, $4)`,
		id,
		account.Name,
		account.Type,
		now)

	if res != nil {
		return nil, err
	}

	return r.FindByID(id.String())
}

func (r *AccountRepository) FindByID(id string) (*domain.Account, error) {
	u := &domain.Account{}
	err := r.db.QueryRow(
		`SELECT id, name, type, created_at FROM accounts WHERE id = $1`, id,
	).Scan(&u.ID, &u.Name, &u.Type, &u.CreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrAccountNotFound
	}
	return u, err
}
