package httpapi

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/ABAlmeida/pennyledger/internal/accounts"
)

type accountService interface {
	CreateAccount(ctx context.Context, input accounts.CreateAccountInput) (accounts.Account, error)
	GetAccountByID(ctx context.Context, id string) (accounts.Account, error)
}

type createAccountRequest struct {
	OwnerName           string `json:"owner_name"`
	OpeningBalancePence int64  `json:"opening_balance_pence"`
}

type accountResponse struct {
	ID           string `json:"id"`
	OwnerName    string `json:"owner_name"`
	Currency     string `json:"currency"`
	BalancePence int64  `json:"balance_pence"`
	Status       string `json:"status"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
}

func createAccountHandler(service accountService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var request createAccountRequest

		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			http.Error(w, "invalid json", http.StatusBadRequest)
			return
		}

		account, err := service.CreateAccount(r.Context(), accounts.CreateAccountInput{
			OwnerName:           request.OwnerName,
			OpeningBalancePence: request.OpeningBalancePence,
		})

		if errors.Is(err, accounts.ErrOwnerNameRequired) {
			http.Error(w, "owner_name is required", http.StatusBadRequest)
			return
		}

		if errors.Is(err, accounts.ErrOpeningBalanceNegative) {
			http.Error(w, "opening_balance_pence cannot be negative", http.StatusBadRequest)
			return
		}

		if err != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		writeJSON(w, http.StatusCreated, toAccountResponse(account))
	}
}

func getAccountHandler(service accountService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		account, err := service.GetAccountByID(r.Context(), r.PathValue("id"))
		if errors.Is(err, accounts.ErrNotFound) {
			http.Error(w, "account not found", http.StatusNotFound)
			return
		}
		if err != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		writeJSON(w, http.StatusOK, toAccountResponse(account))
	}
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}

func toAccountResponse(account accounts.Account) accountResponse {
	return accountResponse{
		ID:           account.ID,
		OwnerName:    account.OwnerName,
		Currency:     account.Currency,
		BalancePence: account.BalancePence,
		Status:       account.Status,
		CreatedAt:    account.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:    account.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}
