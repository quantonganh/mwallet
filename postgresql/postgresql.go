package postgresql

import (
	"context"
	"database/sql"

	_ "github.com/lib/pq"
	"github.com/pkg/errors"

	"github.com/quantonganh/mwallet"
	"github.com/quantonganh/mwallet/account"
	"github.com/quantonganh/mwallet/payment"
)

const (
	paymentDirectionIncoming = "incoming"
	paymentDirectionOutgoing = "outgoing"
)

var (
	errNegativeAmount = errors.New("Cannot transfer negative amount")
	errDifferentCurrencies = errors.New("Cannot transfer between different currencies")
	errNotEnoughMoney = errors.New("Not enough money on your account")
)

const (
	sqlInsertAccount = `INSERT INTO account (id, balance, currency) VALUES ($1, $2, $3)`
	sqlSelectAccount = `SELECT "id", "balance", "currency" FROM account WHERE id=$1`
	sqlSelectForUpdateAccount = `SELECT "id", "balance", "currency" FROM account WHERE id=$1 FOR UPDATE`
	sqlUpdateAccount     = `UPDATE account SET balance=$1 WHERE id=$2`
	sqlSelectAllAccounts = `SELECT "id", "balance", "currency" from account`
	sqlDeleteAccount = `DELETE FROM account WHERE id=$1`

	sqlInsertPayment = `INSERT INTO payment(account, amount, from_account, to_account, direction) VALUES ($1, $2, $3, $4, $5)`
	sqlSelectPaymentsByAccount = `SELECT id, account, amount, from_account, to_account, direction FROM payment WHERE account=$1`
	sqlSelectAllPayments = `SELECT id, account, amount, from_account, to_account FROM payment`
)

type accountRepository struct {
	db *sql.DB
}

func NewAccountRepository(db *sql.DB) account.Repository {
	return &accountRepository{
		db: db,
	}
}

// Create creates new account
func (r *accountRepository) Create(a *mwallet.Account) error {
	_, err := r.db.Exec(sqlInsertAccount, a.ID, a.Balance, a.Currency)
	if err != nil {
		return errors.Wrapf(err, "failed to add account: %s", a.ID)
	}
	return nil
}

// Find finds account base on ID
func (r *accountRepository) Find(id string) (*mwallet.Account, error) {
	var (
		balance float64
		currency string
	)
	err := r.db.QueryRow(sqlSelectAccount, id).Scan(&id, &balance, &currency)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get account: %s", id)
	}

	return &mwallet.Account{
		ID:       id,
		Balance:  balance,
		Currency: currency,
	}, nil
}

func (r *accountRepository) getForUpdate(id string) (*mwallet.Account, error) {
	var (
		balance float64
		currency string
	)
	err := r.db.QueryRow(sqlSelectForUpdateAccount, id).Scan(&id, &balance, &currency)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get account: %s", id)
	}

	return &mwallet.Account{
		ID:       id,
		Balance:  balance,
		Currency: currency,
	}, nil
}

// Transfer send a payment from one account to another
func (r *accountRepository) Transfer(fromAccountID, toAccountID string, amount float64) error {
	if amount < 0 {
		return errNegativeAmount
	}

	fromAccount, err := r.getForUpdate(fromAccountID)
	if err != nil {
		return err
	}
	if fromAccount.Balance < amount {
		return errNotEnoughMoney
	}

	toAccount, err := r.getForUpdate(toAccountID)
	if err != nil {
		return err
	}

	if fromAccount.Currency != toAccount.Currency {
		return errDifferentCurrencies
	}

	ctx := context.Background()
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, sqlUpdateAccount, fromAccount.Balance - amount, fromAccountID)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	_, err = tx.ExecContext(ctx, sqlUpdateAccount, toAccount.Balance + amount, toAccountID)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	_, err = r.db.Exec(sqlInsertPayment, fromAccountID, amount, sql.NullString{}, toAccountID, paymentDirectionOutgoing)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	_, err = r.db.Exec(sqlInsertPayment, toAccountID, amount, fromAccountID, sql.NullString{}, paymentDirectionIncoming)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

// FindAll finds all accounts
func (r *accountRepository) FindAll() ([]*mwallet.Account, error) {
	rows, err := r.db.Query(sqlSelectAllAccounts)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get all accounts")
	}
	defer rows.Close()

	var accounts []*mwallet.Account
	for rows.Next() {
		var (
			id string
			balance float64
			currency string
		)
		err := rows.Scan(&id, &balance, &currency)
		if err != nil {
			return nil, errors.Wrap(err, "failed to scan rows")
		}
		accounts = append(accounts, &mwallet.Account{
			ID:       id,
			Balance:  balance,
			Currency: currency,
		})
	}

	return accounts, nil
}

// Delete deletes an account
func (r *accountRepository) Delete(id string) error {
	_, err := r.db.Exec(sqlDeleteAccount, id)
	if err != nil {
		return errors.Wrapf(err, "failed to delete account: %s", id)
	}
	return nil
}

type paymentRepository struct {
	db *sql.DB
}

// NewPaymentRepository creates new payment repository
func NewPaymentRepository(db *sql.DB) payment.Repository {
	return &paymentRepository{
		db: db,
	}
}

// Find finds payments relate to an account
func (r *paymentRepository) Find(accountID string) ([]*mwallet.Payment, error) {
	rows, err := r.db.Query(sqlSelectPaymentsByAccount, accountID)
	if err != nil {
		return nil, err
	}

	var payments []*mwallet.Payment
	for rows.Next() {
		var (
			id string
			account string
			amount float64
			fromAccount sql.NullString
			toAccount sql.NullString
			direction string
		)
		if err := rows.Scan(&id, &account, &amount, &fromAccount, &toAccount, &direction); err != nil {
			return nil, err
		}
		payments = append(payments, &mwallet.Payment{
			ID:          id,
			Account: account,
			Amount:      amount,
			FromAccount: fromAccount.String,
			ToAccount:   toAccount.String,
			Direction: direction,
		})
	}
	return payments, nil
}

// FindAll finds all payments
func (r *paymentRepository) FindAll() ([]*mwallet.Payment, error) {
	rows, err := r.db.Query(sqlSelectAllPayments)
	if err != nil {
		return nil, err
	}

	var payments []*mwallet.Payment
	for rows.Next() {
		var (
			id string
			fromAccount string
			toAccount string
			amount float64
		)
		if err := rows.Scan(&id, &fromAccount, &toAccount, &amount); err != nil {
			return nil, err
		}
		payments = append(payments, &mwallet.Payment{
			ID:          id,
			FromAccount: fromAccount,
			ToAccount:   toAccount,
			Amount:      amount,
		})
	}
	return payments, nil
}
