package application

type TransferMoneyInput struct {
	FromAccountID string `json:"fromAccountId"`
	ToAccountID   string `json:"toAccountId"`
	Amount        int64  `json:"amount"`
}
