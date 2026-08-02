package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type WidrawStatus string

var (
	AmountMustBeMoreThanFifthThousandError = errors.New("Valor para o levantamento deve ser maior que 500 AOA")
	CodeMustBeDiferentAndMoreThanZeroError = errors.New("Código do levantamento deve ser maior que zero")
	CodeMustHaveEightCharacteresError      = errors.New("A referência de levantamento tem que ter 8 carácteres")
	CodeMustHaveThreeCharacteresError      = errors.New("Código Sécreto tem que ter 3 digitos")
	ErrWithdrawNotPending                  = errors.New("apenas levantamentos pendentes podem alterar de estado")
	ErrInvalidCode                         = errors.New("código de levantamento inválido")
)

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
	CodeHash  string       `json:"code_hash"`
	Code      string       `json:"code"`
	ExpiresAt time.Time    `json:"expires_at"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
}

func (w *Withdraw) ValidateCode() error {
	codeLength := len(w.Code)

	if codeLength != 8 {
		return CodeMustHaveEightCharacteresError
	}

	return nil
}

func (w *Withdraw) ValidateAmountToWithDraw() error {
	if w.Amount < 50000 {
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

func (w *Withdraw) MarkAsExpired() error {
	if w.Status != PENDING {
		return ErrWithdrawNotPending
	}
	w.Status = EXPIRED
	w.UpdatedAt = time.Now()
	return nil
}
