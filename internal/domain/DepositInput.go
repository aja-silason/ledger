package domain

import "errors"

type DepositInput struct {
	AccountId string `json:"accountId"`
	Amount    int64  `json:"amount"`
}

func (a *DepositInput) ValidateDepositInput(input *DepositInput) error {
	switch input {
	case input:
		return nil
	default:
		return errors.New("Precisa fornecer informações ao payload")
	}
}
