package account

import (
	"github.com/quantonganh/mwallet"
)

type Service interface {
	AddAccount(account *mwallet.Account) error
	GetAccount(id string) (*mwallet.Account, error)
	ListAccounts() ([]*mwallet.Account, error)
}

type service struct {
	account Repository
}

func (s *service) AddAccount(account *mwallet.Account) error {
	if err := s.account.Create(account); err != nil {
		return err
	}
	return nil
}

func (s *service) GetAccount(id string) (*mwallet.Account, error) {
	account, err := s.account.Find(id)
	if err != nil {
		return nil, err
	}
	return account, nil
}

func (s *service) ListAccounts() ([]*mwallet.Account, error) {
	accounts, err := s.account.FindAll()
	if err != nil {
		return nil, err
	}
	return accounts, nil
}

func NewService(account Repository) Service {
	return &service{
		account: account,
	}
}