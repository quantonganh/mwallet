package payment

import (
	"github.com/quantonganh/mwallet"
)

// Repository is the interface that wraps the methods related to a payment
type Repository interface {
	Find(accountID string) ([]*mwallet.Payment, error)
	FindAll() ([]*mwallet.Payment, error)
}
