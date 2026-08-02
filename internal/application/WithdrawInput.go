package application

type WithDrawInput struct {
	AccountID  string `json:"accountId"`
	Amount     int64  `json:"amount"`
	SecretCode int64  `json:"secretCode"`
}
