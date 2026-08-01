package postgres

import "database/sql"

type BalanceRepository struct{ db *sql.DB }

func NewPostgresSQLBalanceRepository(db *sql.DB) *BalanceRepository {
	return &BalanceRepository{db: db}
}
