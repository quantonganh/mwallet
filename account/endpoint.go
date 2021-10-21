package account

import (
	"context"

	"github.com/go-kit/kit/endpoint"

	"github.com/quantonganh/mwallet"
)

// swagger:model addAccount
type addAccount struct {
	// required: true
	// example: bob123
	ID string
	// required: true
	// example: 100.00
	Balance float64
	// required: true
	// example: USD
	Currency string
}

func makeAddAccountEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(addAccount)
		if err := s.AddAccount(&mwallet.Account{
			ID:       req.ID,
			Balance:  req.Balance,
			Currency: req.Currency,
		}); err != nil {
			return nil, err
		}
		return req, nil
	}
}

// swagger:model getAccountRequest
type getAccountRequest struct {
	// required: true
	// example: bob123
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

// swagger:model listAccountsResponse
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

// swagger:model deleteAccountRequest
type deleteAccountRequest struct {
	// required: true
	// example: bob123
	ID string
}

// swagger:model deleteAccountResponse
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

