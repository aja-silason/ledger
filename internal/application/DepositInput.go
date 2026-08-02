package application

type DepositInput struct {
	AccountId string `json:"accountId"`
	Amount    int64  `json:"amount"`
}
