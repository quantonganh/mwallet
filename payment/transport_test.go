package payment

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-kit/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/quantonganh/mwallet"
	paymentmocks "github.com/quantonganh/mwallet/payment/mocks"
)

func TestMakeHandler(t *testing.T) {
	t.Run("send payment", testSendPayment)
	t.Run("find payment", testFindPayment)
	t.Run("list payments", testListPayments)
}

func testSendPayment(t *testing.T) {
	ps := new(paymentmocks.Service)
	payment := sendPayment{
		FromAccount: "bob123",
		ToAccount:   "alice456",
		Amount:      50.00,
	}
	body, err := json.Marshal(payment)
	require.NoError(t, err)
	ps.On("Send", "bob123", "alice456", 50.00).Return(nil)

	paymentHandler := MakeHandler(ps, log.NewNopLogger())

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/transferring/payments", bytes.NewBuffer(body))
	paymentHandler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)

	var resp mwallet.Payment
	require.NoError(t, json.NewDecoder(rr.Body).Decode(&resp))
	assert.Equal(t, payment.FromAccount, resp.FromAccount)
	assert.Equal(t, payment.ToAccount, resp.ToAccount)
	assert.Equal(t, payment.Amount, resp.Amount)
}

func testFindPayment(t *testing.T) {
	ps := new(paymentmocks.Service)
	payment := mwallet.Payment{
		Account:     "bob123",
		ToAccount:   "alice456",
		Amount:      50.00,
		Direction:   "outgoing",
	}
	ps.On("Find", "bob123").Return([]*mwallet.Payment{&payment}, nil)

	paymentHandler := MakeHandler(ps, log.NewNopLogger())

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/transferring/payments/bob123", nil)
	paymentHandler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)

	var resp mwallet.Account
	require.NoError(t, json.NewDecoder(rr.Body).Decode(&resp))
}

func testListPayments(t *testing.T) {
	ps := new(paymentmocks.Service)
	bob := mwallet.Account{
		ID:       "bob123",
		Balance:  100.00,
		Currency: "USD",
	}
	alice := mwallet.Account{
		ID:       "alice456",
		Balance:  0.01,
		Currency: "USD",
	}
	payment := mwallet.Payment{
		FromAccount: bob.ID,
		ToAccount:   alice.ID,
		Amount:      50.00,
	}
	ps.On("List").Return([]*mwallet.Payment{&payment}, nil)

	paymentHandler := MakeHandler(ps, log.NewNopLogger())

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/transferring/payments", nil)
	paymentHandler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)

	var resp listPaymentResponse
	require.NoError(t, json.NewDecoder(rr.Body).Decode(&resp))
	require.NoError(t, resp.Err)
	assert.Equal(t, 1, len(resp.Payments))
	assert.Equal(t, payment.FromAccount, resp.Payments[0].FromAccount)
	assert.Equal(t, payment.ToAccount, resp.Payments[0].ToAccount)
	assert.Equal(t, payment.Amount, resp.Payments[0].Amount)
}