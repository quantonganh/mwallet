package postgresql

import (
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
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	bob := &mwallet.Account{
		ID:       "bob123",
		Balance:  100.00,
		Currency: "USD",
	}
	rows := sqlmock.NewRows([]string{"id", "balance", "currency"}).AddRow(bob.ID, bob.Balance, bob.Currency)
	mock.ExpectQuery(regexp.QuoteMeta(sqlSelectForUpdateAccount)).WithArgs(bob.ID).WillReturnRows(rows)

	alice := &mwallet.Account{
		ID:       "alice456",
		Balance:  0.01,
		Currency: "USD",
	}
	rows = sqlmock.NewRows([]string{"id", "balance", "currency"}).AddRow(alice.ID, alice.Balance, alice.Currency)
	mock.ExpectQuery(regexp.QuoteMeta(sqlSelectForUpdateAccount)).WithArgs(alice.ID).WillReturnRows(rows)

	mock.ExpectBegin()

	mock.ExpectExec(regexp.QuoteMeta(sqlUpdateAccount)).WithArgs(50.00, bob.ID).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec(regexp.QuoteMeta(sqlUpdateAccount)).WithArgs(50.01, alice.ID).WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectExec(regexp.QuoteMeta(sqlInsertPayment)).WithArgs(bob.ID, alice.ID, 50.00).WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit()

	accountRepo := NewAccountRepository(db)
	err = accountRepo.Transfer("bob123", "alice456", 50.00)
	require.NoError(t, err)

	rows = sqlmock.NewRows([]string{"id", "balance", "currency"}).AddRow(bob.ID, 50.00, bob.Currency)
	mock.ExpectQuery(regexp.QuoteMeta(sqlSelectAccount)).WithArgs(bob.ID).WillReturnRows(rows)

	newBob, err := accountRepo.Find(bob.ID)
	require.NoError(t, err)
	assert.Equal(t, 50.00, newBob.Balance)

	rows = sqlmock.NewRows([]string{"id", "balance", "currency"}).AddRow(alice.ID, 50.01, alice.Currency)
	mock.ExpectQuery(regexp.QuoteMeta(sqlSelectAccount)).WithArgs(alice.ID).WillReturnRows(rows)

	newAlice, err := accountRepo.Find(alice.ID)
	require.NoError(t, err)
	assert.Equal(t, 50.01, newAlice.Balance)
}

func TestCreatePayment(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	payment := &mwallet.Payment{
		FromAccount: "bob123",
		ToAccount:   "alice456",
		Amount:      50.00,
	}
	mock.ExpectExec(regexp.QuoteMeta(sqlInsertPayment)).WithArgs(payment.FromAccount, payment.ToAccount, payment.Amount).WillReturnResult(sqlmock.NewResult(1, 1))

	paymentRepo := NewPaymentRepository(db)
	require.NoError(t, paymentRepo.Create(payment))
}

func TestFindPayment(t *testing.T) {
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
	mock.ExpectQuery(regexp.QuoteMeta(sqlSelectPaymentsByAccount)).WithArgs("bob123").WillReturnRows(rows)

	paymentRepo := NewPaymentRepository(db)
	payments, err := paymentRepo.Find("bob123")
	require.NoError(t, err)
	assert.Equal(t, "bob123", payments[0].FromAccount)
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