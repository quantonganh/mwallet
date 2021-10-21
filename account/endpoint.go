package account

import (
	"context"

	"github.com/go-kit/kit/endpoint"

	"github.com/quantonganh/mwallet"
)

type addAccountRequest struct {
	id string
	balance float64
	currency string
}

type addAccountResponse struct {}

func makeAddAccountEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(addAccountRequest)
		if err := s.AddAccount(&mwallet.Account{
			ID:       req.id,
			Balance:  req.balance,
			Currency: req.currency,
		}); err != nil {
			return nil, err
		}
		return addAccountResponse{}, nil
	}
}

type getAccountRequest struct {
	ID string
}

func makeFindAccountEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(getAccountRequest)
		account, err := s.GetAccount(req.ID)
		if err != nil {
			return nil, err
		}
		return account, nil
	}
}

type listAccountsRequest struct {}

type listAccountsResponse struct {
	Accounts []*mwallet.Account `json:"accounts"`
	Err error `json:"error"`
}

func makeListAccountsEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		_ = request.(listAccountsRequest)
		accounts, err := s.ListAccounts()
		return listAccountsResponse{
			Accounts: accounts,
			Err:      err,
		}, nil
	}
}

type deleteAccountRequest struct {
	ID string
}

type deleteAccountResponse struct {
	Err error `json:"error"`
}

func makeDeleteAccountEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(deleteAccountRequest)
		err := s.DeleteAccount(req.ID)
		if err != nil {
			return deleteAccountResponse{
				Err: err,
			}, err
		}
		return deleteAccountResponse{
			Err: nil,
		}, nil
	}
}

