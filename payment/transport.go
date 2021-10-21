package payment

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-kit/kit/transport"
	kithttp "github.com/go-kit/kit/transport/http"
	kitlog "github.com/go-kit/log"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"

	"github.com/quantonganh/mwallet"
)

var errBadRoute = errors.New("bad route")

// MakeHandler creates HTTP handler for payment
func MakeHandler(ps Service, logger kitlog.Logger) http.Handler {
	opts := []kithttp.ServerOption{
		kithttp.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
	}

	sendPaymentHandler := kithttp.NewServer(
		makeSendPaymentEndpoint(ps),
		decodeSendPaymentRequest,
		encodeResponse,
		opts...
	)

	findPaymentHandler := kithttp.NewServer(
		makeFindPaymentEndpoint(ps),
		decodeFindPaymentRequest,
		encodeResponse,
		opts...
	)

	listPaymentsHandler := kithttp.NewServer(
		makeListPaymentsEndpoint(ps),
		decodeListRequest,
		encodeResponse,
		opts...
	)

	r := mux.NewRouter()
	// swagger:operation POST /transferring/payments payments sendPayment
	// Transfer payment.
	// ---
	// produces:
	// - application/json
	// parameters:
	// - name: transfer payment
	//   in: body
	//   description: payment payload
	//   schema:
	//     "$ref": "#/definitions/sendPayment"
	//   required: true
	// responses:
	//     '200':
	//         schema:
	//           "$ref": "#/definitions/sendPayment"
	//     '500':
	//         description: Internal server error
	r.Handle("/transferring/payments", sendPaymentHandler).Methods(http.MethodPost)
	// swagger:operation GET /transferring/payments/{id} payments findPayment
	// Find payment.
	// ---
	// produces:
	// - application/json
	// parameters:
	// - name: id
	//   in: path
	//   description: id of the account
	//   type: string
	//   required: true
	// responses:
	//     '200':
	//         schema:
	//           "$ref": "#/definitions/findPaymentResponse"
	//     '500':
	//         description: Internal server error
	r.Handle("/transferring/payments/{id}", findPaymentHandler).Methods(http.MethodGet)
	// swagger:operation GET /transferring/payments payments listPayments
	// List all payments.
	// ---
	// produces:
	// - application/json
	// responses:
	//     '200':
	//         schema:
	//           "$ref": "#/definitions/listPaymentResponse"
	//     '500':
	//         description: Internal server error
	r.Handle("/transferring/payments", listPaymentsHandler).Methods(http.MethodGet)

	return r
}

func decodeSendPaymentRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var body mwallet.Payment
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return nil, err
	}

	return sendPayment{
		FromAccount: body.FromAccount,
		ToAccount:   body.ToAccount,
		Amount:      body.Amount,
	}, nil
}

func decodeFindPaymentRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		return nil, errBadRoute
	}
	return findPaymentRequest{
		AccountID: id,
	}, nil
}

func decodeListRequest(_ context.Context, r *http.Request) (interface{}, error) {
	return listPaymentRequest{}, nil
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(errorer); ok && e.error() != nil {
		encodeError(ctx, e.error(), w)
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

type errorer interface {
	error() error
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	switch err {
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}
