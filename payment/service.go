package payment

import (
	"github.com/quantonganh/mwallet"
	"github.com/quantonganh/mwallet/account"
)

// Service is the interface that wraps the payment methods
type Service interface {
	Send(fromAccountID, toAccountID string, amount float64) error
	Find(accountID string) ([]*mwallet.Payment, error)
	List() ([]*mwallet.Payment, error)
}

type service struct {
	account account.Repository
	payment Repository
}

// Send sends a payment from one account to another
func (s *service) Send(fromAccountID, toAccountID string, amount float64) error {
	if err := s.account.Transfer(fromAccountID, toAccountID, amount); err != nil {
		return err
	}

	return nil
}

// Find finds a payment base on account ID
func (s *service) Find(accountID string) ([]*mwallet.Payment, error) {
	return s.payment.Find(accountID)
}

// List lists all payments
func (s *service) List() ([]*mwallet.Payment, error) {
	return s.payment.FindAll()
}

// NewService creates new payment service
func NewService(account account.Repository, payment Repository) Service {
	return &service{
		account: account,
		payment: payment,
	}
}