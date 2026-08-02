package domain

type SuccessMessage interface {
	Message() string
}

type Success struct {
	Value string `json:"message"`
}

func (t Success) Message() string {
	return t.Value
}

func NewSuccessMessage(msg string) SuccessMessage {
	return Success{Value: msg}
}
