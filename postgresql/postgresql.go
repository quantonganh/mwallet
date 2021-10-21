package postgresql

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/pkg/errors"

	"github.com/quantonganh/mwallet"
	"github.com/quantonganh/mwallet/account"
	"github.com/quantonganh/mwallet/payment"
)

const (
	sqlInsertAccount = `INSERT INTO account (id, balance, currency) VALUES ($1, $2, $3)`
	sqlSelectAccount = `SELECT "id", "balance", "currency" FROM account WHERE id=$1`
	sqlSelectForUpdateAccount = `SELECT "id", "balance", "currency" FROM account WHERE id=$1 FOR UPDATE`
	sqlUpdateAccount     = `UPDATE account SET balance=$1 WHERE id=$2`
	sqlSelectAllAccounts = `SELECT "id", "balance", "currency" from account`
	sqlDeleteAccount = `DELETE FROM account WHERE id=$1`

	sqlInsertPayment = `INSERT INTO payment(from_account, to_account, amount) VALUES ($1, $2, $3)`
	sqlSelectPaymentsByAccount = `SELECT id, from_account, to_account, amount FROM payment WHERE from_account=$1 OR to_account=$1`
	sqlSelectAllPayments = `SELECT id, from_account, to_account, amount FROM payment`
)

type accountRepository struct {
	db *sql.DB
}

func NewAccountRepository(db *sql.DB) account.Repository {
	return &accountRepository{
		db: db,
	}
}

func (r *accountRepository) Create(a *mwallet.Account) error {
	_, err := r.db.Exec(sqlInsertAccount, a.ID, a.Balance, a.Currency)
	if err != nil {
		return errors.Wrapf(err, "failed to add account: %s", a.ID)
	}
	return nil
}

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

func (r *accountRepository) Transfer(fromAccountID, toAccountID string, amount float64) error {
	if amount < 0 {
		return errors.New("Amount must be positive.")
	}

	fromAccount, err := r.getForUpdate(fromAccountID)
	if err != nil {
		return err
	}
	if fromAccount.Balance < amount {
		return fmt.Errorf("%s do not have enough money to transfer %f", fromAccountID, amount)
	}

	toAccount, err := r.getForUpdate(toAccountID)
	if err != nil {
		return err
	}

	if fromAccount.Currency != toAccount.Currency {
		return errors.New("Only payments without the same currency are supported")
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

	_, err = r.db.Exec(sqlInsertPayment, fromAccountID, toAccountID, amount)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

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

func NewPaymentRepository(db *sql.DB) payment.Repository {
	return &paymentRepository{
		db: db,
	}
}

func (r *paymentRepository) Create(payment *mwallet.Payment) error {
	_, err := r.db.Exec(sqlInsertPayment, payment.FromAccount, payment.ToAccount, payment.Amount)
	if err != nil {
		return err
	}
	return nil
}

func (r *paymentRepository) Find(accountID string) ([]*mwallet.Payment, error) {
	rows, err := r.db.Query(sqlSelectPaymentsByAccount, accountID)
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
