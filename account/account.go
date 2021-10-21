package account

import (
	"github.com/quantonganh/mwallet"
)

// Repository is the interface that wraps the methods related to an account
type Repository interface {
	Create(account *mwallet.Account) error
	Find(id string) (*mwallet.Account, error)
	Transfer(fromAccountID, toAccountID string, amount float64) error
	FindAll() ([]*mwallet.Account, error)
	Delete(id string) error
}

