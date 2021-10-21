package account

import (
	"github.com/quantonganh/mwallet"
)

// Service is the interface that wraps the account methods
type Service interface {
	AddAccount(account *mwallet.Account) error
	GetAccount(id string) (*mwallet.Account, error)
	ListAccounts() ([]*mwallet.Account, error)
	DeleteAccount(id string) error
}

type service struct {
	account Repository
}

// AddAccount adds an account
func (s *service) AddAccount(account *mwallet.Account) error {
	if err := s.account.Create(account); err != nil {
		return err
	}
	return nil
}

// GetAccount gets an account base on ID
func (s *service) GetAccount(id string) (*mwallet.Account, error) {
	account, err := s.account.Find(id)
	if err != nil {
		return nil, err
	}
	return account, nil
}

// ListAccounts lists all accounts
func (s *service) ListAccounts() ([]*mwallet.Account, error) {
	accounts, err := s.account.FindAll()
	if err != nil {
		return nil, err
	}
	return accounts, nil
}

// DeleteAccount deletes an account base on ID
func (s *service) DeleteAccount(id string) error {
	return s.account.Delete(id)
}

// NewService creates new account service
func NewService(account Repository) Service {
	return &service{
		account: account,
	}
}