package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type AccountType string

const (
	Asset     AccountType = "ASSET"
	Liability AccountType = "LIABILITY"
	Equity    AccountType = "EQUITY"
	Revenue   AccountType = "REVENUE"
	Expense   AccountType = "EXPENSE"
)

type Account struct {
	ID   uuid.UUID   `json:"id"`
	Name string      `json:"name"`
	Type AccountType `json:"type"`
	// Balance   decimal.Decimal `json:"balance"`
	CreatedAt time.Time `json:"created_at"`
}

type LegDirection string

const (
	Debit  LegDirection = "DEBIT"
	Credit LegDirection = "CREDIT"
)

type TransactionLeg struct {
	ID        uuid.UUID       `json:"id"`
	AccountID uuid.UUID       `json:"account_id"`
	Direction LegDirection    `json:"direction"`
	Amount    decimal.Decimal `json:"amount"`
}

type Transaction struct {
	ID             uuid.UUID        `json:"id"`
	IdempotencyKey string           `json:"idempotency_key"`
	Description    string           `json:"description"`
	Legs           []TransactionLeg `json:"legs"`
	CreatedAt      time.Time        `json:"created_at"`
}

func (t *Transaction) Validate() error {

	if len(t.Legs) < 2 {
		return errors.New("Uma transação deve conter pelo menos dois lados 'conta' (débito e crédito)")
	}

	totalDebit := decimal.Zero
	totalCredit := decimal.Zero

	for _, leg := range t.Legs {
		if leg.Amount.LessThanOrEqual(decimal.Zero) {
			return errors.New("O valor de cada lado deve ser maior que zero (0)")
		}

		if leg.Direction == Debit {
			totalDebit = totalDebit.Add(leg.Amount)
		} else if leg.Direction == Credit {
			totalCredit = totalCredit.Add(leg.Amount)
		}
	}

	if !totalDebit.Equal(totalCredit) {
		return errors.New("Livro razao desbalanceado: a soma dos debitos deve ser igual à soma dos creditos")
	}

	return nil
}

func (a *Account) ValidateAccount() error {

	if a.Name == "" {
		return errors.New("Nome não pode estar vazio")
	}

	return nil
}

func (a *Account) ValidateAccountType(types string) error {
	switch a.Type {
	case Asset, Equity, Liability, Revenue, Expense:
		return nil
	default:
		return errors.New("Tipo de conta inválido")
	}

}
