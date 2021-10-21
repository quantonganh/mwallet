package account

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

// MakeHandler creates new HTTP handler for account
func MakeHandler(as Service, logger kitlog.Logger) http.Handler {
	opts := []kithttp.ServerOption{
		kithttp.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
	}

	addAccountHandler := kithttp.NewServer(
		makeAddAccountEndpoint(as),
		decodeAddAccountRequest,
		encodeResponse,
		opts...
	)

	findAccountHandler := kithttp.NewServer(
		makeFindAccountEndpoint(as),
		decodeFindAccountRequest,
		encodeResponse,
		opts...
	)

	listAccountsHandler := kithttp.NewServer(
		makeListAccountsEndpoint(as),
		decodeListAccountsRequest,
		encodeResponse,
		opts...
	)

	deleteAccountHandler := kithttp.NewServer(
		makeDeleteAccountEndpoint(as),
		decodeDeleteAccountRequest,
		encodeResponse,
		opts...
	)

	r := mux.NewRouter()
	// swagger:operation POST /opening/accounts accounts createAccount
	// Create new account.
	// ---
	// produces:
	// - application/json
	// parameters:
	// - name: create account
	//   in: body
	//   description: account payload
	//   schema:
	//     "$ref": "#/definitions/addAccount"
	//   required: true
	// responses:
	//     '200':
	//         description: Created
	//         schema:
	//           "$ref": "#/definitions/addAccount"
	//     '500':
	//         description: Internal server error
	r.Handle("/opening/accounts", addAccountHandler).Methods(http.MethodPost)
	// swagger:operation GET /opening/accounts/{id} accounts findAccount
	// Find account.
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
	//           "$ref": "#/definitions/addAccount"
	//     '500':
	//         description: Internal server error
	r.Handle("/opening/accounts/{id}", findAccountHandler).Methods(http.MethodGet)
	// swagger:operation GET /opening/accounts accounts listAccounts
	// List all accounts.
	// ---
	// produces:
	// - application/json
	// responses:
	//     '200':
	//         schema:
	//           "$ref": "#/definitions/listAccountsResponse"
	//     '500':
	//         description: Internal server error
	r.Handle("/opening/accounts", listAccountsHandler).Methods(http.MethodGet)
	// swagger:operation DELETE /opening/accounts/{id} accounts deleteAccount
	// Delete an account.
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
	//           "$ref": "#/definitions/deleteAccountResponse"
	//     '500':
	//         description: Internal server error
	r.Handle("/opening/accounts/{id}", deleteAccountHandler).Methods(http.MethodDelete)

	return r
}

func decodeAddAccountRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var body mwallet.Account
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return nil, err
	}

	return addAccount{
		ID:       body.ID,
		Balance:  body.Balance,
		Currency: body.Currency,
	}, nil
}

func decodeFindAccountRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		return nil, errBadRoute
	}
	return getAccountRequest{
		ID: id,
	}, nil
}

func decodeListAccountsRequest(_ context.Context, r *http.Request) (interface{}, error) {
	return listAccountsRequest{}, nil
}

func decodeDeleteAccountRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		return nil, errBadRoute
	}
	return deleteAccountRequest{
		ID: id,
	}, nil
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
