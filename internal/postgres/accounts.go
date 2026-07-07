package postgres

import (
	"context"
	"errors"

	"github.com/ABAlmeida/pennyledger/internal/accounts"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const openingBalanceEquityAccountID = "00000000-0000-0000-0000-000000000001"

type AccountStore struct {
	pool *pgxpool.Pool
}

func NewAccountStore(pool *pgxpool.Pool) *AccountStore {
	return &AccountStore{
		pool: pool,
	}
}

func (s *AccountStore) Create(ctx context.Context, id string, ownerName string, openingBalancePence int64) (accounts.Account, error) {
	var account accounts.Account

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return accounts.Account{}, err
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	err = tx.QueryRow(ctx, `
		INSERT INTO accounts (id, owner_name)
		VALUES ($1, $2)
		RETURNING id::text, owner_name, currency, balance_pence, status, created_at, updated_at
	`, id, ownerName).Scan(
		&account.ID,
		&account.OwnerName,
		&account.Currency,
		&account.BalancePence,
		&account.Status,
		&account.CreatedAt,
		&account.UpdatedAt,
	)
	if err != nil {
		return accounts.Account{}, err
	}

	if openingBalancePence > 0 {
		transactionID := uuid.NewString()
		customerEntryID := uuid.NewString()
		openingBalanceEntryID := uuid.NewString()
		var openingBalanceEquityBalance int64

		_, err = tx.Exec(ctx, `
			INSERT INTO ledger_transactions (id, kind)
			VALUES ($1, 'opening_balance')
		`, transactionID)
		if err != nil {
			return accounts.Account{}, err
		}

		err = tx.QueryRow(ctx, `
			UPDATE accounts
			SET balance_pence = balance_pence + $1,
			    updated_at = now()
			WHERE id = $2
			RETURNING balance_pence, updated_at
		`, openingBalancePence, id).Scan(&account.BalancePence, &account.UpdatedAt)
		if err != nil {
			return accounts.Account{}, err
		}

		err = tx.QueryRow(ctx, `
			UPDATE accounts
			SET balance_pence = balance_pence - $1,
			    updated_at = now()
			WHERE id = $2
			RETURNING balance_pence
		`, openingBalancePence, openingBalanceEquityAccountID).Scan(&openingBalanceEquityBalance)
		if err != nil {
			return accounts.Account{}, err
		}

		_, err = tx.Exec(ctx, `
			INSERT INTO ledger_entries (id, ledger_transaction_id, account_id, amount_pence, balance_after_pence)
			VALUES
				($1, $2, $3, $4, $5),
				($6, $2, $7, $8, $9)
		`,
			customerEntryID,
			transactionID,
			id,
			openingBalancePence,
			account.BalancePence,
			openingBalanceEntryID,
			openingBalanceEquityAccountID,
			-openingBalancePence,
			openingBalanceEquityBalance,
		)
		if err != nil {
			return accounts.Account{}, err
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return accounts.Account{}, err
	}

	return account, nil
}

func (s *AccountStore) GetByID(ctx context.Context, id string) (accounts.Account, error) {
	var account accounts.Account

	err := s.pool.QueryRow(ctx, `
		SELECT id::text, owner_name, currency, balance_pence, status, created_at, updated_at
		FROM accounts
		WHERE id = $1
	`, id).Scan(
		&account.ID,
		&account.OwnerName,
		&account.Currency,
		&account.BalancePence,
		&account.Status,
		&account.CreatedAt,
		&account.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return accounts.Account{}, accounts.ErrNotFound
	}

	return account, err
}
