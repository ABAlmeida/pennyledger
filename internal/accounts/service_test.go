package accounts

import (
	"context"
	"testing"
)

type fakeStore struct {
	createdOwnerName string
	createAccount    Account
	createErr        error
	getAccount       Account
	getErr           error
}

func (f *fakeStore) Create(ctx context.Context, id string, ownerName string) (Account, error) {
	f.createdOwnerName = ownerName
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

func TestCreateAccountCallsStoreAndReturnsAccount(t *testing.T) {
	ownerName := "Fake Name"
	store := fakeStore{
		createAccount: Account{ID: "1", OwnerName: ownerName, Currency: "GBP", BalancePence: 0},
	}
	service := NewService(&store)
	account, err := service.CreateAccount(t.Context(), CreateAccountInput{
		OwnerName: ownerName,
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

	if account.BalancePence != 0 {
		t.Fatalf("expected account balance to be 0, got %d", account.BalancePence)
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
