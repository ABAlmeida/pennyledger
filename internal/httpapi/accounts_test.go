package httpapi

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/ABAlmeida/pennyledger/internal/accounts"
)

type fakeHTTPAccountService struct {
	createInput   accounts.CreateAccountInput
	createAccount accounts.Account
	createErr     error

	getID      string
	getAccount accounts.Account
	getErr     error
}

func (f *fakeHTTPAccountService) CreateAccount(ctx context.Context, input accounts.CreateAccountInput) (accounts.Account, error) {
	f.createInput = input
	return f.createAccount, f.createErr
}

func (f *fakeHTTPAccountService) GetAccountByID(ctx context.Context, id string) (accounts.Account, error) {
	f.getID = id
	return f.getAccount, f.getErr
}

func TestCreateAccountReturnsCreatedAccount(t *testing.T) {
	createdAt := time.Date(2026, 6, 14, 10, 0, 0, 0, time.UTC)
	updatedAt := time.Date(2026, 6, 14, 10, 1, 0, 0, time.UTC)

	service := &fakeHTTPAccountService{
		createAccount: accounts.Account{
			ID:           "id",
			OwnerName:    "Alice",
			Currency:     "GBP",
			BalancePence: 0,
			Status:       "active",
			CreatedAt:    createdAt,
			UpdatedAt:    updatedAt,
		},
	}

	request := httptest.NewRequest(
		http.MethodPost,
		"/v1/accounts",
		strings.NewReader(`{"owner_name":"Alice"}`),
	)
	response := httptest.NewRecorder()

	createAccountHandler(service).ServeHTTP(response, request)

	if response.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, response.Code)
	}

	if service.createInput.OwnerName != "Alice" {
		t.Fatalf("expected owner name %q, got %q", "Alice", service.createInput.OwnerName)
	}

	var body accountResponse
	if err := json.NewDecoder(response.Body).Decode(&body); err != nil {
		t.Fatalf("expected valid json response, got error %v", err)
	}

	if body.ID != "id" {
		t.Fatalf("expected account id %q, got %q", "id", body.ID)
	}

	if body.OwnerName != "Alice" {
		t.Fatalf("expected owner name %q, got %q", "Alice", body.OwnerName)
	}
}

func TestCreateAccountReturnsBadRequestForInvalidJSON(t *testing.T) {
	service := &fakeHTTPAccountService{}

	request := httptest.NewRequest(
		http.MethodPost,
		"/v1/accounts",
		strings.NewReader(`{"owner_name":`),
	)
	response := httptest.NewRecorder()

	createAccountHandler(service).ServeHTTP(response, request)

	if response.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, response.Code)
	}

	expectedBodyResponse := "invalid json\n"
	if response.Body.String() != expectedBodyResponse {
		t.Fatalf("expected error message %q, got %q", expectedBodyResponse, response.Body.String())
	}
}

func TestCreateAccountReturnsBadRequestForEmptyOwnerName(t *testing.T) {
	service := &fakeHTTPAccountService{
		createErr: accounts.ErrOwnerNameRequired,
	}

	request := httptest.NewRequest(
		http.MethodPost,
		"/v1/accounts",
		strings.NewReader(`{"owner_name":""}`),
	)
	response := httptest.NewRecorder()

	createAccountHandler(service).ServeHTTP(response, request)

	if response.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, response.Code)
	}

	expectedBodyResponse := "owner_name is required\n"
	if response.Body.String() != expectedBodyResponse {
		t.Fatalf("expected error message %q, got %q", expectedBodyResponse, response.Body.String())
	}
}

func TestCreateAccountReturnsBadRequestForNegativeOpeningBalance(t *testing.T) {
	service := &fakeHTTPAccountService{
		createErr: accounts.ErrOpeningBalanceNegative,
	}

	request := httptest.NewRequest(
		http.MethodPost,
		"/v1/accounts",
		strings.NewReader(`{"owner_name":"Alice","opening_balance_pence":-100}`),
	)
	response := httptest.NewRecorder()

	createAccountHandler(service).ServeHTTP(response, request)

	if response.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, response.Code)
	}

	expectedBodyResponse := "opening_balance_pence cannot be negative\n"
	if response.Body.String() != expectedBodyResponse {
		t.Fatalf("expected error message %q, got %q", expectedBodyResponse, response.Body.String())
	}
}

func TestCreateAccountReturnsInternalServerErrorForServiceError(t *testing.T) {
	service := &fakeHTTPAccountService{
		createErr: errors.New("database failed"),
	}

	request := httptest.NewRequest(
		http.MethodPost,
		"/v1/accounts",
		strings.NewReader(`{"owner_name":"Alice"}`),
	)
	response := httptest.NewRecorder()

	createAccountHandler(service).ServeHTTP(response, request)

	if response.Code != http.StatusInternalServerError {
		t.Fatalf("expected status %d, got %d", http.StatusInternalServerError, response.Code)
	}

	expectedBodyResponse := "internal server error\n"
	if response.Body.String() != expectedBodyResponse {
		t.Fatalf("expected error message %q, got %q", expectedBodyResponse, response.Body.String())
	}
}

func TestGetAccountByIDReturnsAccount(t *testing.T) {
	createdAt := time.Date(2026, 6, 14, 10, 0, 0, 0, time.UTC)
	updatedAt := time.Date(2026, 6, 14, 10, 1, 0, 0, time.UTC)

	service := &fakeHTTPAccountService{
		getAccount: accounts.Account{
			ID:           "account-id",
			OwnerName:    "Alice",
			Currency:     "GBP",
			BalancePence: 0,
			Status:       "active",
			CreatedAt:    createdAt,
			UpdatedAt:    updatedAt,
		},
	}

	request := httptest.NewRequest(
		http.MethodGet,
		"/v1/accounts/account-id",
		nil,
	)
	response := httptest.NewRecorder()

	getAccountHandler(service).ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, response.Code)
	}

	var body accountResponse
	if err := json.NewDecoder(response.Body).Decode(&body); err != nil {
		t.Fatalf("expected valid json response, got error %v", err)
	}

	if body.ID != "account-id" {
		t.Fatalf("expected account id %q, got %q", "account-id", body.ID)
	}

	if body.OwnerName != "Alice" {
		t.Fatalf("expected owner name %q, got %q", "Alice", body.OwnerName)
	}
}

func TestGetAccountByIDReturnsNotFoundForNonExistentAccount(t *testing.T) {
	service := &fakeHTTPAccountService{
		getErr: accounts.ErrNotFound,
	}

	request := httptest.NewRequest(
		http.MethodGet,
		"/v1/accounts/non-existent-id",
		nil,
	)
	response := httptest.NewRecorder()

	getAccountHandler(service).ServeHTTP(response, request)

	if response.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d", http.StatusNotFound, response.Code)
	}

	expectedBodyResponse := "account not found\n"
	if response.Body.String() != expectedBodyResponse {
		t.Fatalf("expected error message %q, got %q", expectedBodyResponse, response.Body.String())
	}
}

func TestGetAccountByIDReturnsInternalServerErrorForServiceError(t *testing.T) {
	service := &fakeHTTPAccountService{
		getErr: errors.New("database failed"),
	}

	request := httptest.NewRequest(
		http.MethodGet,
		"/v1/accounts/account-id",
		nil,
	)
	response := httptest.NewRecorder()

	getAccountHandler(service).ServeHTTP(response, request)

	if response.Code != http.StatusInternalServerError {
		t.Fatalf("expected status %d, got %d", http.StatusInternalServerError, response.Code)
	}

	expectedBodyResponse := "internal server error\n"
	if response.Body.String() != expectedBodyResponse {
		t.Fatalf("expected error message %q, got %q", expectedBodyResponse, response.Body.String())
	}
}
