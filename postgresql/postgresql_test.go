package postgresql

import (
	"database/sql"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/quantonganh/mwallet"
)

func TestCreateAccount(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	account := &mwallet.Account{
		ID:       "bob123",
		Balance:  100.00,
		Currency: "USD",
	}
	mock.ExpectExec(regexp.QuoteMeta(sqlInsertAccount)).WithArgs(account.ID, account.Balance, account.Currency).WillReturnResult(sqlmock.NewResult(1, 1))

	accountRepo := NewAccountRepository(db)
	require.NoError(t, accountRepo.Create(account))
}

func TestFindAccount(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	account := &mwallet.Account{
		ID:       "bob123",
		Balance:  100.00,
		Currency: "USD",
	}
	rows := sqlmock.NewRows([]string{"id", "balance", "currency"}).AddRow(account.ID, account.Balance, account.Currency)
	mock.ExpectQuery(regexp.QuoteMeta(sqlSelectAccount)).WithArgs(account.ID).WillReturnRows(rows)

	accountRepo := NewAccountRepository(db)
	a, err := accountRepo.Find("bob123")
	require.NoError(t, err)
assert.Equal(t, 100.00, a.Balance)
	assert.Equal(t, "USD", a.Currency)
}

func TestFindAllAccounts(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	account := &mwallet.Account{
		ID:       "bob123",
		Balance:  100.00,
		Currency: "USD",
	}
	rows := sqlmock.NewRows([]string{"id", "balance", "currency"}).AddRow(account.ID, account.Balance, account.Currency)
	mock.ExpectQuery(regexp.QuoteMeta(sqlSelectAllAccounts)).WithArgs().WillReturnRows(rows)

	accountRepo := NewAccountRepository(db)
	accounts, err := accountRepo.FindAll()
	require.NoError(t, err)
	assert.Equal(t, 1, len(accounts))
	assert.Equal(t, 100.00, accounts[0].Balance)
	assert.Equal(t, "USD", accounts[0].Currency)
}

func TestDeleteAccount(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	account := &mwallet.Account{
		ID: "bob123",
	}
	mock.ExpectExec(regexp.QuoteMeta(sqlDeleteAccount)).WithArgs(account.ID).WillReturnResult(sqlmock.NewResult(1, 1))

	accountRepo := NewAccountRepository(db)
	require.NoError(t, accountRepo.Delete("bob123"))
}

func TestTransfer(t *testing.T) {
	testCases := map[string]struct{
		fromAccount mwallet.Account
		toAccount mwallet.Account
		amount float64
		err error
	}{
		"negative amount": {
			amount:      -10.00,
			err: errNegativeAmount,
		},
		"not enough money": {
			fromAccount: mwallet.Account{
				ID:       "bob123",
				Balance:  100.00,
				Currency: "USD",
			},
			toAccount:   mwallet.Account{
				ID:       "alice456",
				Balance:  0.01,
				Currency: "USD",
			},
			amount: 200.00,
			err: errNotEnoughMoney,
		},
		"different currencies": {
			fromAccount: mwallet.Account{
				ID:       "bob123",
				Balance:  100.00,
				Currency: "USD",
			},
			toAccount:   mwallet.Account{
				ID:       "alice456",
				Balance:  0.01,
				Currency: "EUR",
			},
			err: errDifferentCurrencies,
		},
		"happy case": {
			fromAccount: mwallet.Account{
				ID:       "bob123",
				Balance:  100.00,
				Currency: "USD",
			},
			toAccount:   mwallet.Account{
				ID:       "alice456",
				Balance:  0.01,
				Currency: "USD",
			},
			amount: 50.00,
			err: nil,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer db.Close()

			rows := sqlmock.NewRows([]string{"id", "balance", "currency"}).AddRow(tc.fromAccount.ID, tc.fromAccount.Balance, tc.fromAccount.Currency)
			mock.ExpectQuery(regexp.QuoteMeta(sqlSelectForUpdateAccount)).WithArgs(tc.fromAccount.ID).WillReturnRows(rows)

			rows = sqlmock.NewRows([]string{"id", "balance", "currency"}).AddRow(tc.toAccount.ID, tc.toAccount.Balance, tc.toAccount.Currency)
			mock.ExpectQuery(regexp.QuoteMeta(sqlSelectForUpdateAccount)).WithArgs(tc.toAccount.ID).WillReturnRows(rows)

			mock.ExpectBegin()

			mock.ExpectExec(regexp.QuoteMeta(sqlUpdateAccount)).WithArgs(tc.fromAccount.Balance-tc.amount, tc.fromAccount.ID).WillReturnResult(sqlmock.NewResult(1, 1))
			mock.ExpectExec(regexp.QuoteMeta(sqlUpdateAccount)).WithArgs(tc.toAccount.Balance+tc.amount, tc.toAccount.ID).WillReturnResult(sqlmock.NewResult(1, 1))

			mock.ExpectExec(regexp.QuoteMeta(sqlInsertPayment)).WithArgs(tc.fromAccount.ID, tc.amount, sql.NullString{}, tc.toAccount.ID, paymentDirectionOutgoing).WillReturnResult(sqlmock.NewResult(1, 1))
			mock.ExpectExec(regexp.QuoteMeta(sqlInsertPayment)).WithArgs(tc.toAccount.ID, tc.amount, tc.fromAccount.ID, sql.NullString{}, paymentDirectionIncoming).WillReturnResult(sqlmock.NewResult(1, 1))

			mock.ExpectCommit()

			accountRepo := NewAccountRepository(db)
			err = accountRepo.Transfer(tc.fromAccount.ID, tc.toAccount.ID, tc.amount)
			assert.Equal(t, tc.err, err)

			if err == nil {
				rows = sqlmock.NewRows([]string{"id", "balance", "currency"}).AddRow(tc.fromAccount.ID, tc.fromAccount.Balance-tc.amount, tc.fromAccount.Currency)
				mock.ExpectQuery(regexp.QuoteMeta(sqlSelectAccount)).WithArgs(tc.fromAccount.ID).WillReturnRows(rows)

				newBob, err := accountRepo.Find(tc.fromAccount.ID)
				require.NoError(t, err)
				assert.Equal(t, tc.fromAccount.Balance-tc.amount, newBob.Balance)

				rows = sqlmock.NewRows([]string{"id", "balance", "currency"}).AddRow(tc.toAccount.ID, tc.toAccount.Balance+tc.amount, tc.toAccount.Currency)
				mock.ExpectQuery(regexp.QuoteMeta(sqlSelectAccount)).WithArgs(tc.toAccount.ID).WillReturnRows(rows)

				newAlice, err := accountRepo.Find(tc.toAccount.ID)
				require.NoError(t, err)
				assert.Equal(t, tc.toAccount.Balance+tc.amount, newAlice.Balance)
			}
		})
	}
}

func TestFindPayment(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	payment := &mwallet.Payment{
		ID: "1",
		Account: "bob123",
		ToAccount:   "alice456",
		Amount:      50.00,
		Direction: paymentDirectionOutgoing,
	}
	rows := sqlmock.NewRows([]string{"id", "account", "amount", "from_account", "to_account", "direction"}).AddRow(payment.ID, payment.Account, payment.Amount, sql.NullString{}, payment.ToAccount, paymentDirectionOutgoing)
	mock.ExpectQuery(regexp.QuoteMeta(sqlSelectPaymentsByAccount)).WithArgs("bob123").WillReturnRows(rows)

	paymentRepo := NewPaymentRepository(db)
	payments, err := paymentRepo.Find("bob123")
	require.NoError(t, err)
	assert.Equal(t, "bob123", payments[0].Account)
	assert.Equal(t, "alice456", payments[0].ToAccount)
	assert.Equal(t, 50.00, payments[0].Amount)
}

func TestFindAllPayments(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	payment := &mwallet.Payment{
		ID: "1",
		FromAccount: "bob123",
		ToAccount:   "alice456",
		Amount:      50.00,
	}
	rows := sqlmock.NewRows([]string{"id", "from_account", "to_account", "amount"}).AddRow(payment.ID, payment.FromAccount, payment.ToAccount, payment.Amount)
	mock.ExpectQuery(regexp.QuoteMeta(sqlSelectAllPayments)).WithArgs().WillReturnRows(rows)

	paymentRepo := NewPaymentRepository(db)
	payments, err := paymentRepo.FindAll()
	require.NoError(t, err)
	assert.Equal(t, 1, len(payments))
	assert.Equal(t, "bob123", payments[0].FromAccount)
	assert.Equal(t, "alice456", payments[0].ToAccount)
	assert.Equal(t, 50.00, payments[0].Amount)
}