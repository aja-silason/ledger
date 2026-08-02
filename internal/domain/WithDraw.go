package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type WidrawStatus string

var AmountMustBeMoreThanFifthThousandError = errors.New("Valor para o levantamento deve ser maior que 500 AOA")

const (
	PENDING  WidrawStatus = "PENDING"
	DRAWED   WidrawStatus = "DRAWED"
	EXPIRED  WidrawStatus = "EXPIRED"
	CANCELED WidrawStatus = "CANCELED"
)

type Withdraw struct {
	ID        uuid.UUID    `json:"id"`
	AccountID uuid.UUID    `json:"account_id"`
	Amount    int64        `json:"amount"`
	Status    WidrawStatus `json:"status"`
	CodeHash  WidrawStatus `json:"code_hash"`
	ExpiresAt time.Time    `json:"expires_at"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
}

func (w *Withdraw) ValidateAmountToWithDraw() error {
	if w.Amount < 0 {
		return AmountMustBeMoreThanFifthThousandError
	}
	return nil
}

func (w *Withdraw) Drawed() (Withdraw, error) {
	drawed := Withdraw{Status: DRAWED}
	return drawed, nil
}

func (w *Withdraw) Canceled() (Withdraw, error) {
	canceled := Withdraw{Status: CANCELED}
	return canceled, nil
}

func (w *Withdraw) Expired() (Withdraw, error) {
	expired := Withdraw{Status: EXPIRED}
	return expired, nil
}
