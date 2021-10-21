package mwallet

// Account represents a mobile wallet account
type Account struct {
	ID string `json:"id"`
	Balance float64 `json:"balance"`
	Currency string `json:"currency"`
}

