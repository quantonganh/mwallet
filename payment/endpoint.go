package payment

import (
	"context"

	"github.com/go-kit/kit/endpoint"

	"github.com/quantonganh/mwallet"
)

type sendRequest struct {
	fromAccount string
	toAccount string
	amount float64
}

type sendResponse struct {}

func makeSendPaymentEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(sendRequest)
		err := s.Send(req.fromAccount, req.toAccount, req.amount)
		if err != nil {
			return nil, err
		}
		return sendResponse{}, nil
	}
}

type findRequest struct {
	accountID string
}

type FindResponse struct {
	Payments []*mwallet.Payment `json:"payments"`
	Err error `json:"error"`
}

func makeFindPaymentEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(findRequest)
		payments, err := s.Find(req.accountID)
		return FindResponse{
			Payments: payments,
			Err: err,
		}, nil
	}
}

type listRequest struct {}

type listResponse struct {
	Payments []*mwallet.Payment `json:"payments"`
	Err error `json:"error"`
}

func makeListPaymentsEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		_ = request.(listRequest)
		payments, err := s.List()
		return listResponse{
			Payments: payments,
			Err:      err,
		}, nil
	}
}
