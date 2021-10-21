package mwallet

// Payment represents a transfer payment
type Payment struct {
	ID string `json:"id"`
	Account string `json:"account"`
	FromAccount string `json:"from_account"`
	ToAccount string `json:"to_account"`
	Amount float64 `json:"amount"`
	Direction string `json:"direction"`
}
