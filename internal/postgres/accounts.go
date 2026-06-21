package postgres

import (
	"context"
	"errors"

	"github.com/ABAlmeida/pennyledger/internal/accounts"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AccountStore struct {
	pool *pgxpool.Pool
}

func NewAccountStore(pool *pgxpool.Pool) *AccountStore {
	return &AccountStore{
		pool: pool,
	}
}

func (s *AccountStore) Create(ctx context.Context, id string, ownerName string) (accounts.Account, error) {
	var account accounts.Account

	err := s.pool.QueryRow(ctx, `
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

	return account, err
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
