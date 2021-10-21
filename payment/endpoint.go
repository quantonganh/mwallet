package payment

import (
	"context"

	"github.com/go-kit/kit/endpoint"

	"github.com/quantonganh/mwallet"
)

// swagger:model sendPayment
type sendPayment struct {
	// required: true
	// example: bob123
	FromAccount string `json:"from_account"`
	// required: true
	// example: alice456
	ToAccount string `json:"to_account"`
	// required: true
	// example: 50.00
	Amount    float64 `json:"amount"`
}

func makeSendPaymentEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(sendPayment)
		err := s.Send(req.FromAccount, req.ToAccount, req.Amount)
		if err != nil {
			return nil, err
		}
		return req, nil
	}
}

// swagger:model findPaymentRequest
type findPaymentRequest struct {
	// required: true
	// example: bob123
	AccountID string
}

// swagger:model findPaymentResponse
type FindPaymentResponse struct {
	Payments []*mwallet.Payment `json:"payments"`
	Err error `json:"error"`
}

func makeFindPaymentEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(findPaymentRequest)
		payments, err := s.Find(req.AccountID)
		return FindPaymentResponse{
			Payments: payments,
			Err: err,
		}, nil
	}
}

type listPaymentRequest struct {}

// swagger:model listPaymentResponse
type listPaymentResponse struct {
	Payments []*mwallet.Payment `json:"payments"`
	Err error `json:"error"`
}

func makeListPaymentsEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		_ = request.(listPaymentRequest)
		payments, err := s.List()
		return listPaymentResponse{
			Payments: payments,
			Err:      err,
		}, nil
	}
}
