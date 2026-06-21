package accounts

import (
	"context"
	"testing"
)

type fakeStore struct {
	createdOwnerName           string
	createdOpeningBalancePence int64
	createAccount              Account
	createErr                  error
	getAccount                 Account
	getErr                     error
}

func (f *fakeStore) Create(ctx context.Context, id string, ownerName string, openingBalancePence int64) (Account, error) {
	f.createdOwnerName = ownerName
	f.createdOpeningBalancePence = openingBalancePence
	return f.createAccount, f.createErr
}

func (f *fakeStore) GetByID(ctx context.Context, id string) (Account, error) {
	return f.getAccount, f.getErr
}

func TestCreateAccountRejectsEmptyOwnerName(t *testing.T) {
	ownerName := ""
	store := fakeStore{}
	service := NewService(&store)
	_, err := service.CreateAccount(t.Context(), CreateAccountInput{
		OwnerName: ownerName,
	})

	if err != ErrOwnerNameRequired {
		t.Fatalf("expected error for empty owner name, got nil")
	}
}

func TestCreateAccountTrimsOwnerName(t *testing.T) {
	ownerName := "  Fake Name  "
	store := fakeStore{}
	service := NewService(&store)
	_, err := service.CreateAccount(t.Context(), CreateAccountInput{
		OwnerName: ownerName,
	})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if store.createdOwnerName != "Fake Name" {
		t.Fatalf("expected trimmed owner name to be 'Fake Name', got %q", store.createdOwnerName)
	}
}

func TestCreateAccountRejectsNegativeOpeningBalance(t *testing.T) {
	ownerName := "  Fake Name  "
	store := fakeStore{}
	service := NewService(&store)
	_, err := service.CreateAccount(t.Context(), CreateAccountInput{
		OwnerName:           ownerName,
		OpeningBalancePence: -100,
	})

	if err != ErrOpeningBalanceNegative {
		t.Fatalf("expected error for negative opening balance, got nil")
	}
}

func TestCreateAccountCallsStoreAndReturnsAccount(t *testing.T) {
	ownerName := "Fake Name"
	store := fakeStore{
		createAccount: Account{ID: "1", OwnerName: ownerName, Currency: "GBP", BalancePence: 100},
	}
	service := NewService(&store)
	account, err := service.CreateAccount(t.Context(), CreateAccountInput{
		OwnerName:           ownerName,
		OpeningBalancePence: 100,
	})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if account.ID != "1" {
		t.Fatalf("expected account ID to be '1', got %q", account.ID)
	}

	if account.Currency != "GBP" {
		t.Fatalf("expected account currency to be 'GBP', got %q", account.Currency)
	}

	if account.BalancePence != 100 {
		t.Fatalf("expected account balance to be 100, got %d", account.BalancePence)
	}

	if store.createdOpeningBalancePence != 100 {
		t.Fatalf("expected opening balance %d, got %d", 100, store.createdOpeningBalancePence)
	}
}

func TestGetAccountByIDReturnsAccount(t *testing.T) {
	store := fakeStore{
		getAccount: Account{ID: "1", OwnerName: "Fake Name"},
	}
	service := NewService(&store)
	account, err := service.GetAccountByID(t.Context(), "1")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if account.ID != "1" {
		t.Fatalf("expected account ID to be '1', got %q", account.ID)
	}

	if account.OwnerName != "Fake Name" {
		t.Fatalf("expected account owner name to be 'Fake Name', got %q", account.OwnerName)
	}
}
