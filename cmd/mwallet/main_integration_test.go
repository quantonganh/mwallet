// +build integration

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/quantonganh/mwallet"
	"github.com/quantonganh/mwallet/payment"
)

func TestTransferPayment(t *testing.T) {
	t.Run("add bob account", func(t *testing.T) {
		addAccount(t, "bob123", 100.00, "USD")
	})

	t.Run("add alice account", func(t *testing.T) {
		addAccount(t, "alice456", 0.01, "USD")
	})

	t.Run("transfer payment", func(t *testing.T) {
		transferPayment(t, "bob123", "alice456", 50.00)
	})

	bob := getAccount(t, "bob123")
	alice := getAccount(t, "alice456")
	assert.Equal(t, 50.00, bob.Balance)
	assert.Equal(t, 50.01, alice.Balance)

	payments := getPayments(t, "bob123")
	assert.Equal(t, 1, len(payments))
	assert.Equal(t, "bob123", payments[0].Account)
	assert.Equal(t, "alice456", payments[0].ToAccount)
	assert.Equal(t, 50.00, payments[0].Amount)

	t.Cleanup(func() {
		deleteAccount(t, "bob123")
		deleteAccount(t, "alice456")
	})
}

func addAccount(t *testing.T, id string, balance float64, currency string) {
	account := mwallet.Account{
		ID:       id,
		Balance:  balance,
		Currency: currency,
	}
	body, err := json.Marshal(account)
	require.NoError(t, err)

	req, err := http.NewRequest(http.MethodPost, "http://localhost:8080/opening/accounts", bytes.NewBuffer(body))
	require.NoError(t, err)

	client := http.Client{}
	resp, err := client.Do(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func transferPayment(t *testing.T, fromAccount, toAccount string, amount float64) {
	payment := mwallet.Payment{
		FromAccount: fromAccount,
		ToAccount:   toAccount,
		Amount:      amount,
	}
	body, err := json.Marshal(payment)
	require.NoError(t, err)

	req, err := http.NewRequest(http.MethodPost, "http://localhost:8080/transferring/payments", bytes.NewBuffer(body))
	require.NoError(t, err)

	client := http.Client{}
	resp, err := client.Do(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func getAccount(t *testing.T, id string) *mwallet.Account {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://localhost:8080/opening/accounts/%s", id), nil)
	require.NoError(t, err)

	client := http.Client{}
	resp, err := client.Do(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var account *mwallet.Account
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&account))

	return account
}

func deleteAccount(t *testing.T, id string) {
	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("http://localhost:8080/opening/accounts/%s", id), nil)
	require.NoError(t, err)

	client := http.Client{}
	resp, err := client.Do(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func getPayments(t *testing.T, accountID string) []*mwallet.Payment {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://localhost:8080/transferring/payments/%s", accountID), nil)
	require.NoError(t, err)

	client := http.Client{}
	resp, err := client.Do(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var findResp payment.FindPaymentResponse
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&findResp))
	require.NoError(t, findResp.Err)

	return findResp.Payments
}
