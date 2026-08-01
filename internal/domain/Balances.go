package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	CurrencyTypeNotAllowedError = errors.New("Tipo de moeda não permitido")
	AmountCanotBeNegativeError  = errors.New("Montante não pode ser negativo")
)

type CurrencyType string

const (
	AOA CurrencyType = "AOA"
	USD CurrencyType = "USD"
	BRL CurrencyType = "BRL"
	EUR CurrencyType = "EUR"
)

type Balance struct {
	ID            uuid.UUID    `json:"id"`
	AccountID     uuid.UUID    `json:"account_id"`
	CurrencyCode  CurrencyType `json:"currency_code"`
	CurrentAmount int          `json:"current_amount"`
	UpdatedAt     time.Time    `json:"updated_at"`
	CreatedAt     time.Time    `json:"created_at"`
}

func (b *Balance) ValidateCurrency() error {
	switch b.CurrencyCode {
	case AOA, USD, BRL, EUR:
		return nil
	default:
		return CurrencyTypeNotAllowedError
	}
}

func (b *Balance) ValidateAmount() error {
	if b.CurrentAmount < 0 {
		return AmountCanotBeNegativeError
	}
	return nil
}
