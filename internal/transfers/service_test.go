package transfers

import (
	"context"
	"testing"
)

type fakeStore struct {
	createdID    string
	createdInput CreateTransferInput
	transfer     Transfer
	err          error
}

func (f *fakeStore) Create(ctx context.Context, id string, input CreateTransferInput) (Transfer, error) {
	f.createdID = id
	f.createdInput = input
	return f.transfer, f.err
}

func TestCreateTransferReturnsErrorWhenFromAccountIsMissing(t *testing.T) {
	store := fakeStore{}
	service := NewService(&store)
	_, err := service.CreateTransfer(t.Context(), CreateTransferInput{
		FromAccountID: "",
		ToAccountID:   "to-account-id",
		AmountPence:   100,
	})

	if err != ErrFromAccountRequired {
		t.Fatalf("expected error for missing from account, got nil")
	}
}

func TestCreateTransferReturnsErrorWhenToAccountIsMissing(t *testing.T) {
	store := fakeStore{}
	service := NewService(&store)
	_, err := service.CreateTransfer(t.Context(), CreateTransferInput{
		FromAccountID: "from-account-id",
		ToAccountID:   "",
		AmountPence:   100,
	})

	if err != ErrToAccountRequired {
		t.Fatalf("expected error for missing to account, got nil")
	}
}

func TestCreateTransferReturnsErrorWhenAmountIsZero(t *testing.T) {
	store := fakeStore{}
	service := NewService(&store)
	_, err := service.CreateTransfer(t.Context(), CreateTransferInput{
		FromAccountID: "from-account-id",
		ToAccountID:   "to-account-id",
		AmountPence:   0,
	})

	if err != ErrAmountMustBePositive {
		t.Fatalf("expected error for zero amount, got nil")
	}
}

func TestCreateTransferReturnsErrorWhenAccountsAreTheSame(t *testing.T) {
	store := fakeStore{}
	service := NewService(&store)
	_, err := service.CreateTransfer(t.Context(), CreateTransferInput{
		FromAccountID: "same-account-id",
		ToAccountID:   "same-account-id",
		AmountPence:   100,
	})

	if err != ErrSameAccount {
		t.Fatalf("expected error for same accounts, got nil")
	}
}

func TestCreateTransferReturnsErrorWhenInsufficientFunds(t *testing.T) {
	store := fakeStore{
		err: ErrInsufficientFunds,
	}

	service := NewService(&store)
	_, err := service.CreateTransfer(t.Context(), CreateTransferInput{
		FromAccountID: "from-account-id",
		ToAccountID:   "to-account-id",
		AmountPence:   100,
	})

	if err != ErrInsufficientFunds {
		t.Fatalf("expected error for insufficient funds, got nil")
	}
}

func TestCreateTransferReturnsTransferWhenInputIsValid(t *testing.T) {
	store := fakeStore{
		transfer: Transfer{
			ID:            "transfer-id",
			FromAccountID: "from-account-id",
			ToAccountID:   "to-account-id",
			AmountPence:   100,
		},
	}

	service := NewService(&store)
	transfer, err := service.CreateTransfer(t.Context(), CreateTransferInput{
		FromAccountID: "from-account-id",
		ToAccountID:   "to-account-id",
		AmountPence:   100,
	})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if transfer.FromAccountID != "from-account-id" {
		t.Fatalf("expected from account id to be 'from-account-id', got %v", transfer.FromAccountID)
	}

	if transfer.ToAccountID != "to-account-id" {
		t.Fatalf("expected to account id to be 'to-account-id', got %v", transfer.ToAccountID)
	}

	if transfer.AmountPence != 100 {
		t.Fatalf("expected amount to be 100, got %v", transfer.AmountPence)
	}

	if store.createdID == "" {
		t.Fatal("expected generated transfer id")
	}

	if store.createdInput.FromAccountID != "from-account-id" {
		t.Fatalf("expected from account id to be 'from-account-id', got %v", store.createdInput.FromAccountID)
	}

	if store.createdInput.ToAccountID != "to-account-id" {
		t.Fatalf("expected to account id to be 'to-account-id', got %v", store.createdInput.ToAccountID)
	}

	if store.createdInput.AmountPence != 100 {
		t.Fatalf("expected amount pence 100, got %d", store.createdInput.AmountPence)
	}
}
