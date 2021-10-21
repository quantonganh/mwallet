package account

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
	accountmocks "github.com/quantonganh/mwallet/account/mocks"
)

func TestMakeHandler(t *testing.T) {
	t.Run("add account", testAddAccount)
	t.Run("get account", testGetAccount)
	t.Run("list accounts", testListAccounts)
	t.Run("delete account", testDeleteAccount)
}

func testAddAccount(t *testing.T) {
	as := new(accountmocks.Service)
	bob := mwallet.Account{
		ID:       "bob123",
		Balance:  100.00,
		Currency: "USD",
	}
	body, err := json.Marshal(bob)
	require.NoError(t, err)
	as.On("AddAccount", &bob).Return(nil)

	accountHandler := MakeHandler(as, log.NewNopLogger())

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/opening/accounts", bytes.NewBuffer(body))
	accountHandler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)

	var resp mwallet.Account
	require.NoError(t, json.NewDecoder(rr.Body).Decode(&resp))
	assert.Equal(t, bob.ID, resp.ID)
	assert.Equal(t, bob.Balance, resp.Balance)
	assert.Equal(t, bob.Currency, resp.Currency)
}

func testGetAccount(t *testing.T) {
	as := new(accountmocks.Service)
	bob := mwallet.Account{
		ID:       "bob123",
		Balance:  100.00,
		Currency: "USD",
	}
	as.On("GetAccount", "bob123").Return(&bob, nil)

	accountHandler := MakeHandler(as, log.NewNopLogger())

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/opening/accounts/bob123", nil)
	accountHandler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)

	var resp mwallet.Account
	require.NoError(t, json.NewDecoder(rr.Body).Decode(&resp))
	assert.Equal(t, bob.ID, resp.ID)
	assert.Equal(t, bob.Balance, resp.Balance)
	assert.Equal(t, bob.Currency, resp.Currency)
}

func testListAccounts(t *testing.T) {
	as := new(accountmocks.Service)
	bob := mwallet.Account{
		ID:       "bob123",
		Balance:  100.00,
		Currency: "USD",
	}
	alice := mwallet.Account{
		ID:       "alice",
		Balance:  0.01,
		Currency: "USD",
	}
	as.On("ListAccounts").Return([]*mwallet.Account{&bob, &alice}, nil)

	accountHandler := MakeHandler(as, log.NewNopLogger())

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/opening/accounts", nil)
	accountHandler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)

	var resp listAccountsResponse
	require.NoError(t, json.NewDecoder(rr.Body).Decode(&resp))
	require.NoError(t, resp.Err)
	assert.Equal(t, 2, len(resp.Accounts))
	assert.Equal(t, bob.ID, resp.Accounts[0].ID)
	assert.Equal(t, bob.Balance, resp.Accounts[0].Balance)
	assert.Equal(t, bob.Currency, resp.Accounts[0].Currency)
	assert.Equal(t, alice.ID, resp.Accounts[1].ID)
	assert.Equal(t, alice.Balance, resp.Accounts[1].Balance)
	assert.Equal(t, alice.Currency, resp.Accounts[1].Currency)
}

func testDeleteAccount(t *testing.T) {
	as := new(accountmocks.Service)
	as.On("DeleteAccount", "bob123").Return(nil)

	accountHandler := MakeHandler(as, log.NewNopLogger())

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/opening/accounts/bob123", nil)
	accountHandler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
}

